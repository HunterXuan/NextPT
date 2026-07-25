// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/sitein"
	"server/internal/model/out/siteout"

	"github.com/gogf/gf/v2/container/gvar"
)

type (
	ISiteAnnouncementDomain interface {
		ListPublished(ctx context.Context, userId uint64, in sitein.AnnouncementListInp) ([]entity.SiteAnnouncement, int, error)
		QueryReadAnnouncementIdsByUser(ctx context.Context, userId uint64) ([]uint64, error)
		QueryReadIds(ctx context.Context, userId uint64, announcementIds []uint64) (map[uint64]bool, error)
		MarkRead(ctx context.Context, userId uint64, announcementId uint64) error
		AdminList(ctx context.Context, in sitein.AdminAnnouncementListInp) ([]entity.SiteAnnouncement, int, error)
		AdminCreate(ctx context.Context, in sitein.AdminAnnouncementCreateInp, actorId uint64) (uint64, error)
		AdminUpdate(ctx context.Context, in sitein.AdminAnnouncementUpdateInp, actorId uint64) error
		AdminDelete(ctx context.Context, id uint64) error
		DeleteReadRecordsByAnnouncementId(ctx context.Context, announcementId uint64) error
	}
	ISiteAnnouncementUsecase interface {
		List(ctx context.Context, actor *model.Actor, in sitein.AnnouncementListInp) (*siteout.AnnouncementListOut, error)
		MarkRead(ctx context.Context, actor *model.Actor, in sitein.AnnouncementReadInp) error
		AdminList(ctx context.Context, actor *model.Actor, in sitein.AdminAnnouncementListInp) (*siteout.AnnouncementListOut, error)
		AdminCreate(ctx context.Context, actor *model.Actor, in sitein.AdminAnnouncementCreateInp) (*siteout.AnnouncementCreateOut, error)
		AdminUpdate(ctx context.Context, actor *model.Actor, in sitein.AdminAnnouncementUpdateInp) error
		AdminDelete(ctx context.Context, actor *model.Actor, in sitein.AdminAnnouncementDeleteInp) error
	}
	ISiteAuditDomain interface {
		AdminListAudits(ctx context.Context, in sitein.AuditListInp) ([]entity.SiteAudit, int, error)
		Create(ctx context.Context, in sitein.AuditCreateInp) error
	}
	ISiteAuditUsecase interface {
		Record(ctx context.Context, actor *model.Actor, in sitein.AuditRecordInp)
	}
	ISiteConfigDomain interface {
		// Get 获取后台业务配置项，支持传入默认值兜底（直接查库，无缓存）
		Get(ctx context.Context, group string, key string, def ...any) *gvar.Var
		// GetByPath 模仿 GF 官方 gcfg 行为，通过 "group.key" 的格式获取后台业务配置项
		// 如果没有提供 customDef，底层自动从 consts.SiteConfigDefaults 获取兜底默认值（直接查库，无缓存）
		GetByPath(ctx context.Context, path string, customDef ...any) *gvar.Var
		AdminListConfigs(ctx context.Context, in sitein.SiteConfigListInp) (*siteout.SiteConfigListOut, error)
		AdminUpdateConfig(ctx context.Context, in sitein.SiteConfigUpdateInp) error
	}
	ISiteMessageDomain interface {
		ListByReceiver(ctx context.Context, receiverId uint64, in sitein.MessageListInp) ([]entity.SiteMessage, int, error)
		MarkRead(ctx context.Context, receiverId uint64, id uint64) error
		MarkAllRead(ctx context.Context, receiverId uint64) error
		Create(ctx context.Context, in sitein.MessageCreateInp) (uint64, error)
		BatchCreate(ctx context.Context, items []sitein.MessageCreateInp) error
		AdminList(ctx context.Context, in sitein.AdminMessageListInp) ([]entity.SiteMessage, int, error)
	}
	ISiteMessageUsecase interface {
		List(ctx context.Context, actor *model.Actor, in sitein.MessageListInp) (*siteout.MessageListOut, error)
		MarkRead(ctx context.Context, actor *model.Actor, in sitein.MessageReadInp) error
		MarkAllRead(ctx context.Context, actor *model.Actor, in sitein.MessageReadAllInp) error
		Notify(ctx context.Context, in sitein.MessageNotifyInp)
		AdminList(ctx context.Context, actor *model.Actor, in sitein.AdminMessageListInp) (*siteout.MessageListOut, error)
		AdminCreate(ctx context.Context, actor *model.Actor, in sitein.AdminMessageCreateInp) (*siteout.MessageCreateOut, error)
	}
)

var (
	localSiteAnnouncementDomain  ISiteAnnouncementDomain
	localSiteAnnouncementUsecase ISiteAnnouncementUsecase
	localSiteAuditDomain         ISiteAuditDomain
	localSiteAuditUsecase        ISiteAuditUsecase
	localSiteConfigDomain        ISiteConfigDomain
	localSiteMessageDomain       ISiteMessageDomain
	localSiteMessageUsecase      ISiteMessageUsecase
)

func SiteAnnouncementDomain() ISiteAnnouncementDomain {
	if localSiteAnnouncementDomain == nil {
		panic("implement not found for interface ISiteAnnouncementDomain, forgot register?")
	}
	return localSiteAnnouncementDomain
}

func RegisterSiteAnnouncementDomain(i ISiteAnnouncementDomain) {
	localSiteAnnouncementDomain = i
}

func SiteAnnouncementUsecase() ISiteAnnouncementUsecase {
	if localSiteAnnouncementUsecase == nil {
		panic("implement not found for interface ISiteAnnouncementUsecase, forgot register?")
	}
	return localSiteAnnouncementUsecase
}

func RegisterSiteAnnouncementUsecase(i ISiteAnnouncementUsecase) {
	localSiteAnnouncementUsecase = i
}

func SiteAuditDomain() ISiteAuditDomain {
	if localSiteAuditDomain == nil {
		panic("implement not found for interface ISiteAuditDomain, forgot register?")
	}
	return localSiteAuditDomain
}

func RegisterSiteAuditDomain(i ISiteAuditDomain) {
	localSiteAuditDomain = i
}

func SiteAuditUsecase() ISiteAuditUsecase {
	if localSiteAuditUsecase == nil {
		panic("implement not found for interface ISiteAuditUsecase, forgot register?")
	}
	return localSiteAuditUsecase
}

func RegisterSiteAuditUsecase(i ISiteAuditUsecase) {
	localSiteAuditUsecase = i
}

func SiteConfigDomain() ISiteConfigDomain {
	if localSiteConfigDomain == nil {
		panic("implement not found for interface ISiteConfigDomain, forgot register?")
	}
	return localSiteConfigDomain
}

func RegisterSiteConfigDomain(i ISiteConfigDomain) {
	localSiteConfigDomain = i
}

func SiteMessageDomain() ISiteMessageDomain {
	if localSiteMessageDomain == nil {
		panic("implement not found for interface ISiteMessageDomain, forgot register?")
	}
	return localSiteMessageDomain
}

func RegisterSiteMessageDomain(i ISiteMessageDomain) {
	localSiteMessageDomain = i
}

func SiteMessageUsecase() ISiteMessageUsecase {
	if localSiteMessageUsecase == nil {
		panic("implement not found for interface ISiteMessageUsecase, forgot register?")
	}
	return localSiteMessageUsecase
}

func RegisterSiteMessageUsecase(i ISiteMessageUsecase) {
	localSiteMessageUsecase = i
}
