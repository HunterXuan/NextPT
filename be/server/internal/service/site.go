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
)

var (
	localSiteAuditDomain  ISiteAuditDomain
	localSiteAuditUsecase ISiteAuditUsecase
	localSiteConfigDomain ISiteConfigDomain
)

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
