package admin

import (
	"context"
	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserPermissionRevoke(ctx context.Context, req *v1.IamUserPermissionRevokeReq) (res *v1.IamUserPermissionRevokeRes, err error) {
	err = service.AdminIamPermissionUsecase().RevokeUserAcl(ctx, contexts.GetActor(ctx), req.IamUserPermissionRevokeInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamUserPermissionRevokeRes{}, nil
}
