package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserUpdate(ctx context.Context, req *v1.IamUserUpdateReq) (res *v1.IamUserUpdateRes, err error) {
	err = service.AdminIamUserUsecase().Update(ctx, contexts.GetActor(ctx), req.IamUserUpdateInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamUserUpdateRes{}, nil
}
