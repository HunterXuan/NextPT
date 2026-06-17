package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamRoleList(ctx context.Context, req *v1.IamRoleListReq) (res *v1.IamRoleListRes, err error) {
	out, err := service.AdminIamRoleUsecase().List(ctx, contexts.GetActor(ctx))
	if err != nil {
		return nil, err
	}
	res = &v1.IamRoleListRes{IamRoleListOut: *out}
	return
}
