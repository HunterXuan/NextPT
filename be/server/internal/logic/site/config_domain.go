package site

import (
	"context"
	"math"
	"regexp"
	"strconv"
	"strings"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/sitein"
	"server/internal/model/out/siteout"
	"server/internal/service"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
)

const siteConfigValueField = "val"

type sSiteConfigDomain struct{}

func NewSiteConfigDomain() *sSiteConfigDomain {
	return &sSiteConfigDomain{}
}

func init() {
	service.RegisterSiteConfigDomain(NewSiteConfigDomain())
}

// Get 获取后台业务配置项，支持传入默认值兜底（直接查库，无缓存）
func (s *sSiteConfigDomain) Get(ctx context.Context, group, key string, def ...any) *gvar.Var {
	cfg, err := s.getConfigByGroupAndKey(ctx, group, key)

	if err != nil {
		return gvar.New(s.defaultConfigValue(def...))
	}
	return gvar.New(s.getConfigValue(cfg, def...))
}

func (s *sSiteConfigDomain) getConfigByGroupAndKey(ctx context.Context, group, key string) (*entity.SiteConfig, error) {
	var cfg *entity.SiteConfig
	err := dao.SiteConfig.Ctx(ctx).
		Where(dao.SiteConfig.Columns().Group, group).
		Where(dao.SiteConfig.Columns().Key, key).
		Scan(&cfg)
	return cfg, err
}

func (s *sSiteConfigDomain) getConfigValue(cfg *entity.SiteConfig, def ...any) any {
	if cfg == nil || cfg.Value == nil || cfg.Value.IsNil() {
		return s.defaultConfigValue(def...)
	}
	if value := cfg.Value.Get(siteConfigValueField); value != nil {
		return value.Val()
	}
	return s.defaultConfigValue(def...)
}

