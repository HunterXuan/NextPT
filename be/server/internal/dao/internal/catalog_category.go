// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogCategoryDao is the data access object for the table catalog_category.
type CatalogCategoryDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  CatalogCategoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// CatalogCategoryColumns defines and stores column names for the table catalog_category.
type CatalogCategoryColumns struct {
	Id        string //
	NameI18N  string // 多语言名称映射
	Slug      string // URL-friendly
	SortOrder string //
	Enabled   string //
	CreatedAt string //
	UpdatedAt string //
}

// catalogCategoryColumns holds the columns for the table catalog_category.
var catalogCategoryColumns = CatalogCategoryColumns{
	Id:        "id",
	NameI18N:  "name_i18n",
	Slug:      "slug",
	SortOrder: "sort_order",
	Enabled:   "enabled",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewCatalogCategoryDao creates and returns a new DAO object for table data access.
func NewCatalogCategoryDao(handlers ...gdb.ModelHandler) *CatalogCategoryDao {
	return &CatalogCategoryDao{
		group:    "default",
		table:    "catalog_category",
		columns:  catalogCategoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogCategoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogCategoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogCategoryDao) Columns() CatalogCategoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogCategoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogCategoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogCategoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
