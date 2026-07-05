package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserPermissionDetail(ctx context.Context, req *v1.IamUserPermissionDetailReq) (res *v1.IamUserPermissionDetailRes, err error) {
	out, err := service.AdminIamPermissionUsecase().UserDetail(ctx, contexts.GetActor(ctx), req.IamUserPermissionDetailInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamUserPermissionDetailRes{IamUserPermissionDetailOut: *out}, nil
}
