package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamRoleDelete(ctx context.Context, req *v1.IamRoleDeleteReq) (res *v1.IamRoleDeleteRes, err error) {
	err = service.AdminIamRoleUsecase().Delete(ctx, contexts.GetActor(ctx), req.IamRoleDeleteInp)
	if err == nil {
		res = &v1.IamRoleDeleteRes{}
	}
	return
}
