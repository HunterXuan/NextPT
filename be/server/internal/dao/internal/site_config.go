// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SiteConfigDao is the data access object for the table site_config.
type SiteConfigDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SiteConfigColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SiteConfigColumns defines and stores column names for the table site_config.
type SiteConfigColumns struct {
	Id          string //
	Group       string //
	Key         string //
	Value       string // 支持存 boolean/number/array
	Description string //
	CreatedAt   string //
	UpdatedAt   string //
}

// siteConfigColumns holds the columns for the table site_config.
var siteConfigColumns = SiteConfigColumns{
	Id:          "id",
	Group:       "group",
	Key:         "key",
	Value:       "value",
	Description: "description",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewSiteConfigDao creates and returns a new DAO object for table data access.
func NewSiteConfigDao(handlers ...gdb.ModelHandler) *SiteConfigDao {
	return &SiteConfigDao{
		group:    "default",
		table:    "site_config",
		columns:  siteConfigColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SiteConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SiteConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SiteConfigDao) Columns() SiteConfigColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SiteConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SiteConfigDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SiteConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
