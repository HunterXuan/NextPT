package admin

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/adminout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/util/gconv"
)

type sAdminSiteConfigUsecase struct{}

func init() {
	service.RegisterAdminSiteConfigUsecase(NewAdminSiteConfigUsecase())
}

func NewAdminSiteConfigUsecase() *sAdminSiteConfigUsecase {
	return &sAdminSiteConfigUsecase{}
}

func (s *sAdminSiteConfigUsecase) List(ctx context.Context, actor *model.Actor, in adminin.SiteConfigListInp) (*adminout.SiteConfigListOut, error) {
	out, err := service.SiteConfigDomain().AdminListConfigs(ctx, sitein.SiteConfigListInp{Group: in.Group})
	if err != nil {
		return nil, err
	}
	return &adminout.SiteConfigListOut{Configs: out.Configs}, nil
}

func (s *sAdminSiteConfigUsecase) Update(ctx context.Context, actor *model.Actor, in adminin.SiteConfigUpdateInp) error {
	if s.isRoleIdConfig(in.Group, in.Key) {
		role, err := service.IamRoleDomain().GetRoleById(ctx, gconv.Uint(in.Value))
		if err != nil {
			return err
		}
		if role == nil {
			return gerror.New(gi18n.T(ctx, "admin.role.not_found"))
		}
	}
	if err := service.SiteConfigDomain().AdminUpdateConfig(ctx, sitein.SiteConfigUpdateInp{
		Group: in.Group,
		Key:   in.Key,
		Value: in.Value,
	}); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeSiteConfig,
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"group": in.Group,
			"key":   in.Key,
		},
	})
	return nil
}

func (s *sAdminSiteConfigUsecase) isRoleIdConfig(group string, key string) bool {
	path := group + "." + key
	return path == consts.SiteConfigIamDefaultRegisterRole
}
