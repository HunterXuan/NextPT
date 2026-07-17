package economyout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
)

type ShopProductItem = model.EconomyShopProductConfig

type ShopProductListOut struct {
	List []ShopProductItem `json:"list"`
}

type ShopOrderCreateOut struct {
	OrderId      uint64  `json:"orderId"`
	ProductKey   string  `json:"productKey"`
	Price        float64 `json:"price"`
	TargetType   string  `json:"targetType"`
	TargetId     uint64  `json:"targetId"`
	BalanceAfter float64 `json:"balanceAfter"`
}

type ShopOrderItem struct {
	Id          uint64          `json:"id"`
	Product     ShopProductItem `json:"product"`
	Price       float64         `json:"price"`
	Status      int             `json:"status"`
	TargetType  string          `json:"targetType"`
	TargetId    uint64          `json:"targetId"`
	CreatedAt   *gtime.Time     `json:"createdAt"`
	CompletedAt *gtime.Time     `json:"completedAt"`
}

type ShopOrderListOut struct {
	Page  int             `json:"page"`
	Size  int             `json:"size"`
	Total int             `json:"total"`
	List  []ShopOrderItem `json:"list"`
}
