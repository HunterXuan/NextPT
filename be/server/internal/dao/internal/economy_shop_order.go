// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EconomyShopOrderDao is the data access object for the table economy_shop_order.
type EconomyShopOrderDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  EconomyShopOrderColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// EconomyShopOrderColumns defines and stores column names for the table economy_shop_order.
type EconomyShopOrderColumns struct {
	Id              string //
	UserId          string //
	ProductKey      string // 商品标识快照
	ProductType     string // 商品类型快照
	ProductSnapshot string // 购买时商品配置快照
	Price           string // 成交价格快照
	Status          string // 0=待处理, 1=已完成
	TargetType      string // 履约目标类型
	TargetId        string // 履约目标ID
	CreatedAt       string //
	CompletedAt     string //
}

// economyShopOrderColumns holds the columns for the table economy_shop_order.
var economyShopOrderColumns = EconomyShopOrderColumns{
	Id:              "id",
	UserId:          "user_id",
	ProductKey:      "product_key",
	ProductType:     "product_type",
	ProductSnapshot: "product_snapshot",
	Price:           "price",
	Status:          "status",
	TargetType:      "target_type",
	TargetId:        "target_id",
	CreatedAt:       "created_at",
	CompletedAt:     "completed_at",
}

// NewEconomyShopOrderDao creates and returns a new DAO object for table data access.
func NewEconomyShopOrderDao(handlers ...gdb.ModelHandler) *EconomyShopOrderDao {
	return &EconomyShopOrderDao{
		group:    "default",
		table:    "economy_shop_order",
		columns:  economyShopOrderColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *EconomyShopOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *EconomyShopOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *EconomyShopOrderDao) Columns() EconomyShopOrderColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *EconomyShopOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *EconomyShopOrderDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *EconomyShopOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
