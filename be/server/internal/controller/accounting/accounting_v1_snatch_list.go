package accounting

import (
	"context"
	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/accounting/v1"
)

func (c *ControllerV1) SnatchList(ctx context.Context, req *v1.SnatchListReq) (res *v1.SnatchListRes, err error) {
	actor := contexts.GetActor(ctx)
	out, err := service.AccountingSnatchUsecase().ListMySnatches(ctx, actor, req.SnatchListInp)
	if err != nil {
		return nil, err
	}
	return &v1.SnatchListRes{SnatchListOut: *out}, nil
}
