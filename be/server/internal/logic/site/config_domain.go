package site

import (
	"context"
	"strings"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gcache"
)

type sSiteConfigDomain struct{}

func NewSiteConfigDomain() *sSiteConfigDomain {
	return &sSiteConfigDomain{}
}

func init() {
	service.RegisterSiteConfigDomain(NewSiteConfigDomain())
}

func (s *sSiteConfigDomain) GetConfigByGroupAndKey(ctx context.Context, group, key string) (*entity.SiteConfig, error) {
	var cfg *entity.SiteConfig
	err := dao.SiteConfig.Ctx(ctx).
		Where(dao.SiteConfig.Columns().Group, group).
		Where(dao.SiteConfig.Columns().Key, key).
		Scan(&cfg)
	return cfg, err
}

// Get 获取后台业务配置项，支持传入默认值兜底（直接查库，无缓存）
func (s *sSiteConfigDomain) Get(ctx context.Context, group, key string, def ...any) *gvar.Var {
	cfg, err := s.GetConfigByGroupAndKey(ctx, group, key)

	if err != nil || cfg == nil || cfg.Value == nil || cfg.Value.IsNil() {
		if len(def) > 0 {
			return gvar.New(def[0])
		}
		return gvar.New(nil)
	}

	return cfg.Value.Var()
}

// GetByPath 模仿 GF 官方 gcfg 行为，通过 "group.key" 的格式获取后台业务配置项
// 如果没有提供 customDef，底层自动从 consts.SiteConfigDefaults 获取兜底默认值（直接查库，无缓存）
func (s *sSiteConfigDomain) GetByPath(ctx context.Context, path string, customDef ...any) *gvar.Var {
	parts := strings.SplitN(path, ".", 2)

	var def any
	if len(customDef) > 0 {
		def = customDef[0]
	} else {
		def = consts.SiteConfigDefaults[path]
	}

	if len(parts) != 2 {
		return gvar.New(def)
	}
	return s.Get(ctx, parts[0], parts[1], def)
}

func (s *sSiteConfigDomain) AdminListConfigs(ctx context.Context, group string) ([]*entity.SiteConfig, error) {
	m := dao.SiteConfig.Ctx(ctx)
	if group != "" {
		m = m.Where(dao.SiteConfig.Columns().Group, group)
	}
	var configs []*entity.SiteConfig
	err := m.Scan(&configs)
	return configs, err
}

func (s *sSiteConfigDomain) AdminUpdateConfig(ctx context.Context, group string, key string, value string) error {
	_, err := dao.SiteConfig.Ctx(ctx).
		Data(dao.SiteConfig.Columns().Value, value).
		Where(dao.SiteConfig.Columns().Group, group).
		Where(dao.SiteConfig.Columns().Key, key).
		Update()
	if err != nil {
		return err
	}
	cacheKey := service.SysCache().KeySiteConfigFullPath(ctx, group+"."+key)
	_, _ = gcache.Remove(ctx, cacheKey)
	_ = service.SysCache().PublishInvalidate(ctx, cacheKey)
	return nil
}
