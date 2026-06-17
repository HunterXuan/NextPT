package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserStatUpdate(ctx context.Context, req *v1.IamUserStatUpdateReq) (res *v1.IamUserStatUpdateRes, err error) {
	err = service.AdminIamUserUsecase().StatUpdate(ctx, contexts.GetActor(ctx), req.IamUserStatUpdateInp)
	if err == nil {
		res = &v1.IamUserStatUpdateRes{}
	}
	return
}
