package economy

import (
	"context"

	v1 "server/api/economy/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ShopProductList(ctx context.Context, req *v1.ShopProductListReq) (res *v1.ShopProductListRes, err error) {
	out, err := service.EconomyShopUsecase().ListProducts(ctx, contexts.GetActor(ctx), req.ShopProductListInp)
	if err != nil {
		return nil, err
	}
	return &v1.ShopProductListRes{ShopProductListOut: *out}, nil
}
