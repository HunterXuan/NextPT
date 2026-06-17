package tracker

import (
	"context"
	"encoding/json"
	"time"

	libtracker "server/internal/library/tracker"
	"server/internal/model/entity"
	"server/internal/service"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model/in/trackerin"

	"github.com/gogf/gf/v2/os/glog"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type sTrackerPeerDomain struct{}

func init() {
	service.RegisterTrackerPeerDomain(NewTrackerPeerDomain())
}

func NewTrackerPeerDomain() *sTrackerPeerDomain {
	return &sTrackerPeerDomain{}
}

// GetActivePeers 获取种子当前在线的全部做种者和下载者
func (s *sTrackerPeerDomain) GetActivePeers(ctx context.Context, torrentId uint64) ([]*entity.TrackerPeer, error) {
	now := time.Now().Unix()

	seedersKey := service.SysCache().KeyTrackerTorrentSeeders(ctx, torrentId)
	leechersKey := service.SysCache().KeyTrackerTorrentLeechers(ctx, torrentId)

	sValues, err := g.Redis().Do(ctx, "ZRANGEBYSCORE", seedersKey, now, "+inf")
	if err != nil {
		return nil, err
	}
	lValues, err := g.Redis().Do(ctx, "ZRANGEBYSCORE", leechersKey, now, "+inf")
	if err != nil {
		return nil, err
	}

	var peerIds []string
	if !sValues.IsNil() {
		peerIds = append(peerIds, sValues.Strings()...)
	}
	if !lValues.IsNil() {
		peerIds = append(peerIds, lValues.Strings()...)
	}

	if len(peerIds) == 0 {
		return nil, nil
	}

	// 去重并提取 peerId
	var detailKeys []any
	seenKeys := make(map[string]bool)

	for _, member := range peerIds {
		_, pId, ok := libtracker.PeerParseZsetMember(member)
		if !ok {
			continue
		}

		key := service.SysCache().KeyTrackerPeerDetail(ctx, torrentId, pId)
		if !seenKeys[key] {
			seenKeys[key] = true
			detailKeys = append(detailKeys, key)
		}
	}

	if len(detailKeys) == 0 {
		return nil, nil
	}

	// 批量拉取详情
	bytesList, err := g.Redis().Do(ctx, "MGET", detailKeys...)
	if err != nil {
		return nil, err
	}

	var activePeers []*entity.TrackerPeer
	for _, val := range bytesList.Vars() {
		if val.IsNil() || val.String() == "" {
			continue
		}
		var peer entity.TrackerPeer
		if err := json.Unmarshal(val.Bytes(), &peer); err != nil {
			continue
		}
		activePeers = append(activePeers, &peer)
	}

	return activePeers, nil
}

func (s *sTrackerPeerDomain) GetSeedingUsers(ctx context.Context) ([]uint64, error) {
	seedingUsersKey := service.SysCache().KeyTrackerSeedingUsers(ctx)
	res, err := g.Redis().Do(ctx, "SMEMBERS", seedingUsersKey)
	if err != nil || res.IsNil() {
		return nil, err
	}
	var activeUserIds []uint64
	for _, val := range res.Strings() {
		if uId := gconv.Uint64(val); uId > 0 {
			activeUserIds = append(activeUserIds, uId)
		}
	}
	return activeUserIds, nil
}

func (s *sTrackerPeerDomain) GetUserSeedingPeers(ctx context.Context, userId uint64) ([]entity.TrackerPeer, error) {
	userSeedingKey := service.SysCache().KeyTrackerUserSeeding(ctx, userId)
	members, err := g.Redis().Do(ctx, "SMEMBERS", userSeedingKey)
	if err != nil || members.IsNil() || len(members.Strings()) == 0 {
		return nil, err
	}

	var detailKeys []any
	for _, m := range members.Strings() {
		tId, pId, ok := libtracker.PeerParseUserSetValue(m)
		if ok {
			detailKeys = append(detailKeys, service.SysCache().KeyTrackerPeerDetail(ctx, tId, pId))
		}
	}

	if len(detailKeys) == 0 {
		return nil, nil
	}

	bytesList, err := g.Redis().Do(ctx, "MGET", detailKeys...)
	if err != nil {
		return nil, err
	}

	var seeders []entity.TrackerPeer
	for _, val := range bytesList.Vars() {
		if val.IsNil() || val.String() == "" {
			continue
		}
		var peer entity.TrackerPeer
		if err := json.Unmarshal(val.Bytes(), &peer); err == nil {
			seeders = append(seeders, peer)
		}
	}
	return seeders, nil
}

// GetEnabledClientWhitelists 获取所有启用的客户端白名单规则
func (s *sTrackerPeerDomain) GetEnabledClientWhitelists(ctx context.Context) ([]*entity.TrackerAgentWhitelist, error) {
	var list []*entity.TrackerAgentWhitelist
	err := dao.TrackerAgentWhitelist.Ctx(ctx).Where(dao.TrackerAgentWhitelist.Columns().Enabled, true).Scan(&list)
	return list, err
}

// calculateTrafficDiff 计算本次汇报的上传下载增量
func (s *sTrackerPeerDomain) CalculateTrafficDiff(ctx context.Context, event *trackerin.AnnounceEvent, oldPeer *entity.TrackerPeer) (diffUp, diffDn int64) {
	if oldPeer == nil {
		// 防刷核心：首次出现的 Peer，由于无法计算时间差（绕过防作弊限速），并且为了防止跨站偷流量，必须返回 0 增量。
		// 此次汇报的绝对值只作为后续计算的基准水位线。
		return 0, 0
	}

	diffUp = event.Uploaded - gconv.Int64(oldPeer.Uploaded)
	diffDn = event.Downloaded - gconv.Int64(oldPeer.Downloaded)

	// 处理客户端重启重置或回退的异常情况
	if diffUp < 0 {
		diffUp = event.Uploaded
	}
	if diffDn < 0 {
		diffDn = event.Downloaded
	}

	// 1. 防作弊检查：硬限速 10MB/s
	if oldPeer.LastAction != nil {
		timeDiff := event.Now.Unix() - oldPeer.LastAction.Unix()
		maxSpeedBytes := int64(10 * 1024 * 1024) // 10MB/s

		var currentSpeed int64
		if timeDiff > 0 {
			currentSpeed = (diffUp + diffDn) / timeDiff
		} else {
			// 同一秒内或时钟回退，如果是小量可以放行，大流量直接判定超速
			if diffUp+diffDn > maxSpeedBytes {
				currentSpeed = maxSpeedBytes + 1
			}
		}

		if currentSpeed > maxSpeedBytes {
			glog.Warningf(context.Background(), "[Anti-Cheat] User %d Torrent %d Peer %s Exceeded 10MB/s limit. Up: %d, Dn: %d, timeDiff: %ds",
				event.UserId, event.TorrentId, event.PeerId, diffUp, diffDn, timeDiff)
			// 超速流量全部清零（静默丢弃，不报错，保护数据纯净）
			return 0, 0
		}
	}

	return diffUp, diffDn
}

// 辅助方法：维护 Redis Peer 缓存与 ZSET 索引以及用户做种下载索引
func (s *sTrackerPeerDomain) UpsertPeer(ctx context.Context, event *trackerin.AnnounceEvent, oldPeer *entity.TrackerPeer, interval int) error {
	var peer entity.TrackerPeer
	if oldPeer != nil {
		peer = *oldPeer
	} else {
		peer.TorrentId = event.TorrentId
		peer.UserId = event.UserId
		peer.PeerId = []byte(event.PeerId)
		peer.StartedAt = event.Now
		peer.IsConnectable = true
	}

	peer.Ipv4 = event.Ipv4
	peer.Ipv6 = event.Ipv6
	peer.Port = uint(event.Port)
	peer.Uploaded = uint64(event.Uploaded)
	peer.Downloaded = uint64(event.Downloaded)
	peer.Remaining = uint64(event.Left)
	peer.IsSeeder = event.IsSeeder
	peer.Agent = event.UserAgent
	peer.LastAction = event.Now

	if event.Event == consts.TrackerAnnounceEventCompleted && oldPeer != nil && oldPeer.FinishedAt == nil {
		peer.FinishedAt = event.Now
	}
	if event.IsSeeder && peer.FinishedAt == nil {
		peer.FinishedAt = event.Now
	}

	// 序列化详情并写入 Redis
	peerDetailKey := service.SysCache().KeyTrackerPeerDetail(ctx, event.TorrentId, event.PeerId)
	bytes, err := json.Marshal(peer)
	if err != nil {
		return err
	}

	// 计算过期时间：宣告间隔的 1.5 倍
	if interval <= 0 {
		interval = 1800
	}
	expireSeconds := int(float64(interval) * 1.5)

	_, err = g.Redis().Do(ctx, "SET", peerDetailKey, bytes, "EX", expireSeconds)
	if err != nil {
		return err
	}

	// 维护 ZSET 和用户索引
	seedersKey := service.SysCache().KeyTrackerTorrentSeeders(ctx, event.TorrentId)
	leechersKey := service.SysCache().KeyTrackerTorrentLeechers(ctx, event.TorrentId)
	userSeedingKey := service.SysCache().KeyTrackerUserSeeding(ctx, event.UserId)
	userLeechingKey := service.SysCache().KeyTrackerUserLeeching(ctx, event.UserId)
	activeTorrentsKey := service.SysCache().KeyTrackerActiveTorrents(ctx)
	seedingUsersKey := service.SysCache().KeyTrackerSeedingUsers(ctx)

	expireAt := event.Now.Unix() + int64(expireSeconds)
	userVal := libtracker.PeerFormatUserSetValue(event.TorrentId, event.PeerId)
	zsetMember := libtracker.PeerFormatZsetMember(event.UserId, event.PeerId)

	// 激活种子 ID
	_, _ = g.Redis().Do(ctx, "SADD", activeTorrentsKey, event.TorrentId)

	if event.IsSeeder {
		// 放入做种者列表
		_, err = g.Redis().Do(ctx, "ZADD", seedersKey, expireAt, zsetMember)
		if err != nil {
			return err
		}
		_, err = g.Redis().Do(ctx, "ZREM", leechersKey, zsetMember)
		if err != nil {
			return err
		}

		// 放入做种反向索引
		_, err = g.Redis().Do(ctx, "SADD", userSeedingKey, userVal)
		if err != nil {
			return err
		}
		_, err = g.Redis().Do(ctx, "SREM", userLeechingKey, userVal)
		if err != nil {
			return err
		}

		// 激活做种用户 ID
		_, _ = g.Redis().Do(ctx, "SADD", seedingUsersKey, event.UserId)
	} else {
		// 放入下载者列表
		_, err = g.Redis().Do(ctx, "ZADD", leechersKey, expireAt, zsetMember)
		if err != nil {
			return err
		}
		_, err = g.Redis().Do(ctx, "ZREM", seedersKey, zsetMember)
		if err != nil {
			return err
		}

		// 放入下载反向索引
		_, err = g.Redis().Do(ctx, "SADD", userLeechingKey, userVal)
		if err != nil {
			return err
		}
		_, err = g.Redis().Do(ctx, "SREM", userSeedingKey, userVal)
		if err != nil {
			return err
		}

		// 如果之前也是做种用户，但现在只剩下载了，我们可以检查并移除其做种用户状态
		seedingCard, _ := g.Redis().Do(ctx, "SCARD", userSeedingKey)
		if seedingCard.Int() == 0 {
			_, _ = g.Redis().Do(ctx, "SREM", seedingUsersKey, event.UserId)
		}
	}

	return nil
}

// RemovePeer 清理指定的 Peer 及所有相关索引
func (s *sTrackerPeerDomain) RemovePeer(ctx context.Context, torrentId, userId uint64, peerId string) error {
	peerDetailKey := service.SysCache().KeyTrackerPeerDetail(ctx, torrentId, peerId)
	seedersKey := service.SysCache().KeyTrackerTorrentSeeders(ctx, torrentId)
	leechersKey := service.SysCache().KeyTrackerTorrentLeechers(ctx, torrentId)
	userSeedingKey := service.SysCache().KeyTrackerUserSeeding(ctx, userId)
	userLeechingKey := service.SysCache().KeyTrackerUserLeeching(ctx, userId)
	activeTorrentsKey := service.SysCache().KeyTrackerActiveTorrents(ctx)
	seedingUsersKey := service.SysCache().KeyTrackerSeedingUsers(ctx)

	_, err := g.Redis().Do(ctx, "DEL", peerDetailKey)
	if err != nil {
		return err
	}
	zsetMember := libtracker.PeerFormatZsetMember(userId, peerId)
	_, err = g.Redis().Do(ctx, "ZREM", seedersKey, zsetMember)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "ZREM", leechersKey, zsetMember)
	if err != nil {
		return err
	}
	userVal := libtracker.PeerFormatUserSetValue(torrentId, peerId)
	_, err = g.Redis().Do(ctx, "SREM", userSeedingKey, userVal)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SREM", userLeechingKey, userVal)
	if err != nil {
		return err
	}

	// 检查该种子是否还有任何 peer
	sCard, _ := g.Redis().Do(ctx, "ZCARD", seedersKey)
	lCard, _ := g.Redis().Do(ctx, "ZCARD", leechersKey)
	if sCard.Int() == 0 && lCard.Int() == 0 {
		_, _ = g.Redis().Do(ctx, "SREM", activeTorrentsKey, torrentId)
	}

	// 检查用户是否还有做种
	seedingCard, _ := g.Redis().Do(ctx, "SCARD", userSeedingKey)
	if seedingCard.Int() == 0 {
		_, _ = g.Redis().Do(ctx, "SREM", seedingUsersKey, userId)
	}
	return nil
}

