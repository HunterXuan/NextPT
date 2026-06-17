package admin

import (
	"context"
	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserPermissionGrant(ctx context.Context, req *v1.IamUserPermissionGrantReq) (res *v1.IamUserPermissionGrantRes, err error) {
	err = service.AdminIamPermissionUsecase().GrantUserAcl(ctx, contexts.GetActor(ctx), req.IamUserPermissionGrantInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamUserPermissionGrantRes{}, nil
}
