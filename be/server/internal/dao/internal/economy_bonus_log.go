// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EconomyBonusLogDao is the data access object for the table economy_bonus_log.
type EconomyBonusLogDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  EconomyBonusLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// EconomyBonusLogColumns defines and stores column names for the table economy_bonus_log.
type EconomyBonusLogColumns struct {
	Id           string //
	UserId       string //
	Amount       string // 正=获得 负=消耗
	BalanceAfter string // 变动后余额
	Action       string // torrent_reward/post_reward/daily_bonus/...
	TargetType   string //
	TargetId     string //
	Period       string // 结算周期/幂等键
	Remark       string //
	CreatedAt    string //
}

// economyBonusLogColumns holds the columns for the table economy_bonus_log.
var economyBonusLogColumns = EconomyBonusLogColumns{
	Id:           "id",
	UserId:       "user_id",
	Amount:       "amount",
	BalanceAfter: "balance_after",
	Action:       "action",
	TargetType:   "target_type",
	TargetId:     "target_id",
	Period:       "period",
	Remark:       "remark",
	CreatedAt:    "created_at",
}

// NewEconomyBonusLogDao creates and returns a new DAO object for table data access.
func NewEconomyBonusLogDao(handlers ...gdb.ModelHandler) *EconomyBonusLogDao {
	return &EconomyBonusLogDao{
		group:    "default",
		table:    "economy_bonus_log",
		columns:  economyBonusLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *EconomyBonusLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *EconomyBonusLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *EconomyBonusLogDao) Columns() EconomyBonusLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *EconomyBonusLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *EconomyBonusLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *EconomyBonusLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