func (s *sSiteConfigDomain) defaultConfigValue(def ...any) any {
	if len(def) > 0 {
		return def[0]
	}
	return nil
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

func (s *sSiteConfigDomain) AdminListConfigs(ctx context.Context, in sitein.SiteConfigListInp) (*siteout.SiteConfigListOut, error) {
	m := dao.SiteConfig.Ctx(ctx)
	if in.Group != "" {
		m = m.Where(dao.SiteConfig.Columns().Group, in.Group)
	}
	var configs []*entity.SiteConfig
	err := m.Scan(&configs)
	if err != nil {
		return nil, err
	}

	items := make([]*siteout.SiteConfigItem, 0, len(configs))
	for _, cfg := range configs {
		if cfg == nil {
			continue
		}
		items = append(items, s.buildSiteConfigItem(cfg))
	}
	return &siteout.SiteConfigListOut{Configs: items}, nil
}

func (s *sSiteConfigDomain) buildSiteConfigItem(cfg *entity.SiteConfig) *siteout.SiteConfigItem {
	return &siteout.SiteConfigItem{
		Id:        cfg.Id,
		Group:     cfg.Group,
		Key:       cfg.Key,
		Value:     s.getConfigValue(cfg),
		ValueType: string(s.getConfigValueType(cfg.Group, cfg.Key)),
		CreatedAt: cfg.CreatedAt,
		UpdatedAt: cfg.UpdatedAt,
	}
}

func (s *sSiteConfigDomain) AdminUpdateConfig(ctx context.Context, in sitein.SiteConfigUpdateInp) error {
	normalizedValue, err := s.normalizeConfigValue(in.Group, in.Key, in.Value)
	if err != nil {
		return err
	}
	if in.Group+"."+in.Key == consts.SiteConfigSiteTasks {
		if err := s.validateTaskRoleLevels(ctx, normalizedValue.(model.SiteTasks)); err != nil {
			return err
		}
	}
	valueJson := s.encodeConfigValue(normalizedValue)
	_, err = dao.SiteConfig.Ctx(ctx).
		Data(dao.SiteConfig.Columns().Value, valueJson).
		Where(dao.SiteConfig.Columns().Group, in.Group).
		Where(dao.SiteConfig.Columns().Key, in.Key).
		Update()
	if err != nil {
		return err
	}
	cacheKey := service.SysCache().KeySiteConfigFullPath(ctx, in.Group+"."+in.Key)
	_, _ = gcache.Remove(ctx, cacheKey)
	_ = service.SysCache().PublishInvalidate(ctx, cacheKey)
	return nil
}

func (s *sSiteConfigDomain) encodeConfigValue(value any) string {
	return gjson.MustEncodeString(map[string]any{
		siteConfigValueField: value,
	})
}

func (s *sSiteConfigDomain) getConfigValueType(group, key string) consts.SiteConfigValueType {
	value, ok := consts.SiteConfigDefaults[group+"."+key]
	if !ok {
		return consts.SiteConfigValueTypeJSON
	}

	switch value.(type) {
	case bool:
		return consts.SiteConfigValueTypeBoolean
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return consts.SiteConfigValueTypeInt
	case float32, float64:
		return consts.SiteConfigValueTypeFloat
	case string:
		return consts.SiteConfigValueTypeString
	default:
		return consts.SiteConfigValueTypeJSON
	}
}

func (s *sSiteConfigDomain) normalizeConfigValue(group, key string, value any) (any, error) {
	fullKey := group + "." + key
	if fullKey == consts.SiteConfigIamInviteBypassEmailPattern {
		pattern := strings.TrimSpace(gvar.New(value).String())
		if pattern == "" {
			return "", nil
		}
		if _, err := regexp.Compile(pattern); err != nil {
			return nil, gerror.New("invalid invite registration email pattern")
		}
		return pattern, nil
	}

	if fullKey == consts.SiteConfigSiteTasks {
		var tasks model.SiteTasks
		if err := gvar.New(value).Scan(&tasks); err != nil {
			return nil, gerror.New("invalid site tasks config")
		}
		tasks = tasks.Normalized()
		if err := tasks.Validate(); err != nil {
			return nil, err
		}
		return tasks, nil
	}

	if fullKey == consts.SiteConfigSiteAdvertisements {
		var advertisements model.SiteAdvertisements
		if err := gvar.New(value).Scan(&advertisements); err != nil {
			return nil, gerror.New("invalid site advertisements config")
		}
		if err := advertisements.Validate(); err != nil {
			return nil, err
		}
		advertisements = advertisements.Normalized()
		if err := advertisements.Validate(); err != nil {
			return nil, err
		}
		return advertisements, nil
	}

	switch s.getConfigValueType(group, key) {
	case consts.SiteConfigValueTypeBoolean:
		return s.normalizeConfigBool(value)
	case consts.SiteConfigValueTypeInt:
		return s.normalizeConfigInt(value)
	case consts.SiteConfigValueTypeFloat:
		return s.normalizeConfigFloat(value)
	default:
		return value, nil
	}
}

func (s *sSiteConfigDomain) validateTaskRoleLevels(ctx context.Context, tasks model.SiteTasks) error {
	roles, err := service.IamRoleDomain().ListRoles(ctx)
	if err != nil {
		return err
	}
	normalLevels := make(map[uint64]struct{}, len(roles))
	for _, role := range roles {
		if !role.IsStaff {
			normalLevels[uint64(role.Level)] = struct{}{}
		}
	}
	for _, task := range tasks {
		if task.Rule.Type != consts.SiteTaskRuleTypeRoleLevelReached {
			continue
		}
		if _, exists := normalLevels[task.Rule.Target]; !exists {
			return gerror.New("task role level must belong to a non-staff role")
		}
	}
	return nil
}

func (s *sSiteConfigDomain) normalizeConfigBool(value any) (bool, error) {
	valueVar := s.configValueVar(value)
	if _, err := strconv.ParseBool(valueVar.String()); err != nil {
		return false, gerror.New("invalid boolean config value")
	}
	return valueVar.Bool(), nil
}

func (s *sSiteConfigDomain) normalizeConfigInt(value any) (int64, error) {
	valueVar := s.configValueVar(value)
	if _, err := strconv.ParseInt(valueVar.String(), 10, 64); err != nil {
		return 0, gerror.New("invalid integer config value")
	}
	return valueVar.Int64(), nil
}

func (s *sSiteConfigDomain) normalizeConfigFloat(value any) (float64, error) {
	valueVar, _, err := s.configNumberVar(value, "invalid float config value")
	if err != nil {
		return 0, err
	}
	return valueVar.Float64(), nil
}

func (s *sSiteConfigDomain) configValueVar(value any) *gvar.Var {
	if text, ok := value.(string); ok {
		return gvar.New(strings.TrimSpace(text))
	}
	return gvar.New(value)
}

func (s *sSiteConfigDomain) configNumberVar(value any, errorMessage string) (*gvar.Var, float64, error) {
	valueVar := s.configValueVar(value)
	number, err := strconv.ParseFloat(valueVar.String(), 64)
	if err != nil || math.IsNaN(number) || math.IsInf(number, 0) {
		return nil, 0, gerror.New(errorMessage)
	}
	return valueVar, number, nil
}
