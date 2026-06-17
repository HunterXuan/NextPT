package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"
	"server/internal/service"
)

type sAdminIamPermissionUsecase struct{}

func init() {
	service.RegisterAdminIamPermissionUsecase(NewAdminIamPermissionUsecase())
}

func NewAdminIamPermissionUsecase() *sAdminIamPermissionUsecase {
	return &sAdminIamPermissionUsecase{}
}

func (s *sAdminIamPermissionUsecase) GrantUserAcl(ctx context.Context, actor *model.Actor, in adminin.IamUserPermissionGrantInp) error {
	err := service.IamPermissionDomain().GrantUserPermission(ctx, in.UserId, in.PermKey, in.IsDeny)
	if err == nil {
		service.IamUserUsecase().InvalidateUserCache(ctx, in.UserId)
	}
	return err
}

func (s *sAdminIamPermissionUsecase) RevokeUserAcl(ctx context.Context, actor *model.Actor, in adminin.IamUserPermissionRevokeInp) error {
	err := service.IamPermissionDomain().RevokeUserPermission(ctx, in.UserId, in.PermKey, in.IsDeny)
	if err == nil {
		service.IamUserUsecase().InvalidateUserCache(ctx, in.UserId)
	}
	return err
}

func (s *sAdminIamPermissionUsecase) List(ctx context.Context, actor *model.Actor, in adminin.IamPermissionListInp) (*adminout.IamPermissionListOut, error) {
	return &adminout.IamPermissionListOut{Permissions: service.IamPermissionDomain().GetAllPermissions(ctx)}, nil
}
