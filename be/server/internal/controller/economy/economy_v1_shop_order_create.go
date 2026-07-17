package economy

import (
	"context"

	v1 "server/api/economy/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ShopOrderCreate(ctx context.Context, req *v1.ShopOrderCreateReq) (res *v1.ShopOrderCreateRes, err error) {
	out, err := service.EconomyShopUsecase().CreateOrder(ctx, contexts.GetActor(ctx), req.ShopOrderCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.ShopOrderCreateRes{ShopOrderCreateOut: *out}, nil
}
