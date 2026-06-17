package accounting

import (
	"context"
	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/accounting/v1"
)

func (c *ControllerV1) SnatchGet(ctx context.Context, req *v1.SnatchGetReq) (res *v1.SnatchGetRes, err error) {
	actor := contexts.GetActor(ctx)
	out, err := service.AccountingSnatchUsecase().GetMySnatch(ctx, actor, req.SnatchGetInp)
	if err != nil {
		return nil, err
	}
	return &v1.SnatchGetRes{SnatchGetOut: *out}, nil
}
