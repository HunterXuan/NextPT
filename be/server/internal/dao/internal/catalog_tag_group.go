// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogTagGroupDao is the data access object for the table catalog_tag_group.
type CatalogTagGroupDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  CatalogTagGroupColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// CatalogTagGroupColumns defines and stores column names for the table catalog_tag_group.
type CatalogTagGroupColumns struct {
	Id          string //
	NameI18N    string // 多语言名称映射
	Slug        string // 英文标识如 resolution
	CategoryIds string // 适用的分类 ID 数组 (例: [1,2])，为空则全站通用
	SortOrder   string //
	CreatedAt   string //
	UpdatedAt   string //
}

// catalogTagGroupColumns holds the columns for the table catalog_tag_group.
var catalogTagGroupColumns = CatalogTagGroupColumns{
	Id:          "id",
	NameI18N:    "name_i18n",
	Slug:        "slug",
	CategoryIds: "category_ids",
	SortOrder:   "sort_order",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewCatalogTagGroupDao creates and returns a new DAO object for table data access.
func NewCatalogTagGroupDao(handlers ...gdb.ModelHandler) *CatalogTagGroupDao {
	return &CatalogTagGroupDao{
		group:    "default",
		table:    "catalog_tag_group",
		columns:  catalogTagGroupColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogTagGroupDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogTagGroupDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogTagGroupDao) Columns() CatalogTagGroupColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogTagGroupDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogTagGroupDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogTagGroupDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
