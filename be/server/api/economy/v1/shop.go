package v1

import (
	"server/internal/model/in/economyin"
	"server/internal/model/out/economyout"

	"github.com/gogf/gf/v2/frame/g"
)

type ShopProductListReq struct {
	g.Meta `path:"/shop/products" method:"get" tags:"EconomyShop" summary:"获取魔力商城商品" perm:"read:economy/shop-product:*"`
	economyin.ShopProductListInp
}

type ShopProductListRes struct {
	economyout.ShopProductListOut
}

type ShopOrderCreateReq struct {
	g.Meta `path:"/shop/orders" method:"post" tags:"EconomyShop" summary:"兑换魔力商城商品" perm:"create:economy/shop-order:*"`
	economyin.ShopOrderCreateInp
}

type ShopOrderCreateRes struct {
	economyout.ShopOrderCreateOut
}

type ShopOrderListReq struct {
	g.Meta `path:"/shop/orders" method:"get" tags:"EconomyShop" summary:"获取我的商城订单" perm:"read:economy/shop-order:*"`
	economyin.ShopOrderListInp
}

type ShopOrderListRes struct {
	economyout.ShopOrderListOut
}
