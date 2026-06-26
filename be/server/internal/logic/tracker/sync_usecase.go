package tracker

import (
	"context"
	"time"

	libtracker "server/internal/library/tracker"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sTrackerSyncUsecase struct {
}

func init() {
	service.RegisterTrackerSyncUsecase(NewTrackerSyncUsecase())
}

const evictPeerLuaScript = `
local zset_key = KEYS[1]
local detail_key = KEYS[2]
local user_set_key = KEYS[3]
local member = ARGV[1]
local now = tonumber(ARGV[2])
local user_val = ARGV[3]

local score = redis.call('ZSCORE', zset_key, member)
if score and tonumber(score) <= now then
    redis.call('ZREM', zset_key, member)
    redis.call('DEL', detail_key)
    redis.call('SREM', user_set_key, user_val)
    return 1
end
return 0
`

func NewTrackerSyncUsecase() *sTrackerSyncUsecase {
	return &sTrackerSyncUsecase{}
}

// CleanupGhostPeers 幽灵节点清理
func (s *sTrackerSyncUsecase) CleanupGhostPeers(ctx context.Context) error {
	// 从 Redis 获取所有活跃种子 ID
	activeTorrentsKey := service.SysCache().KeyTrackerActiveTorrents(ctx)
	torrentIdsVal, err := g.Redis().Do(ctx, "SMEMBERS", activeTorrentsKey)
	if err != nil {
		glog.Error(ctx, "[Cron] Ghost peer cleanup failed to fetch active torrent IDs:", err)
		return err
	}
	if torrentIdsVal.IsNil() || len(torrentIdsVal.Strings()) == 0 {
		return nil
	}

	now := time.Now().Unix()
	removedCount := 0
	seedingUsersKey := service.SysCache().KeyTrackerSeedingUsers(ctx)

	for _, tVal := range torrentIdsVal.Strings() {
		torrentId := gconv.Uint64(tVal)
		seedersKey := service.SysCache().KeyTrackerTorrentSeeders(ctx, torrentId)
		leechersKey := service.SysCache().KeyTrackerTorrentLeechers(ctx, torrentId)

		// 找出过期的 Peer ID
		expiredSeeders, _ := g.Redis().Do(ctx, "ZRANGEBYSCORE", seedersKey, 0, now)
		expiredLeechers, _ := g.Redis().Do(ctx, "ZRANGEBYSCORE", leechersKey, 0, now)

		affectedUserIds := make(map[uint64]bool)

		cleanupPeers := func(members []string, isSeeder bool) {
			for _, member := range members {
				uId, pId, ok := libtracker.PeerParseZsetMember(member)
				if !ok {
					continue
				}

				peerDetailKey := service.SysCache().KeyTrackerPeerDetail(ctx, torrentId, pId)
				userVal := libtracker.PeerFormatUserSetValue(torrentId, pId)

				var zsetKey string
				var userSetKey string
				if isSeeder {
					zsetKey = seedersKey
					userSetKey = service.SysCache().KeyTrackerUserSeeding(ctx, uId)
				} else {
					zsetKey = leechersKey
					userSetKey = service.SysCache().KeyTrackerUserLeeching(ctx, uId)
				}

				res, err := g.Redis().Do(ctx, "EVAL", evictPeerLuaScript, 3, zsetKey, peerDetailKey, userSetKey, member, now, userVal)
				if err == nil && !res.IsNil() && res.Int() == 1 {
					removedCount++
					if isSeeder {
						affectedUserIds[uId] = true
					}
				}
			}
		}

		if !expiredSeeders.IsNil() && len(expiredSeeders.Strings()) > 0 {
			cleanupPeers(expiredSeeders.Strings(), true)
		}
		if !expiredLeechers.IsNil() && len(expiredLeechers.Strings()) > 0 {
			cleanupPeers(expiredLeechers.Strings(), false)
		}

		// 检查该种子是否还有任何 peer，如果没有，从 active_torrents 中移除
		sCard, _ := g.Redis().Do(ctx, "ZCARD", seedersKey)
		lCard, _ := g.Redis().Do(ctx, "ZCARD", leechersKey)
		if sCard.Int() == 0 && lCard.Int() == 0 {
			_, _ = g.Redis().Do(ctx, "SREM", activeTorrentsKey, torrentId)
		}

		// 检查受影响的做种用户是否还有做种，如果没有，从 seeding_users 中移除
		for uId := range affectedUserIds {
			userSeedingKey := service.SysCache().KeyTrackerUserSeeding(ctx, uId)
			seedingCard, _ := g.Redis().Do(ctx, "SCARD", userSeedingKey)
			if seedingCard.Int() == 0 {
				_, _ = g.Redis().Do(ctx, "SREM", seedingUsersKey, uId)
			}
		}
	}

	if removedCount > 0 {
		glog.Infof(ctx, "[Cron] Ghost peer cleanup completed. Removed %d dead peers.", removedCount)
	}
	return nil
}

// SyncTorrentData 种子数据对齐
func (s *sTrackerSyncUsecase) SyncTorrentData(ctx context.Context) error {
	// 1. 获取 Redis 中的活跃种子 ID
	activeTorrentsKey := service.SysCache().KeyTrackerActiveTorrents(ctx)
	redisTorrentsVal, err := g.Redis().Do(ctx, "SMEMBERS", activeTorrentsKey)
	if err != nil {
		glog.Error(ctx, "[Cron] Sync torrent data failed to fetch active torrents from Redis:", err)
		return err
	}

	activeTorrentIds := make(map[uint64]bool)
	if !redisTorrentsVal.IsNil() {
		for _, v := range redisTorrentsVal.Strings() {
			activeTorrentIds[gconv.Uint64(v)] = true
		}
	}

	// 2. 获取 MySQL 中 seeders > 0 或 leechers > 0 的种子 ID 并合并
	var dbTorrents []entity.CatalogTorrent
	dbTorrents, err = service.CatalogTorrentDomain().QueryActiveTorrentIds(ctx)
	if err != nil {
		glog.Error(ctx, "[Cron] Sync torrent data failed to fetch active torrents from MySQL:", err)
		return err
	}

	for _, t := range dbTorrents {
		activeTorrentIds[t.Id] = true
	}

	if len(activeTorrentIds) == 0 {
		return nil
	}

	now := time.Now().Unix()
	updatedCount := 0

	for torrentId := range activeTorrentIds {
		seedersKey := service.SysCache().KeyTrackerTorrentSeeders(ctx, torrentId)
		leechersKey := service.SysCache().KeyTrackerTorrentLeechers(ctx, torrentId)

		sCountVal, err := g.Redis().Do(ctx, "ZCOUNT", seedersKey, now, "+inf")
		if err != nil {
			continue
		}
		lCountVal, err := g.Redis().Do(ctx, "ZCOUNT", leechersKey, now, "+inf")
		if err != nil {
			continue
		}

		if err = service.CatalogTorrentDomain().UpdateTorrentPeerStats(ctx, torrentId, sCountVal.Int(), lCountVal.Int()); err != nil {
			glog.Errorf(ctx, "[Cron] Sync torrent data update failed for torrent %d: %v", torrentId, err)
			return err
		}
		updatedCount++
	}

	glog.Infof(ctx, "[Cron] Sync torrent data completed for %d torrents.", updatedCount)
	return nil
}

// CleanupIdempotency 清理 Tracker 幂等表，保留最近 2 天
func (s *sTrackerSyncUsecase) CleanupIdempotency(ctx context.Context) error {
	retentionDays := 2
	beforeTime := gtime.Now().Add(-time.Duration(retentionDays*24) * time.Hour)
	rows, err := service.TrackerEventDomain().DeleteExpiredIdempotency(ctx, beforeTime)
	if err != nil {
		glog.Error(ctx, "[Cron] Failed to cleanup tracker event idempotency:", err)
		return err
	}
	if rows > 0 {
		glog.Infof(ctx, "[Cron] Cleanup tracker event idempotency completed. Removed %d expired entries before %s.", rows, beforeTime)
	}
	return nil
}
