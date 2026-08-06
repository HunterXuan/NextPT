// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"server/internal/model/entity"
)

type (
	ISysCache interface {
		// PublishInvalidate 广播缓存失效事件
		PublishInvalidate(ctx context.Context, key string) error
		KeySiteConfig(ctx context.Context, group string, key string) string
		KeySiteConfigFullPath(ctx context.Context, key string) string
		KeyTrackerPeers(ctx context.Context, torrentId uint64) string
		KeyTrackerClientWhitelist(ctx context.Context) string
		KeySysCronLock(ctx context.Context, lockName string) string
		KeyTrackerPeerDetail(ctx context.Context, torrentId uint64, peerId string) string
		KeyTrackerTorrentSeeders(ctx context.Context, torrentId uint64) string
		KeyTrackerTorrentLeechers(ctx context.Context, torrentId uint64) string
		KeyIamUserAcls(ctx context.Context, userId uint64) string
		KeyIamUserPublic(ctx context.Context, userId uint64, language string) string
		KeyIamEmailVerificationToken(ctx context.Context, tokenHash string) string
		KeyIamEmailVerificationUser(ctx context.Context, userId uint64) string
		KeyIamEmailVerificationRateIp(ctx context.Context, ipHash string) string
		KeyIamEmailVerificationRateEmail(ctx context.Context, emailHash string) string
		KeyIamPasswordResetToken(ctx context.Context, tokenHash string) string
		KeyIamPasswordResetUser(ctx context.Context, userId uint64) string
		KeyIamPasswordResetRateIp(ctx context.Context, ipHash string) string
		KeyIamPasswordResetRateEmail(ctx context.Context, emailHash string) string
		KeyIamTwoStepSetup(ctx context.Context, userId uint64) string
		KeyIamTwoStepLoginChallenge(ctx context.Context, challengeHash string) string
		KeyIamTwoStepLoginRate(ctx context.Context, userId uint64) string
		KeyIamSession(ctx context.Context, sessionId string) string
		KeyIamUserSessions(ctx context.Context, userId uint64) string
		KeyIamRolePerms(ctx context.Context, roleId uint) string
		KeyIamRoleActorVersion(ctx context.Context, roleId uint) string
		KeyIamActor(ctx context.Context, userId uint64) string
		KeyTrackerUserSeeding(ctx context.Context, userId uint64) string
		KeyTrackerUserLeeching(ctx context.Context, userId uint64) string
		KeyTrackerActiveTorrents(ctx context.Context) string
		KeyTrackerSeedingUsers(ctx context.Context) string
		KeyIamPasskeyActor(ctx context.Context, passkey string) string
		KeyCatalogTorrentInfoHash(ctx context.Context, hexInfoHash string) string
		KeyCatalogMetadataTmdb(ctx context.Context, tmdbType string, tmdbId string, locale string) string
		KeyCatalogMetadataImdb(ctx context.Context, imdbId string) string
		KeyCatalogMetadataDouban(ctx context.Context, doubanId string) string
		KeyCatalogMetadataBangumi(ctx context.Context, bangumiId string) string
		KeyCatalogHotTorrents(ctx context.Context) string
		KeyForumHotTopics(ctx context.Context) string
		KeyTrackerLockPeer(ctx context.Context, torrentId uint64, userId uint64, peerId string) string
		KeyTrackerAnnounceQueue(ctx context.Context) string
		KeyTrackerAnnounceDedup(ctx context.Context, digest string) string
		KeyTrackerDlqCounts(ctx context.Context) string
		KeyTrackerAnnounceQueueDlq(ctx context.Context) string
	}
	ISysCron interface {
		// Start 启动所有系统定时任务，在全局 Init 中调用
		Start(ctx context.Context)
		AdminListCronLogs(ctx context.Context, jobName string, status int, page int, size int) ([]*entity.SysCronLog, int, error)
	}
	ISysMailgun interface {
		SendTextMail(ctx context.Context, subject string, body string, recipient string) error
		SendHtmlMail(ctx context.Context, subject string, textBody string, htmlBody string, recipient string) error
	}
	ISysStorage interface {
		// Upload 将数据上传到存储引擎
		Upload(ctx context.Context, key string, data []byte, contentType string) error
		// Download 从存储引擎下载文件并保存到本地指定路径
		Download(ctx context.Context, key string, destPath string) error
		// Delete 从存储引擎删除文件
		Delete(ctx context.Context, key string) error
		// GetLocalPath 获取文件的本地访问路径。
		// 如果是 local 驱动，返回文件在磁盘的真实绝对路径；
		// 如果是 s3 等远程驱动，返回本地缓存目录中的绝对路径，供调用方检查是否存在并按需调用 Download。
		GetLocalPath(ctx context.Context, key string) string
	}
)

var (
	localSysCache   ISysCache
	localSysCron    ISysCron
	localSysMailgun ISysMailgun
	localSysStorage ISysStorage
)

func SysCache() ISysCache {
	if localSysCache == nil {
		panic("implement not found for interface ISysCache, forgot register?")
	}
	return localSysCache
}

func RegisterSysCache(i ISysCache) {
	localSysCache = i
}

func SysCron() ISysCron {
	if localSysCron == nil {
		panic("implement not found for interface ISysCron, forgot register?")
	}
	return localSysCron
}

func RegisterSysCron(i ISysCron) {
	localSysCron = i
}

func SysMailgun() ISysMailgun {
	if localSysMailgun == nil {
		panic("implement not found for interface ISysMailgun, forgot register?")
	}
	return localSysMailgun
}

func RegisterSysMailgun(i ISysMailgun) {
	localSysMailgun = i
}

func SysStorage() ISysStorage {
	if localSysStorage == nil {
		panic("implement not found for interface ISysStorage, forgot register?")
	}
	return localSysStorage
}

func RegisterSysStorage(i ISysStorage) {
	localSysStorage = i
}
