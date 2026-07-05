// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogCommentDao is the data access object for the table catalog_comment.
type CatalogCommentDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  CatalogCommentColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// CatalogCommentColumns defines and stores column names for the table catalog_comment.
type CatalogCommentColumns struct {
	Id          string //
	TargetType  string // catalog_torrent
	TargetId    string //
	UserId      string //
	Content     string // 评论内容 (Markdown/BBCode)
	RewardCount string // 获得的魔力值打赏
	LikeCount   string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// catalogCommentColumns holds the columns for the table catalog_comment.
var catalogCommentColumns = CatalogCommentColumns{
	Id:          "id",
	TargetType:  "target_type",
	TargetId:    "target_id",
	UserId:      "user_id",
	Content:     "content",
	RewardCount: "reward_count",
	LikeCount:   "like_count",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewCatalogCommentDao creates and returns a new DAO object for table data access.
func NewCatalogCommentDao(handlers ...gdb.ModelHandler) *CatalogCommentDao {
	return &CatalogCommentDao{
		group:    "default",
		table:    "catalog_comment",
		columns:  catalogCommentColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogCommentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogCommentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogCommentDao) Columns() CatalogCommentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogCommentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogCommentDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogCommentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
