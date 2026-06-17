package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamSessionDelete(ctx context.Context, req *v1.IamSessionDeleteReq) (res *v1.IamSessionDeleteRes, err error) {
	err = service.AdminIamSessionUsecase().DeleteByUser(ctx, contexts.GetActor(ctx), req.IamSessionDeleteInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamSessionDeleteRes{}, nil
}
