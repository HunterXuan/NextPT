// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ForumCategoryDao is the data access object for the table forum_category.
type ForumCategoryDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  ForumCategoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// ForumCategoryColumns defines and stores column names for the table forum_category.
type ForumCategoryColumns struct {
	Id          string //
	NameI18N    string // 多语言名称映射
	DescI18N    string // 多语言描述映射
	SortOrder   string //
	MinRoleView string // 最低可见等级
	CreatedAt   string //
	UpdatedAt   string //
}

// forumCategoryColumns holds the columns for the table forum_category.
var forumCategoryColumns = ForumCategoryColumns{
	Id:          "id",
	NameI18N:    "name_i18n",
	DescI18N:    "desc_i18n",
	SortOrder:   "sort_order",
	MinRoleView: "min_role_view",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewForumCategoryDao creates and returns a new DAO object for table data access.
func NewForumCategoryDao(handlers ...gdb.ModelHandler) *ForumCategoryDao {
	return &ForumCategoryDao{
		group:    "default",
		table:    "forum_category",
		columns:  forumCategoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ForumCategoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ForumCategoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ForumCategoryDao) Columns() ForumCategoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ForumCategoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ForumCategoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ForumCategoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
