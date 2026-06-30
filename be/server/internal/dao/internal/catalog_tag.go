// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogTagDao is the data access object for the table catalog_tag.
type CatalogTagDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  CatalogTagColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// CatalogTagColumns defines and stores column names for the table catalog_tag.
type CatalogTagColumns struct {
	Id        string //
	GroupId   string // 所属分组
	NameI18N  string // 多语言名称映射
	Value     string // 稳定值
	SortOrder string //
	CreatedAt string //
	UpdatedAt string //
}

// catalogTagColumns holds the columns for the table catalog_tag.
var catalogTagColumns = CatalogTagColumns{
	Id:        "id",
	GroupId:   "group_id",
	NameI18N:  "name_i18n",
	Value:     "value",
	SortOrder: "sort_order",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewCatalogTagDao creates and returns a new DAO object for table data access.
func NewCatalogTagDao(handlers ...gdb.ModelHandler) *CatalogTagDao {
	return &CatalogTagDao{
		group:    "default",
		table:    "catalog_tag",
		columns:  catalogTagColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogTagDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogTagDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogTagDao) Columns() CatalogTagColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogTagDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogTagDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogTagDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
