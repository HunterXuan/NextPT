package tracker

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/accountingin"
	"server/internal/model/in/trackerin"
	"server/internal/service"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
)

type sTrackerEventUsecase struct {
	consumerGroup string
}

const (
	announceDedupTTLSeconds = 30
	announceQueueMaxLen     = 100000
)

const pushAnnounceEventLuaScript = `
if redis.call("EXISTS", KEYS[2]) == 1 then
	return ""
end
local msgId = redis.call("XADD", KEYS[1], "MAXLEN", "~", ARGV[1], "*", "data", ARGV[3])
redis.call("SET", KEYS[2], msgId, "EX", ARGV[2])
return msgId
`

func init() {
	service.RegisterTrackerEventUsecase(NewTrackerEventUsecase())
}

func NewTrackerEventUsecase() *sTrackerEventUsecase {
	s := &sTrackerEventUsecase{
		consumerGroup: "accounting_group",
	}

	// Initialize Consumer Group
	ctx := context.Background()
	queueKey := service.SysCache().KeyTrackerAnnounceQueue(ctx)
	_, err := g.Redis().Do(ctx, "XGROUP", "CREATE", queueKey, s.consumerGroup, "0", "MKSTREAM")
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		glog.Errorf(ctx, "Failed to create consumer group: %v", err)
	}

	// 启动常驻后台的异步消费协程池 (Redis Stream)
	// 在多实例部署时，需确保每个实例的 consumerName 唯一，否则会共享 PEL 导致抢占冲突或 XAUTOCLAIM 失效
	instanceId := guid.S()
	for i := 0; i < 10; i++ {
		go s.consumeEvents(fmt.Sprintf("consumer-%s-%d", instanceId, i))
	}
	// 启动死信/超时重试协程
	go s.watchPendingEvents(fmt.Sprintf("pending_watcher-%s", instanceId))

	return s
}

// PushAnnounceEvent 接收并投递 Announce 事件到队列 (Redis Stream)
func (s *sTrackerEventUsecase) PushAnnounceEvent(event *trackerin.AnnounceEvent) error {
	ctx := context.Background()
	data, err := json.Marshal(event)
	if err != nil {
		glog.Errorf(ctx, "failed to marshal announce event: %v", err)
		return err
	}

	queueKey := service.SysCache().KeyTrackerAnnounceQueue(ctx)
	dedupKey := s.announceDedupKey(ctx, event)
	_, err = g.Redis().Do(ctx, "EVAL", pushAnnounceEventLuaScript, 2, queueKey, dedupKey, announceQueueMaxLen, announceDedupTTLSeconds, data)
	if err != nil {
		glog.Errorf(ctx, "failed to push announce event to redis stream: %v", err)
		return err
	}
	return nil
}

func (s *sTrackerEventUsecase) announceDedupKey(ctx context.Context, event *trackerin.AnnounceEvent) string {
	h := sha1.New()
	_, _ = fmt.Fprintf(
		h,
		"%d:%d:%s:%s:%d:%d:%d:%t:%s:%s:%d",
		event.TorrentId,
		event.UserId,
		event.PeerId,
		event.Event,
		event.Uploaded,
		event.Downloaded,
		event.Left,
		event.IsSeeder,
		event.Ipv4,
		event.Ipv6,
		event.Port,
	)
	return service.SysCache().KeyTrackerAnnounceDedup(ctx, hex.EncodeToString(h.Sum(nil)))
}

