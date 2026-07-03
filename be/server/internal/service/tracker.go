// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	v1 "server/api/tracker/v1"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/trackerin"
	"server/internal/model/out/trackerout"

	"github.com/gogf/gf/v2/os/gtime"
)

type (
	ITrackerEventDomain interface {
		// DeleteExpiredIdempotency 删除过期的幂等记录
		DeleteExpiredIdempotency(ctx context.Context, beforeTime *gtime.Time) (int64, error)
	}
	ITrackerEventUsecase interface {
		// PushAnnounceEvent 接收并投递 Announce 事件到队列 (Redis Stream)
		PushAnnounceEvent(event *trackerin.AnnounceEvent) error
	}
	ITrackerPeerDomain interface {
		// GetPeers 获取种子当前在线的全部做种者和下载者
		GetPeers(ctx context.Context, torrentId uint64) ([]*entity.TrackerPeer, error)
		GetSeedingUsers(ctx context.Context) ([]uint64, error)
		GetUserSeedingPeers(ctx context.Context, userId uint64) ([]entity.TrackerPeer, error)
		GetUserLeechingPeers(ctx context.Context, userId uint64) ([]entity.TrackerPeer, error)
		// GetEnabledClientWhitelists 获取所有启用的客户端白名单规则
		GetEnabledClientWhitelists(ctx context.Context) ([]*entity.TrackerAgentWhitelist, error)
		// calculateTrafficDiff 计算本次汇报的上传下载增量
		CalculateTrafficDiff(ctx context.Context, event *trackerin.AnnounceEvent, oldPeer *entity.TrackerPeer) (diffUp int64, diffDn int64)
		// 辅助方法：维护 Redis Peer 缓存与 ZSET 索引以及用户做种下载索引
		UpsertPeer(ctx context.Context, event *trackerin.AnnounceEvent, oldPeer *entity.TrackerPeer, interval int) error
		// RemovePeer 清理指定的 Peer 及所有相关索引
		RemovePeer(ctx context.Context, torrentId uint64, userId uint64, peerId string) error
		DeletePeersByTorrentId(ctx context.Context, torrentId uint64) error
	}
	ITrackerPeerUsecase interface {
		Announce(ctx context.Context, actor *model.Actor, req *v1.AnnounceReq) (*trackerout.AnnounceOut, error)
		// Scrape 批量查询种子的做种、下载和完成数
		Scrape(ctx context.Context, actor *model.Actor, infoHashes []string) (*trackerout.ScrapeOut, error)
		// CheckClientWhitelist 检查客户端白名单 (内存全量缓存)
		CheckClientWhitelist(ctx context.Context, peerId string, userAgent string) bool
	}
	ITrackerSyncUsecase interface {
		// CleanupGhostPeers 幽灵节点清理
		CleanupGhostPeers(ctx context.Context) error
		// SyncTorrentData 种子数据对齐
		SyncTorrentData(ctx context.Context) error
		// CleanupIdempotency 清理 Tracker 幂等表，保留最近 2 天
		CleanupIdempotency(ctx context.Context) error
	}
)

var (
	localTrackerEventDomain  ITrackerEventDomain
	localTrackerEventUsecase ITrackerEventUsecase
	localTrackerPeerDomain   ITrackerPeerDomain
	localTrackerPeerUsecase  ITrackerPeerUsecase
	localTrackerSyncUsecase  ITrackerSyncUsecase
)

func TrackerEventDomain() ITrackerEventDomain {
	if localTrackerEventDomain == nil {
		panic("implement not found for interface ITrackerEventDomain, forgot register?")
	}
	return localTrackerEventDomain
}

func RegisterTrackerEventDomain(i ITrackerEventDomain) {
	localTrackerEventDomain = i
}

func TrackerEventUsecase() ITrackerEventUsecase {
	if localTrackerEventUsecase == nil {
		panic("implement not found for interface ITrackerEventUsecase, forgot register?")
	}
	return localTrackerEventUsecase
}

func RegisterTrackerEventUsecase(i ITrackerEventUsecase) {
	localTrackerEventUsecase = i
}

func TrackerPeerDomain() ITrackerPeerDomain {
	if localTrackerPeerDomain == nil {
		panic("implement not found for interface ITrackerPeerDomain, forgot register?")
	}
	return localTrackerPeerDomain
}

func RegisterTrackerPeerDomain(i ITrackerPeerDomain) {
	localTrackerPeerDomain = i
}

func TrackerPeerUsecase() ITrackerPeerUsecase {
	if localTrackerPeerUsecase == nil {
		panic("implement not found for interface ITrackerPeerUsecase, forgot register?")
	}
	return localTrackerPeerUsecase
}

func RegisterTrackerPeerUsecase(i ITrackerPeerUsecase) {
	localTrackerPeerUsecase = i
}

func TrackerSyncUsecase() ITrackerSyncUsecase {
	if localTrackerSyncUsecase == nil {
		panic("implement not found for interface ITrackerSyncUsecase, forgot register?")
	}
	return localTrackerSyncUsecase
}

func RegisterTrackerSyncUsecase(i ITrackerSyncUsecase) {
	localTrackerSyncUsecase = i
}
