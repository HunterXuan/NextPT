package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamPermissionList(ctx context.Context, req *v1.IamPermissionListReq) (res *v1.IamPermissionListRes, err error) {
	res = &v1.IamPermissionListRes{}
	out, err := service.AdminIamPermissionUsecase().List(ctx, contexts.GetActor(ctx), req.IamPermissionListInp)
	if err == nil {
		res.IamPermissionListOut = *out
	}
	return
}
