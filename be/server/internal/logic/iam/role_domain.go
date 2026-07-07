package iam

import (
	"context"

	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sIamRoleDomain struct{}

func init() {
	service.RegisterIamRoleDomain(NewIamRoleDomain())
}

func NewIamRoleDomain() *sIamRoleDomain {
	return &sIamRoleDomain{}
}

func (s *sIamRoleDomain) GetRoleById(ctx context.Context, roleId uint) (*entity.IamRole, error) {
	var role *entity.IamRole
	err := dao.IamRole.Ctx(ctx).Where(dao.IamRole.Columns().Id, roleId).Scan(&role)
	return role, err
}

func (s *sIamRoleDomain) GetRolesByIds(ctx context.Context, roleIds []uint) ([]entity.IamRole, error) {
	if len(roleIds) == 0 {
		return nil, nil
	}
	var roles []entity.IamRole
	err := dao.IamRole.Ctx(ctx).
		WhereIn(dao.IamRole.Columns().Id, roleIds).
		Scan(&roles)
	return roles, err
}

func (s *sIamRoleDomain) ListRoles(ctx context.Context) ([]entity.IamRole, error) {
	var roles []entity.IamRole
	err := dao.IamRole.Ctx(ctx).
		OrderAsc(dao.IamRole.Columns().Level).
		OrderAsc(dao.IamRole.Columns().Id).
		Scan(&roles)
	return roles, err
}

func (s *sIamRoleDomain) AdminListRoles(ctx context.Context) ([]entity.IamRole, error) {
	var roles []entity.IamRole
	err := dao.IamRole.Ctx(ctx).
		OrderAsc(dao.IamRole.Columns().Level).
		OrderAsc(dao.IamRole.Columns().Id).
		Scan(&roles)
	return roles, err
}

func (s *sIamRoleDomain) AdminCreateRole(ctx context.Context, level int, nameI18N, rules, permissions []byte, isStaff bool) (uint, error) {
	id, err := dao.IamRole.Ctx(ctx).Data(g.Map{
		dao.IamRole.Columns().Level:       level,
		dao.IamRole.Columns().NameI18N:    nameI18N,
		dao.IamRole.Columns().Rules:       rules,
		dao.IamRole.Columns().Permissions: permissions,
		dao.IamRole.Columns().IsStaff:     isStaff,
	}).InsertAndGetId()
	return uint(id), err
}

func (s *sIamRoleDomain) AdminUpdateRole(ctx context.Context, id uint, level *int, nameI18N, rules, permissions []byte, isStaff *bool) error {
	data := g.Map{}
	if level != nil {
		data[dao.IamRole.Columns().Level] = *level
	}
	if nameI18N != nil {
		data[dao.IamRole.Columns().NameI18N] = nameI18N
	}
	if rules != nil {
		data[dao.IamRole.Columns().Rules] = rules
	}
	if permissions != nil {
		data[dao.IamRole.Columns().Permissions] = permissions
	}
	if isStaff != nil {
		data[dao.IamRole.Columns().IsStaff] = *isStaff
	}
	if len(data) == 0 {
		return nil
	}
	_, err := dao.IamRole.Ctx(ctx).Where(dao.IamRole.Columns().Id, id).Data(data).Update()
	return err
}

func (s *sIamRoleDomain) AdminDeleteRole(ctx context.Context, id uint) (int, error) {
	count, err := dao.IamUser.Ctx(ctx).Where(dao.IamUser.Columns().Role, id).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return count, nil
	}
	_, err = dao.IamRole.Ctx(ctx).Where(dao.IamRole.Columns().Id, id).Delete()
	return 0, err
}
