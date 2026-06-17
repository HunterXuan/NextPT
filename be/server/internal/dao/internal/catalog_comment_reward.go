// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogCommentRewardDao is the data access object for the table catalog_comment_reward.
type CatalogCommentRewardDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  CatalogCommentRewardColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// CatalogCommentRewardColumns defines and stores column names for the table catalog_comment_reward.
type CatalogCommentRewardColumns struct {
	Id        string //
	UserId    string //
	CommentId string //
	Amount    string // 打赏金额
	CreatedAt string //
}

// catalogCommentRewardColumns holds the columns for the table catalog_comment_reward.
var catalogCommentRewardColumns = CatalogCommentRewardColumns{
	Id:        "id",
	UserId:    "user_id",
	CommentId: "comment_id",
	Amount:    "amount",
	CreatedAt: "created_at",
}

// NewCatalogCommentRewardDao creates and returns a new DAO object for table data access.
func NewCatalogCommentRewardDao(handlers ...gdb.ModelHandler) *CatalogCommentRewardDao {
	return &CatalogCommentRewardDao{
		group:    "default",
		table:    "catalog_comment_reward",
		columns:  catalogCommentRewardColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogCommentRewardDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogCommentRewardDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogCommentRewardDao) Columns() CatalogCommentRewardColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogCommentRewardDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogCommentRewardDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogCommentRewardDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
