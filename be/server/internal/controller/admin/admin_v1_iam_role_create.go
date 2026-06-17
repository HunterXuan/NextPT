package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamRoleCreate(ctx context.Context, req *v1.IamRoleCreateReq) (res *v1.IamRoleCreateRes, err error) {
	id, err := service.AdminIamRoleUsecase().Create(ctx, contexts.GetActor(ctx), req.IamRoleCreateInp)
	if err != nil {
		return nil, err
	}
	res = &v1.IamRoleCreateRes{}
	res.Id = id
	return
}
