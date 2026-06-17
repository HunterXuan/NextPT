package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ModCheaterList(ctx context.Context, req *v1.ModCheaterListReq) (res *v1.ModCheaterListRes, err error) {
	out, err := service.AdminModCheaterUsecase().List(ctx, contexts.GetActor(ctx), req.ModCheaterListInp)
	if err != nil {
		return nil, err
	}
	return &v1.ModCheaterListRes{
		ListCheaterLogsOut: *out,
	}, nil
}
