package iam

import (
	"context"

	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/out/iamout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sIamRoleUsecase struct{}

func init() {
	service.RegisterIamRoleUsecase(NewIamRoleUsecase())
}

func NewIamRoleUsecase() *sIamRoleUsecase {
	return &sIamRoleUsecase{}
}

func (s *sIamRoleUsecase) List(ctx context.Context, actor *model.Actor) (*iamout.RoleListOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	roles, err := service.IamRoleDomain().ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]iamout.RoleItem, 0, len(roles))
	for _, role := range roles {
		items = append(items, iamout.RoleItem{
			Id:       role.Id,
			Level:    role.Level,
			Name:     s.localizeName(ctx, role),
			NameI18N: s.scanNameI18N(role),
			Rules:    s.scanRules(role),
			IsStaff:  role.IsStaff,
		})
	}

	return &iamout.RoleListOut{Roles: items}, nil
}

func (s *sIamRoleUsecase) scanNameI18N(role entity.IamRole) map[string]string {
	names := make(map[string]string)
	if role.NameI18N == nil {
		return names
	}
	_ = role.NameI18N.Scan(&names)
	return names
}

func (s *sIamRoleUsecase) scanRules(role entity.IamRole) map[string]any {
	rules := make(map[string]any)
	if role.Rules == nil {
		return rules
	}
	_ = role.Rules.Scan(&rules)
	return rules
}

func (s *sIamRoleUsecase) localizeName(ctx context.Context, role entity.IamRole) string {
	names := s.scanNameI18N(role)
	if len(names) == 0 {
		return ""
	}

	lang := gi18n.LanguageFromCtx(ctx)
	if lang != "" && names[lang] != "" {
		return names[lang]
	}

	defaultLang := g.Cfg().MustGet(ctx, "i18n.default", "zh-CN").String()
	if defaultLang != "" && names[defaultLang] != "" {
		return names[defaultLang]
	}

	for _, name := range names {
		if name != "" {
			return name
		}
	}
	return ""
}