// consumeEvents 消费协程 (阻塞拉取 Redis Stream)
func (s *sTrackerEventUsecase) consumeEvents(consumerName string) {
	ctx := context.Background()
	queueKey := service.SysCache().KeyTrackerAnnounceQueue(ctx)

	for {
		// XREADGROUP GROUP accounting_group consumerName BLOCK 0 COUNT 10 STREAMS queueKey >
		res, err := g.Redis().Do(ctx, "XREADGROUP", "GROUP", s.consumerGroup, consumerName, "BLOCK", 0, "COUNT", 10, "STREAMS", queueKey, ">")
		if err != nil {
			glog.Errorf(ctx, "redis xreadgroup error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if res.IsNil() {
			continue
		}

		arr := res.Interfaces()
		if len(arr) == 0 {
			continue
		}

		streamData := gconv.Interfaces(arr[0])
		if len(streamData) < 2 {
			continue
		}

		messages := gconv.Interfaces(streamData[1])
		for _, msgInfo := range messages {
			msgData := gconv.Interfaces(msgInfo)
			msgId := gconv.String(msgData[0])
			kv := gconv.Strings(msgData[1])

			if len(kv) >= 2 {
				var event trackerin.AnnounceEvent
				if err := json.Unmarshal([]byte(kv[1]), &event); err != nil {
					// Invalid JSON, drop it
					g.Redis().Do(ctx, "XACK", queueKey, s.consumerGroup, msgId)
					g.Redis().Do(ctx, "XDEL", queueKey, msgId)
					continue
				}

				if err := s.handleAnnounceEvent(msgId, &event); err == nil {
					g.Redis().Do(ctx, "XACK", queueKey, s.consumerGroup, msgId)
					g.Redis().Do(ctx, "XDEL", queueKey, msgId)
				} else {
					s.handleDlq(ctx, queueKey, msgId, event)
				}
			}
		}
	}
}

func (s *sTrackerEventUsecase) watchPendingEvents(consumerName string) {
	ctx := context.Background()
	queueKey := service.SysCache().KeyTrackerAnnounceQueue(ctx)

	for {
		time.Sleep(30 * time.Second)

		// 既然已升级到 Redis 7.0，可直接使用更高效的原生 XAUTOCLAIM
		// XAUTOCLAIM queueKey consumerGroup consumerName min-idle-time start count
		// 60000ms (60s) timeout
		res, err := g.Redis().Do(ctx, "XAUTOCLAIM", queueKey, s.consumerGroup, consumerName, 60000, "0-0", "COUNT", 100)
		if err != nil {
			glog.Errorf(ctx, "redis xautoclaim error: %v", err)
			continue
		}

		if res.IsNil() {
			continue
		}

		arr := res.Interfaces()
		if len(arr) < 2 {
			continue
		}

		messages := gconv.Interfaces(arr[1])
		for _, msgInfo := range messages {
			msgData := gconv.Interfaces(msgInfo)
			msgId := gconv.String(msgData[0])
			kv := gconv.Strings(msgData[1])

			if len(kv) >= 2 {
				var event trackerin.AnnounceEvent
				if err := json.Unmarshal([]byte(kv[1]), &event); err != nil {
					g.Redis().Do(ctx, "XACK", queueKey, s.consumerGroup, msgId)
					g.Redis().Do(ctx, "XDEL", queueKey, msgId)
					continue
				}

				if err := s.handleAnnounceEvent(msgId, &event); err == nil {
					g.Redis().Do(ctx, "XACK", queueKey, s.consumerGroup, msgId)
					g.Redis().Do(ctx, "XDEL", queueKey, msgId)
				} else {
					s.handleDlq(ctx, queueKey, msgId, event)
				}
			}
		}
	}
}

func (s *sTrackerEventUsecase) handleDlq(ctx context.Context, queueKey, msgId string, event trackerin.AnnounceEvent) {
	dlqCountKey := service.SysCache().KeyTrackerDlqCounts(ctx)
	res, err := g.Redis().Do(ctx, "HINCRBY", dlqCountKey, msgId, 1)
	if err != nil {
		glog.Errorf(ctx, "redis hincrby dlq count error: %v", err)
		return
	}

	count := res.Int()
	if count >= 3 {
		// DLQ threshold reached
		dlqQueueKey := service.SysCache().KeyTrackerAnnounceQueueDlq(ctx)
		glog.Warningf(ctx, "Event %s failed 3 times, moving to DLQ", msgId)

		data, _ := json.Marshal(event)
		g.Redis().Do(ctx, "XADD", dlqQueueKey, "MAXLEN", "~", 10000, "*", "data", data, "original_id", msgId)

		// Ack and Del from main queue
		g.Redis().Do(ctx, "XACK", queueKey, s.consumerGroup, msgId)
		g.Redis().Do(ctx, "XDEL", queueKey, msgId)
		// Clean up dlq count
		g.Redis().Do(ctx, "HDEL", dlqCountKey, msgId)
	}
}

// handleAnnounceEvent 协调 MySQL 账务与 Redis 状态的更新
func (s *sTrackerEventUsecase) handleAnnounceEvent(msgId string, event *trackerin.AnnounceEvent) error {
	ctx := context.Background()

	if event.TorrentId == 0 {
		return fmt.Errorf("announce event missing torrent id")
	}
	if event.Now == nil {
		return fmt.Errorf("announce event missing event time")
	}

	// 0. 获取 Redis 分布式锁，防止并发刷流双花攻击
	// 锁冲突时必须返回 error，让消息留在 pending 中，等锁释放后由 XAUTOCLAIM 重新认领处理
	if !s.tryAcquireLock(ctx, event) {
		return fmt.Errorf("peer lock conflict, will retry via XAUTOCLAIM")
	}

	// 2. 提取历史状态
	var oldPeer *entity.TrackerPeer
	peerDetailKey := service.SysCache().KeyTrackerPeerDetail(ctx, event.TorrentId, event.PeerId)
	detailVal, err := g.Redis().Do(ctx, "GET", peerDetailKey)
	if err != nil {
		glog.Error(ctx, "async announce event failed to get old peer: ", err)
		return err // Let it retry
	}
	if !detailVal.IsNil() && detailVal.String() != "" {
		oldPeer = &entity.TrackerPeer{}
		if err := json.Unmarshal(detailVal.Bytes(), oldPeer); err != nil {
			glog.Error(ctx, "async announce event failed to unmarshal old peer: ", err)
			return err
		}
	}
	if s.isStaleAnnounceEvent(event, oldPeer) {
		glog.Infof(ctx, "stale announce event skipped [torrent:%d, user:%d, peer:%s, event_time:%s, last_action:%s]",
			event.TorrentId, event.UserId, event.PeerId, event.Now, oldPeer.LastAction)
		return nil
	}

	// 3. 计算账务增量与时间增量
	diffUp, diffDn := service.TrackerPeerDomain().CalculateTrafficDiff(ctx, event, oldPeer)

	timeDiff := 0
	if oldPeer != nil && oldPeer.LastAction != nil {
		timeDiff = max(int(event.Now.Unix()-oldPeer.LastAction.Unix()), 0)
	}

	// 3. MySQL 账务核心事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 3.1 幂等写入，防重放
		res, err := dao.TrackerEventIdempotency.Ctx(ctx).InsertIgnore(do.TrackerEventIdempotency{
			MsgId:     msgId,
			CreatedAt: event.Now,
		})
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil || affected == 0 {
			return nil // already processed
		}
		finishedTransition, err := service.AccountingSnatchDomain().RecordSnatch(ctx, accountingin.RecordSnatchInp{
			TorrentId:      event.TorrentId,
			UserId:         event.UserId,
			Ipv4:           event.Ipv4,
			Ipv6:           event.Ipv6,
			Port:           event.Port,
			UploadedDiff:   diffUp,
			DownloadedDiff: diffDn,
			Remaining:      event.Left,
			IsSeeder:       event.IsSeeder,
			IsFinished:     event.IsSeeder,
			EventTime:      event.Now,
			TimeDiff:       timeDiff,
		})
		if err != nil {
			return err
		}
		if err := service.AccountingTrafficDomain().RecordTraffic(ctx, event.UserId, diffUp, diffDn, event.IsSeeder, timeDiff, event.Now); err != nil {
			return err
		}
		if event.Event == consts.TrackerAnnounceEventCompleted && finishedTransition {
			if err := service.CatalogTorrentDomain().IncrementTorrentStats(ctx, event.TorrentId, "times_completed", 1); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		glog.Errorf(ctx, "mysql accounting failed for announce event [torrent:%d, user:%d]: %v", event.TorrentId, event.UserId, err)
		return err // Retry on failure
	}

	// 5. Redis 节点与在线态维护
	switch event.Event {
	case consts.TrackerAnnounceEventStopped:
		if err := service.TrackerPeerDomain().RemovePeer(ctx, event.TorrentId, event.UserId, event.PeerId); err != nil {
			glog.Errorf(ctx, "redis remove peer failed: %v", err)
			return err
		}
	default:
		interval := s.getTrackerConfigCache(ctx, consts.SiteConfigTrackerAnnounceInterval).Int()
		if err := service.TrackerPeerDomain().UpsertPeer(ctx, event, oldPeer, interval); err != nil {
			glog.Errorf(ctx, "redis upsert peer failed: %v", err)
			return err
		}
	}

	return nil
}

func (s *sTrackerEventUsecase) isStaleAnnounceEvent(event *trackerin.AnnounceEvent, oldPeer *entity.TrackerPeer) bool {
	if oldPeer == nil || oldPeer.LastAction == nil {
		return false
	}
	return event.Now.Time.Before(oldPeer.LastAction.Time)
}

// tryAcquireLock 尝试获取 Redis 分布式锁，用于防并发双花
func (s *sTrackerEventUsecase) tryAcquireLock(ctx context.Context, event *trackerin.AnnounceEvent) bool {
	lockKey := service.SysCache().KeyTrackerLockPeer(ctx, event.TorrentId, event.UserId, event.PeerId)
	v, err := g.Redis().Do(ctx, "SET", lockKey, 1, "NX", "EX", 5)
	if err != nil {
		glog.Error(ctx, "redis lock error: ", err)
		return false
	}
	if v.IsNil() || v.String() != "OK" {
		glog.Warningf(ctx, "concurrent request dropped for torrent:%d user:%d peer_id: %s", event.TorrentId, event.UserId, event.PeerId)
		return false
	}
	return true
}

func (s *sTrackerEventUsecase) getTrackerConfigCache(ctx context.Context, key string) *gvar.Var {
	cacheKey := service.SysCache().KeySiteConfigFullPath(ctx, key)
	val, err := gcache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		return service.SiteConfigDomain().GetByPath(ctx, key).Val(), nil
	}, 5*time.Minute)
	if err != nil || val.IsNil() {
		return service.SiteConfigDomain().GetByPath(ctx, key)
	}
	return gvar.New(val.Val())
}
