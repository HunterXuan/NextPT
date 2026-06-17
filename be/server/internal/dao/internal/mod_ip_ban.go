// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ModIpBanDao is the data access object for the table mod_ip_ban.
type ModIpBanDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ModIpBanColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ModIpBanColumns defines and stores column names for the table mod_ip_ban.
type ModIpBanColumns struct {
	Id        string //
	Network   string // IP或CIDR，如 192.168.1.0/24
	Reason    string //
	BannedBy  string //
	CreatedAt string //
}

// modIpBanColumns holds the columns for the table mod_ip_ban.
var modIpBanColumns = ModIpBanColumns{
	Id:        "id",
	Network:   "network",
	Reason:    "reason",
	BannedBy:  "banned_by",
	CreatedAt: "created_at",
}

// NewModIpBanDao creates and returns a new DAO object for table data access.
func NewModIpBanDao(handlers ...gdb.ModelHandler) *ModIpBanDao {
	return &ModIpBanDao{
		group:    "default",
		table:    "mod_ip_ban",
		columns:  modIpBanColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ModIpBanDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ModIpBanDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ModIpBanDao) Columns() ModIpBanColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ModIpBanDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ModIpBanDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ModIpBanDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
