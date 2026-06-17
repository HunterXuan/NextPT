package iam

import (
	"context"
	"strings"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sIamPermissionDomain struct{}

func init() {
	service.RegisterIamPermissionDomain(NewIamPermissionDomain())
}

func NewIamPermissionDomain() *sIamPermissionDomain {
	return &sIamPermissionDomain{}
}

func (s *sIamPermissionDomain) CheckPermissionWithList(ctx context.Context, rolePerms []string, userAcls []string, permKey string) (bool, error) {
	hasPerm := false
	for _, p := range rolePerms {
		if s.matchPermissionKey(p, permKey) {
			hasPerm = true
			break
		}
	}

	for _, acl := range userAcls {
		if strings.HasPrefix(acl, "-") {
			denyKey := strings.TrimPrefix(acl, "-")
			if s.matchPermissionKey(denyKey, permKey) {
				return false, nil
			}
			continue
		}

		if s.matchPermissionKey(acl, permKey) {
			hasPerm = true
		}
	}

	return hasPerm, nil
}

func (s *sIamPermissionDomain) matchPermissionKey(pattern string, permKey string) bool {
	if pattern == "" {
		return false
	}
	if pattern == consts.IamPermissionAll || pattern == permKey {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(permKey, strings.TrimSuffix(pattern, "*"))
	}
	return false
}

func (s *sIamPermissionDomain) GrantUserPermission(ctx context.Context, userId uint64, permKey string, isDeny bool) error {
	permKey = s.normalizePermissionKey(permKey, isDeny)

	columns := dao.IamUserPermission.Columns()
	_, err := dao.IamUserPermission.Ctx(ctx).Data(g.Map{
		columns.UserId:     userId,
		columns.PermKey:    permKey,
		columns.SourceType: consts.IamUserPermissionSourceManual,
		columns.SourceId:   0,
		columns.IsActive:   true,
	}).InsertIgnore()
	return err
}

func (s *sIamPermissionDomain) RevokeUserPermission(ctx context.Context, userId uint64, permKey string, isDeny bool) error {
	permKey = s.normalizePermissionKey(permKey, isDeny)

	columns := dao.IamUserPermission.Columns()
	_, err := dao.IamUserPermission.Ctx(ctx).
		Where(columns.UserId, userId).
		Where(columns.PermKey, permKey).
		Where(columns.SourceType, consts.IamUserPermissionSourceManual).
		Where(columns.SourceId, 0).
		Delete()
	return err
}

func (s *sIamPermissionDomain) GrantUserPermissionsBySource(ctx context.Context, userId uint64, permKeys []string, sourceType int, sourceId uint64, expireAt *gtime.Time) error {
	if userId == 0 || len(permKeys) == 0 || sourceType == 0 {
		return nil
	}

	columns := dao.IamUserPermission.Columns()
	for _, permKey := range permKeys {
		if permKey == "" {
			continue
		}
		_, err := dao.IamUserPermission.Ctx(ctx).Data(g.Map{
			columns.UserId:     userId,
			columns.PermKey:    permKey,
			columns.SourceType: sourceType,
			columns.SourceId:   sourceId,
			columns.ExpireAt:   expireAt,
			columns.IsActive:   true,
		}).InsertIgnore()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sIamPermissionDomain) DeactivateUserPermissionsBySource(ctx context.Context, userId uint64, sourceType int, sourceId uint64) error {
	if userId == 0 || sourceType == 0 {
		return nil
	}

	columns := dao.IamUserPermission.Columns()
	_, err := dao.IamUserPermission.Ctx(ctx).
		Where(columns.UserId, userId).
		Where(columns.SourceType, sourceType).
		Where(columns.SourceId, sourceId).
		Data(g.Map{
			columns.IsActive:  false,
			columns.UpdatedAt: gtime.Now(),
		}).
		Update()
	return err
}

func (s *sIamPermissionDomain) normalizePermissionKey(permKey string, isDeny bool) string {
	if isDeny && !strings.HasPrefix(permKey, "-") {
		return "-" + permKey
	}
	if !isDeny && strings.HasPrefix(permKey, "-") {
		return strings.TrimPrefix(permKey, "-")
	}
	return permKey
}

func (s *sIamPermissionDomain) GetAllPermissions(ctx context.Context) []string {
	return consts.IamPermissionList
}

func (s *sIamPermissionDomain) GetUserPermissions(ctx context.Context, userId uint64) ([]entity.IamUserPermission, error) {
	var acls []entity.IamUserPermission
	columns := dao.IamUserPermission.Columns()
	m := dao.IamUserPermission.Ctx(ctx)
	err := m.
		Fields(columns.PermKey).
		Where(columns.UserId, userId).
		Where(columns.IsActive, true).
		Where(m.Builder().WhereNull(columns.ExpireAt).WhereOrGT(columns.ExpireAt, gtime.Now())).
		Scan(&acls)
	return acls, err
}
