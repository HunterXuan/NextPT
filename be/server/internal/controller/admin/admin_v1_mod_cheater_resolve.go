package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ModCheaterResolve(ctx context.Context, req *v1.ModCheaterResolveReq) (res *v1.ModCheaterResolveRes, err error) {
	err = service.AdminModCheaterUsecase().Resolve(ctx, contexts.GetActor(ctx), req.ModCheaterResolveInp)
	if err != nil {
		return nil, err
	}
	return &v1.ModCheaterResolveRes{}, nil
}
