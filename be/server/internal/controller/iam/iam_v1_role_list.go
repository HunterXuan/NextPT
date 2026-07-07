package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) RoleList(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error) {
	out, err := service.IamRoleUsecase().List(ctx, contexts.GetActor(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.RoleListRes{RoleListOut: out}, nil
}
