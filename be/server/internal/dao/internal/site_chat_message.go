// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SiteChatMessageDao is the data access object for the table site_chat_message.
type SiteChatMessageDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  SiteChatMessageColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// SiteChatMessageColumns defines and stores column names for the table site_chat_message.
type SiteChatMessageColumns struct {
	Id        string //
	UserId    string //
	Content   string //
	CreatedAt string //
}

// siteChatMessageColumns holds the columns for the table site_chat_message.
var siteChatMessageColumns = SiteChatMessageColumns{
	Id:        "id",
	UserId:    "user_id",
	Content:   "content",
	CreatedAt: "created_at",
}

// NewSiteChatMessageDao creates and returns a new DAO object for table data access.
func NewSiteChatMessageDao(handlers ...gdb.ModelHandler) *SiteChatMessageDao {
	return &SiteChatMessageDao{
		group:    "default",
		table:    "site_chat_message",
		columns:  siteChatMessageColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SiteChatMessageDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SiteChatMessageDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SiteChatMessageDao) Columns() SiteChatMessageColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SiteChatMessageDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SiteChatMessageDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SiteChatMessageDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
