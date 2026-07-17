// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EconomyShopOrder is the golang structure of table economy_shop_order for DAO operations like Where/Data.
type EconomyShopOrder struct {
	g.Meta          `orm:"table:economy_shop_order, do:true"`
	Id              any         //
	UserId          any         //
	ProductKey      any         // 商品标识快照
	ProductType     any         // 商品类型快照
	ProductSnapshot *gjson.Json // 购买时商品配置快照
	Price           any         // 成交价格快照
	Status          any         // 0=待处理, 1=已完成
	TargetType      any         // 履约目标类型
	TargetId        any         // 履约目标ID
	CreatedAt       *gtime.Time //
	CompletedAt     *gtime.Time //
}
