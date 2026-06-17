package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamRoleUpdate(ctx context.Context, req *v1.IamRoleUpdateReq) (res *v1.IamRoleUpdateRes, err error) {
	err = service.AdminIamRoleUsecase().Update(ctx, contexts.GetActor(ctx), req.IamRoleUpdateInp)
	if err == nil {
		res = &v1.IamRoleUpdateRes{}
	}
	return
}
