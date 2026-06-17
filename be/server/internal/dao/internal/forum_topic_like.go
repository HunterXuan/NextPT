// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ForumTopicLikeDao is the data access object for the table forum_topic_like.
type ForumTopicLikeDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  ForumTopicLikeColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// ForumTopicLikeColumns defines and stores column names for the table forum_topic_like.
type ForumTopicLikeColumns struct {
	Id        string //
	UserId    string //
	TopicId   string //
	CreatedAt string //
}

// forumTopicLikeColumns holds the columns for the table forum_topic_like.
var forumTopicLikeColumns = ForumTopicLikeColumns{
	Id:        "id",
	UserId:    "user_id",
	TopicId:   "topic_id",
	CreatedAt: "created_at",
}

// NewForumTopicLikeDao creates and returns a new DAO object for table data access.
func NewForumTopicLikeDao(handlers ...gdb.ModelHandler) *ForumTopicLikeDao {
	return &ForumTopicLikeDao{
		group:    "default",
		table:    "forum_topic_like",
		columns:  forumTopicLikeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ForumTopicLikeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ForumTopicLikeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ForumTopicLikeDao) Columns() ForumTopicLikeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ForumTopicLikeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ForumTopicLikeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ForumTopicLikeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
