package admin

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/adminout"
	"server/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gcache"
)

type sAdminIamRoleUsecase struct{}

func NewAdminIamRoleUsecase() *sAdminIamRoleUsecase {
	return &sAdminIamRoleUsecase{}
}

func init() {
	service.RegisterAdminIamRoleUsecase(NewAdminIamRoleUsecase())
}

func (s *sAdminIamRoleUsecase) List(ctx context.Context, actor *model.Actor) (*adminout.IamRoleListOut, error) {
	roles, err := service.IamRoleDomain().AdminListRoles(ctx)
	if err != nil {
		return nil, err
	}
	return &adminout.IamRoleListOut{Roles: roles}, nil
}

func (s *sAdminIamRoleUsecase) Create(ctx context.Context, actor *model.Actor, in adminin.IamRoleCreateInp) (uint, error) {
	nameI18N, _ := gjson.Encode(in.NameI18N)
	rules, _ := gjson.Encode(in.Rules)
	permissions, _ := gjson.Encode(in.Permissions)

	id, err := service.IamRoleDomain().AdminCreateRole(ctx, in.Level, nameI18N, rules, permissions, in.IsStaff)
	if err != nil {
		return 0, err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionCreate,
		TargetType: consts.SiteAuditTargetTypeIamRole,
		TargetId:   uint64(id),
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"level":   in.Level,
			"isStaff": in.IsStaff,
		},
	})
	return id, nil
}

func (s *sAdminIamRoleUsecase) Update(ctx context.Context, actor *model.Actor, in adminin.IamRoleUpdateInp) error {
	var nameI18N, rules, permissions []byte
	if in.NameI18N != nil {
		nameI18N, _ = gjson.Encode(in.NameI18N)
	}
	if in.Rules != nil {
		rules, _ = gjson.Encode(in.Rules)
	}
	if in.Permissions != nil {
		permissions, _ = gjson.Encode(in.Permissions)
	}

	err := service.IamRoleDomain().AdminUpdateRole(ctx, in.Id, in.Level, nameI18N, rules, permissions, in.IsStaff)
	if err == nil {
		s.invalidateRolePermsCache(ctx, in.Id)

		if in.Level != nil || in.IsStaff != nil {
			if err := s.bumpRoleActorVersion(ctx, in.Id); err != nil {
				return err
			}
		}
		service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
			Action:     consts.SiteAuditActionUpdate,
			TargetType: consts.SiteAuditTargetTypeIamRole,
			TargetId:   uint64(in.Id),
			Level:      consts.SiteAuditLevelCritical,
			Detail: map[string]any{
				"levelChanged":       in.Level != nil,
				"nameChanged":        in.NameI18N != nil,
				"rulesChanged":       in.Rules != nil,
				"permissionsChanged": in.Permissions != nil,
				"isStaffChanged":     in.IsStaff != nil,
			},
		})
	}
	return err
}

func (s *sAdminIamRoleUsecase) invalidateRolePermsCache(ctx context.Context, roleId uint) {
	s.invalidateCacheKey(ctx, service.SysCache().KeyIamRolePerms(ctx, roleId))
}

func (s *sAdminIamRoleUsecase) bumpRoleActorVersion(ctx context.Context, roleId uint) error {
	_, err := g.Redis().Do(ctx, "INCR", service.SysCache().KeyIamRoleActorVersion(ctx, roleId))
	return err
}

func (s *sAdminIamRoleUsecase) invalidateCacheKey(ctx context.Context, key string) {
	_, _ = gcache.Remove(ctx, key)
	_ = service.SysCache().PublishInvalidate(ctx, key)
}

func (s *sAdminIamRoleUsecase) Delete(ctx context.Context, actor *model.Actor, in adminin.IamRoleDeleteInp) error {
	role, err := service.IamRoleDomain().GetRoleById(ctx, in.Id)
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "admin.role.check_users_failed"))
	}
	count, err := service.IamRoleDomain().AdminDeleteRole(ctx, in.Id)
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "admin.role.check_users_failed"))
	}
	if count > 0 {
		return gerror.Newf(gi18n.T(ctx, "admin.role.delete_has_users"), count)
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionDelete,
		TargetType: consts.SiteAuditTargetTypeIamRole,
		TargetId:   uint64(in.Id),
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"snapshot": s.iamRoleAuditSnapshot(role),
		},
	})
	return nil
}

func (s *sAdminIamRoleUsecase) iamRoleAuditSnapshot(role *entity.IamRole) map[string]any {
	if role == nil {
		return nil
	}
	return map[string]any{
		"id":       role.Id,
		"level":    role.Level,
		"nameI18N": s.iamRoleJSONMap(role.NameI18N),
		"isStaff":  role.IsStaff,
	}
}

func (s *sAdminIamRoleUsecase) iamRoleJSONMap(value *gjson.Json) map[string]any {
	if value == nil {
		return nil
	}
	var data map[string]any
	if err := value.Scan(&data); err != nil {
		return nil
	}
	return data
}
