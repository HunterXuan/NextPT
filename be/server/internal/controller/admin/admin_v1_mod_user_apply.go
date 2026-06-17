package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ModUserApply(ctx context.Context, req *v1.ModUserApplyReq) (res *v1.ModUserApplyRes, err error) {
	err = service.AdminModUserUsecase().Apply(ctx, contexts.GetActor(ctx), req.ModUserApplyInp)
	if err != nil {
		return nil, err
	}
	return &v1.ModUserApplyRes{}, nil
}
