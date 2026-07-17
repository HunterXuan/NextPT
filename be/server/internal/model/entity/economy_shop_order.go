// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// EconomyShopOrder is the golang structure for table economy_shop_order.
type EconomyShopOrder struct {
	Id              uint64      `json:"id"              orm:"id"               description:""`
	UserId          uint64      `json:"userId"          orm:"user_id"          description:""`
	ProductKey      string      `json:"productKey"      orm:"product_key"      description:"商品标识快照"`
	ProductType     string      `json:"productType"     orm:"product_type"     description:"商品类型快照"`
	ProductSnapshot *gjson.Json `json:"productSnapshot" orm:"product_snapshot" description:"购买时商品配置快照"`
	Price           float64     `json:"price"           orm:"price"            description:"成交价格快照"`
	Status          int         `json:"status"          orm:"status"           description:"0=待处理, 1=已完成"`
	TargetType      string      `json:"targetType"      orm:"target_type"      description:"履约目标类型"`
	TargetId        uint64      `json:"targetId"        orm:"target_id"        description:"履约目标ID"`
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:""`
	CompletedAt     *gtime.Time `json:"completedAt"     orm:"completed_at"     description:""`
}
