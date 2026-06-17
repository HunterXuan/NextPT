package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ModUserRemove(ctx context.Context, req *v1.ModUserRemoveReq) (res *v1.ModUserRemoveRes, err error) {
	err = service.AdminModUserUsecase().Remove(ctx, contexts.GetActor(ctx), req.ModUserRemoveInp)
	if err != nil {
		return nil, err
	}
	return &v1.ModUserRemoveRes{}, nil
}
