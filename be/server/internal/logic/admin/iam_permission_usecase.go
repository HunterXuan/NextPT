package admin

import (
	"context"
	"sort"
	"strings"

	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sAdminIamPermissionUsecase struct{}

func init() {
	service.RegisterAdminIamPermissionUsecase(NewAdminIamPermissionUsecase())
}

func NewAdminIamPermissionUsecase() *sAdminIamPermissionUsecase {
	return &sAdminIamPermissionUsecase{}
}

func (s *sAdminIamPermissionUsecase) GrantUserAcl(ctx context.Context, actor *model.Actor, in adminin.IamUserPermissionGrantInp) error {
	permKeys, err := s.normalizeWildcardPermissions(ctx, in.PermKeys)
	if err != nil {
		return err
	}

	err = service.IamPermissionDomain().GrantUserPermissions(ctx, in.UserId, permKeys, in.IsDeny)
	if err == nil {
		service.IamUserUsecase().InvalidateUserCache(ctx, in.UserId)
	}
	return err
}

func (s *sAdminIamPermissionUsecase) RevokeUserAcl(ctx context.Context, actor *model.Actor, in adminin.IamUserPermissionRevokeInp) error {
	ids := s.normalizeAclIds(in.Ids)
	if len(ids) == 0 {
		return gerror.New(gi18n.T(ctx, "admin.permission.empty"))
	}

	err := service.IamPermissionDomain().RevokeUserPermissionsByIds(ctx, in.UserId, ids)
	if err == nil {
		service.IamUserUsecase().InvalidateUserCache(ctx, in.UserId)
	}
	return err
}

func (s *sAdminIamPermissionUsecase) List(ctx context.Context, actor *model.Actor, in adminin.IamPermissionListInp) (*adminout.IamPermissionListOut, error) {
	return &adminout.IamPermissionListOut{Permissions: service.IamPermissionDomain().GetAllPermissions(ctx)}, nil
}

func (s *sAdminIamPermissionUsecase) isWildcardPermission(permission string) bool {
	return permission == "*" || strings.HasSuffix(permission, ":*")
}

func (s *sAdminIamPermissionUsecase) normalizeWildcardPermissions(ctx context.Context, permissions []string) ([]string, error) {
	list := make([]string, 0, len(permissions))
	seen := make(map[string]struct{}, len(permissions))
	for _, permission := range permissions {
		permission = strings.TrimSpace(strings.TrimPrefix(permission, "-"))
		if permission == "" {
			continue
		}
		if !s.isWildcardPermission(permission) {
			return nil, gerror.New(gi18n.T(ctx, "admin.permission.wildcard_only"))
		}
		if _, ok := seen[permission]; ok {
			continue
		}
		seen[permission] = struct{}{}
		list = append(list, permission)
	}
	if len(list) == 0 {
		return nil, gerror.New(gi18n.T(ctx, "admin.permission.empty"))
	}
	sort.Strings(list)
	return list, nil
}

func (s *sAdminIamPermissionUsecase) normalizeAclIds(ids []uint64) []uint64 {
	list := make([]uint64, 0, len(ids))
	seen := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		list = append(list, id)
	}
	return list
}

func (s *sAdminIamPermissionUsecase) UserDetail(ctx context.Context, actor *model.Actor, in adminin.IamUserPermissionDetailInp) (*adminout.IamUserPermissionDetailOut, error) {
	user, err := service.IamUserDomain().GetUserById(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.user.not_found"))
	}

	acls, total, err := service.IamPermissionDomain().ListUserPermissions(ctx, in.UserId, model.IamUserPermissionListOptions{
		SourceType:   in.SourceType,
		WildcardOnly: in.WildcardOnly,
		Page:         in.Page,
		Size:         in.Size,
	})
	if err != nil {
		return nil, err
	}

	items := make([]adminout.IamUserAclSummary, 0, len(acls))
	for _, acl := range acls {
		rawPermKey := acl.PermKey
		isDeny := strings.HasPrefix(rawPermKey, "-")
		permKey := strings.TrimPrefix(rawPermKey, "-")
		items = append(items, adminout.IamUserAclSummary{
			Id:         acl.Id,
			UserId:     acl.UserId,
			PermKey:    permKey,
			RawPermKey: rawPermKey,
			IsDeny:     isDeny,
			SourceType: acl.SourceType,
			SourceId:   acl.SourceId,
			ExpireAt:   acl.ExpireAt,
			IsActive:   acl.IsActive,
			CreatedAt:  acl.CreatedAt,
		})
	}

	return &adminout.IamUserPermissionDetailOut{
		UserAcls: items,
		Total:    total,
		Page:     in.Page,
		Size:     in.Size,
	}, nil
}
