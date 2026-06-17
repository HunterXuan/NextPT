package accounting

import (
	"context"
	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/accounting/v1"
)

func (c *ControllerV1) TrafficHistoryList(ctx context.Context, req *v1.TrafficHistoryListReq) (res *v1.TrafficHistoryListRes, err error) {
	actor := contexts.GetActor(ctx)
	out, err := service.AccountingTrafficUsecase().ListMyTrafficHistory(ctx, actor, req.TrafficHistoryListInp)
	if err != nil {
		return nil, err
	}
	return &v1.TrafficHistoryListRes{TrafficHistoryListOut: *out}, nil
}
