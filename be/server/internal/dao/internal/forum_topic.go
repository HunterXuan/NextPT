// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ForumTopicDao is the data access object for the table forum_topic.
type ForumTopicDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ForumTopicColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ForumTopicColumns defines and stores column names for the table forum_topic.
type ForumTopicColumns struct {
	Id          string //
	NodeId      string //
	UserId      string //
	Subject     string //
	Content     string // 主题正文
	Appends     string // 追加内容 (V2EX风格) [{"content": "...", "created_at": "..."}]
	IsLocked    string //
	IsSticky    string //
	Views       string //
	ReplyCount  string //
	LikeCount   string //
	LastReplyId string //
	LastReplyAt string //
	LastReplyBy string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// forumTopicColumns holds the columns for the table forum_topic.
var forumTopicColumns = ForumTopicColumns{
	Id:          "id",
	NodeId:      "node_id",
	UserId:      "user_id",
	Subject:     "subject",
	Content:     "content",
	Appends:     "appends",
	IsLocked:    "is_locked",
	IsSticky:    "is_sticky",
	Views:       "views",
	ReplyCount:  "reply_count",
	LikeCount:   "like_count",
	LastReplyId: "last_reply_id",
	LastReplyAt: "last_reply_at",
	LastReplyBy: "last_reply_by",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewForumTopicDao creates and returns a new DAO object for table data access.
func NewForumTopicDao(handlers ...gdb.ModelHandler) *ForumTopicDao {
	return &ForumTopicDao{
		group:    "default",
		table:    "forum_topic",
		columns:  forumTopicColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ForumTopicDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ForumTopicDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ForumTopicDao) Columns() ForumTopicColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ForumTopicDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ForumTopicDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ForumTopicDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
