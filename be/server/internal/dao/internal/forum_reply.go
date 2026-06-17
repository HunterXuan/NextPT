// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ForumReplyDao is the data access object for the table forum_reply.
type ForumReplyDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ForumReplyColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ForumReplyColumns defines and stores column names for the table forum_reply.
type ForumReplyColumns struct {
	Id          string //
	TopicId     string //
	UserId      string //
	Content     string // 回复内容
	RewardCount string // 获得的魔力值打赏
	LikeCount   string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// forumReplyColumns holds the columns for the table forum_reply.
var forumReplyColumns = ForumReplyColumns{
	Id:          "id",
	TopicId:     "topic_id",
	UserId:      "user_id",
	Content:     "content",
	RewardCount: "reward_count",
	LikeCount:   "like_count",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewForumReplyDao creates and returns a new DAO object for table data access.
func NewForumReplyDao(handlers ...gdb.ModelHandler) *ForumReplyDao {
	return &ForumReplyDao{
		group:    "default",
		table:    "forum_reply",
		columns:  forumReplyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ForumReplyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ForumReplyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ForumReplyDao) Columns() ForumReplyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ForumReplyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ForumReplyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ForumReplyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
