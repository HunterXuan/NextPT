package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserPermissionList(ctx context.Context, req *v1.IamUserPermissionListReq) (res *v1.IamUserPermissionListRes, err error) {
	out, err := service.IamUserUsecase().Permissions(ctx, contexts.GetActor(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.IamUserPermissionListRes{UserPermissionListOut: *out}, nil
}
