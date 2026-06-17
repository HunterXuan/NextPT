package accounting

import (
	"context"
	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/accounting/v1"
)

func (c *ControllerV1) TrafficGetMe(ctx context.Context, req *v1.TrafficGetMeReq) (res *v1.TrafficGetMeRes, err error) {
	actor := contexts.GetActor(ctx)
	out, err := service.AccountingTrafficUsecase().GetMyTraffic(ctx, actor)
	if err != nil {
		return nil, err
	}
	return &v1.TrafficGetMeRes{TrafficGetMeOut: *out}, nil
}