func (s *sTrackerPeerDomain) DeletePeersByTorrentId(ctx context.Context, torrentId uint64) error {
	seedersKey := service.SysCache().KeyTrackerTorrentSeeders(ctx, torrentId)
	leechersKey := service.SysCache().KeyTrackerTorrentLeechers(ctx, torrentId)

	// Get all members to clean up user sets and peer details
	sValues, _ := g.Redis().Do(ctx, "ZRANGE", seedersKey, 0, -1)
	lValues, _ := g.Redis().Do(ctx, "ZRANGE", leechersKey, 0, -1)

	var members []string
	if !sValues.IsNil() {
		members = append(members, sValues.Strings()...)
	}
	if !lValues.IsNil() {
		members = append(members, lValues.Strings()...)
	}

	for _, member := range members {
		uId, pId, ok := libtracker.PeerParseZsetMember(member)
		if !ok {
			continue
		}
		// 1. Delete peer detail
		peerDetailKey := service.SysCache().KeyTrackerPeerDetail(ctx, torrentId, pId)
		_, _ = g.Redis().Do(ctx, "DEL", peerDetailKey)

		// 2. Remove from user sets
		userVal := libtracker.PeerFormatUserSetValue(torrentId, pId)
		userSeedingKey := service.SysCache().KeyTrackerUserSeeding(ctx, uId)
		userLeechingKey := service.SysCache().KeyTrackerUserLeeching(ctx, uId)
		_, _ = g.Redis().Do(ctx, "SREM", userSeedingKey, userVal)
		_, _ = g.Redis().Do(ctx, "SREM", userLeechingKey, userVal)

		// 3. Check and clean seeding users list
		seedingCard, _ := g.Redis().Do(ctx, "SCARD", userSeedingKey)
		if seedingCard.Int() == 0 {
			seedingUsersKey := service.SysCache().KeyTrackerSeedingUsers(ctx)
			_, _ = g.Redis().Do(ctx, "SREM", seedingUsersKey, uId)
		}
	}

	// 4. Delete torrent sets
	_, _ = g.Redis().Do(ctx, "DEL", seedersKey, leechersKey)

	// 5. Remove from active torrents
	activeTorrentsKey := service.SysCache().KeyTrackerActiveTorrents(ctx)
	_, _ = g.Redis().Do(ctx, "SREM", activeTorrentsKey, torrentId)

	return nil
}
