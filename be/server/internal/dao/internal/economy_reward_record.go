// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EconomyRewardRecordDao is the data access object for the table economy_reward_record.
type EconomyRewardRecordDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  EconomyRewardRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// EconomyRewardRecordColumns defines and stores column names for the table economy_reward_record.
type EconomyRewardRecordColumns struct {
	Id         string //
	TargetType string // catalog_torrent/catalog_comment/forum_topic/forum_reply
	TargetId   string //
	FromUserId string // 赞赏者ID
	ToUserId   string // 接收者ID
	Amount     string // 赞赏金额
	CreatedAt  string //
}

// economyRewardRecordColumns holds the columns for the table economy_reward_record.
var economyRewardRecordColumns = EconomyRewardRecordColumns{
	Id:         "id",
	TargetType: "target_type",
	TargetId:   "target_id",
	FromUserId: "from_user_id",
	ToUserId:   "to_user_id",
	Amount:     "amount",
	CreatedAt:  "created_at",
}

// NewEconomyRewardRecordDao creates and returns a new DAO object for table data access.
func NewEconomyRewardRecordDao(handlers ...gdb.ModelHandler) *EconomyRewardRecordDao {
	return &EconomyRewardRecordDao{
		group:    "default",
		table:    "economy_reward_record",
		columns:  economyRewardRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *EconomyRewardRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *EconomyRewardRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *EconomyRewardRecordDao) Columns() EconomyRewardRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *EconomyRewardRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *EconomyRewardRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *EconomyRewardRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
