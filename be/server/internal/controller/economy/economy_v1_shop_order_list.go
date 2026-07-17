package economy

import (
	"context"

	v1 "server/api/economy/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ShopOrderList(ctx context.Context, req *v1.ShopOrderListReq) (res *v1.ShopOrderListRes, err error) {
	out, err := service.EconomyShopUsecase().ListMyOrders(ctx, contexts.GetActor(ctx), req.ShopOrderListInp)
	if err != nil {
		return nil, err
	}
	return &v1.ShopOrderListRes{ShopOrderListOut: *out}, nil
}
