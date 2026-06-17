// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ForumNodeDao is the data access object for the table forum_node.
type ForumNodeDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ForumNodeColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ForumNodeColumns defines and stores column names for the table forum_node.
type ForumNodeColumns struct {
	Id            string //
	CategoryId    string // 所属节点分类
	Slug          string // 节点英文标识符
	NameI18N      string // 多语言名称映射
	DescI18N      string // 多语言描述映射
	SortOrder     string //
	MinRoleRead   string //
	MinRoleWrite  string //
	MinRoleCreate string //
	TopicCount    string //
	ReplyCount    string //
	LastTopicId   string //
	Moderators    string // 板块版主ID列表(前端展示用)
	LastReplyAt   string //
	CreatedAt     string //
	UpdatedAt     string //
}

// forumNodeColumns holds the columns for the table forum_node.
var forumNodeColumns = ForumNodeColumns{
	Id:            "id",
	CategoryId:    "category_id",
	Slug:          "slug",
	NameI18N:      "name_i18n",
	DescI18N:      "desc_i18n",
	SortOrder:     "sort_order",
	MinRoleRead:   "min_role_read",
	MinRoleWrite:  "min_role_write",
	MinRoleCreate: "min_role_create",
	TopicCount:    "topic_count",
	ReplyCount:    "reply_count",
	LastTopicId:   "last_topic_id",
	Moderators:    "moderators",
	LastReplyAt:   "last_reply_at",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewForumNodeDao creates and returns a new DAO object for table data access.
func NewForumNodeDao(handlers ...gdb.ModelHandler) *ForumNodeDao {
	return &ForumNodeDao{
		group:    "default",
		table:    "forum_node",
		columns:  forumNodeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ForumNodeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ForumNodeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ForumNodeDao) Columns() ForumNodeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ForumNodeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ForumNodeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ForumNodeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
