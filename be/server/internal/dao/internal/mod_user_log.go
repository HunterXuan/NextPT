// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ModUserLogDao is the data access object for the table mod_user_log.
type ModUserLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ModUserLogColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ModUserLogColumns defines and stores column names for the table mod_user_log.
type ModUserLogColumns struct {
	Id         string //
	UserId     string //
	ModType    string // 1=warned 2=banned 3=leech_warned 4=upload_banned 5=download_banned 6=forum_banned
	Reason     string //
	ExpireAt   string // 过期时间，NULL=永久
	ModBy      string // 操作人ID
	ModComment string //
	IsActive   string //
	CreatedAt  string //
	UpdatedAt  string //
}

// modUserLogColumns holds the columns for the table mod_user_log.
var modUserLogColumns = ModUserLogColumns{
	Id:         "id",
	UserId:     "user_id",
	ModType:    "mod_type",
	Reason:     "reason",
	ExpireAt:   "expire_at",
	ModBy:      "mod_by",
	ModComment: "mod_comment",
	IsActive:   "is_active",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewModUserLogDao creates and returns a new DAO object for table data access.
func NewModUserLogDao(handlers ...gdb.ModelHandler) *ModUserLogDao {
	return &ModUserLogDao{
		group:    "default",
		table:    "mod_user_log",
		columns:  modUserLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ModUserLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ModUserLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ModUserLogDao) Columns() ModUserLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ModUserLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ModUserLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ModUserLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
