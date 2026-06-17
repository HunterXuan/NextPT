package sys

import (
	"context"
	"fmt"
	"time"

	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

type sSysCache struct {
	prefix                      string
	channelCacheInvalidate      string
	prefixSiteConfig            string
	prefixTrackerPeers          string
	prefixTrackerClientWL       string
	prefixSysCronLock           string
	prefixTrackerTorrent        string
	prefixTrackerUser           string
	prefixTrackerActiveTorrents string
	prefixTrackerSeedingUsers   string
	prefixIamPasskey            string
	prefixCatalogTorrentHash    string
	prefixTrackerLock           string
	prefixTrackerAnnounceQueue  string
	prefixTrackerAnnounceDedup  string
}

func NewSysCache() *sSysCache {
	ctx := gctx.New()
	p := g.Cfg().MustGet(ctx, "cache.prefix").String()
	if p != "" && p[len(p)-1] != ':' {
		p += ":"
	}

	return &sSysCache{
		prefix:                      p,
		channelCacheInvalidate:      p + "sys:cache:invalidate",
		prefixSiteConfig:            p + "site:config:",
		prefixTrackerPeers:          p + "tracker:peers:torrent:",
		prefixTrackerClientWL:       p + "tracker:client_wl:compiled",
		prefixSysCronLock:           p + "sys:cron:lock:",
		prefixTrackerTorrent:        p + "tracker:torrent:",
		prefixTrackerUser:           p + "tracker:user:",
		prefixTrackerActiveTorrents: p + "tracker:active_torrents",
		prefixTrackerSeedingUsers:   p + "tracker:seeding_users",
		prefixIamPasskey:            p + "iam:passkey:",
		prefixCatalogTorrentHash:    p + "catalog:torrent:hash:",
		prefixTrackerLock:           p + "tracker:lock:",
		prefixTrackerAnnounceQueue:  p + "tracker:announce_queue",
		prefixTrackerAnnounceDedup:  p + "tracker:announce_dedup:",
	}
}

func init() {
	cache := NewSysCache()
	service.RegisterSysCache(cache)

	// 启动 Redis PubSub 监听，用于多实例下的本地缓存一致性失效
	go cache.listenCacheInvalidation()
}

// listenCacheInvalidation 后台协程订阅 Redis
func (s *sSysCache) listenCacheInvalidation() {
	ctx := gctx.New()

	for {
		// 每次外层循环都会尝试获取新的 Redis 连接
		conn, err := g.Redis().Conn(ctx)
		if err != nil {
			g.Log().Warningf(ctx, "sysCache: wait for redis ready: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		_, err = conn.Subscribe(ctx, s.channelCacheInvalidate)
		if err != nil {
			g.Log().Warningf(ctx, "sysCache: subscribe to %s failed: %v", s.channelCacheInvalidate, err)
			conn.Close(ctx)
			time.Sleep(3 * time.Second)
			continue
		}

		g.Log().Infof(ctx, "sysCache: successfully subscribed to %s for local cache invalidation", s.channelCacheInvalidate)

		// 阻塞读取消息
		for {
			msg, err := conn.ReceiveMessage(ctx)
			if err != nil {
				g.Log().Errorf(ctx, "sysCache: receive message error: %v, will reconnect...", err)
				break
			}

			// 收到了失效通知，清理本地内存缓存
			if msg != nil && msg.Payload != "" {
				keyToInvalidate := msg.Payload
				_, _ = gcache.Remove(ctx, keyToInvalidate)
				g.Log().Debugf(ctx, "sysCache: local cache invalidated for key [%s]", keyToInvalidate)
			}
		}

		conn.Close(ctx)
		time.Sleep(1 * time.Second) // 断开后略作延迟重连
	}
}

// PublishInvalidate 广播缓存失效事件
func (s *sSysCache) PublishInvalidate(ctx context.Context, key string) error {
	_, err := g.Redis().Publish(ctx, s.channelCacheInvalidate, key)
	if err != nil {
		g.Log().Errorf(ctx, "sysCache: failed to publish cache invalidation for key [%s]: %v", key, err)
		return err
	}
	return nil
}

func (s *sSysCache) KeySiteConfig(ctx context.Context, group, key string) string {
	return s.prefixSiteConfig + group + ":" + key
}

func (s *sSysCache) KeySiteConfigFullPath(ctx context.Context, key string) string {
	return s.prefixSiteConfig + key
}

func (s *sSysCache) KeyTrackerPeers(ctx context.Context, torrentId uint64) string {
	return fmt.Sprintf("%s%d", s.prefixTrackerPeers, torrentId)
}

func (s *sSysCache) KeyTrackerClientWhitelist(ctx context.Context) string {
	return s.prefixTrackerClientWL
}

func (s *sSysCache) KeySysCronLock(ctx context.Context, lockName string) string {
	return s.prefixSysCronLock + lockName
}

func (s *sSysCache) KeyTrackerPeerDetail(ctx context.Context, torrentId uint64, peerId string) string {
	return fmt.Sprintf("%s%d:peer:%s", s.prefixTrackerTorrent, torrentId, peerId)
}

func (s *sSysCache) KeyTrackerTorrentSeeders(ctx context.Context, torrentId uint64) string {
	return fmt.Sprintf("%s%d:seeders", s.prefixTrackerTorrent, torrentId)
}

func (s *sSysCache) KeyTrackerTorrentLeechers(ctx context.Context, torrentId uint64) string {
	return fmt.Sprintf("%s%d:leechers", s.prefixTrackerTorrent, torrentId)
}

func (s *sSysCache) KeyIamUserAcls(ctx context.Context, userId uint64) string {
	return fmt.Sprintf("iam:user:%d:acls", userId)
}

func (s *sSysCache) KeyIamRolePerms(ctx context.Context, roleId uint) string {
	return fmt.Sprintf("iam:role:%d:perms", roleId)
}

func (s *sSysCache) KeyIamRoleActorVersion(ctx context.Context, roleId uint) string {
	return fmt.Sprintf("iam:role:%d:actor_version", roleId)
}

func (s *sSysCache) KeyIamActor(ctx context.Context, userId uint64) string {
	return fmt.Sprintf("iam:actor:%d", userId)
}

func (s *sSysCache) KeyTrackerUserSeeding(ctx context.Context, userId uint64) string {
	return fmt.Sprintf("%s%d:seeding", s.prefixTrackerUser, userId)
}

func (s *sSysCache) KeyTrackerUserLeeching(ctx context.Context, userId uint64) string {
	return fmt.Sprintf("%s%d:leeching", s.prefixTrackerUser, userId)
}

func (s *sSysCache) KeyTrackerActiveTorrents(ctx context.Context) string {
	return s.prefixTrackerActiveTorrents
}

func (s *sSysCache) KeyTrackerSeedingUsers(ctx context.Context) string {
	return s.prefixTrackerSeedingUsers
}

func (s *sSysCache) KeyIamPasskeyActor(ctx context.Context, passkey string) string {
	return s.prefixIamPasskey + passkey
}

func (s *sSysCache) KeyCatalogTorrentInfoHash(ctx context.Context, hexInfoHash string) string {
	return s.prefixCatalogTorrentHash + hexInfoHash
}

func (s *sSysCache) KeyTrackerLockPeer(ctx context.Context, torrentId, userId uint64, peerId string) string {
	return fmt.Sprintf("%speer:%d:%d:%s", s.prefixTrackerLock, torrentId, userId, peerId)
}

func (s *sSysCache) KeyTrackerAnnounceQueue(ctx context.Context) string {
	return s.prefixTrackerAnnounceQueue
}

func (s *sSysCache) KeyTrackerAnnounceDedup(ctx context.Context, digest string) string {
	return s.prefixTrackerAnnounceDedup + digest
}

func (s *sSysCache) KeyTrackerDlqCounts(ctx context.Context) string {
	return s.prefix + "tracker:dlq_counts"
}

func (s *sSysCache) KeyTrackerAnnounceQueueDlq(ctx context.Context) string {
	return s.prefixTrackerAnnounceQueue + "_dlq"
}
