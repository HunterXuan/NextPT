// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SiteMessageDao is the data access object for the table site_message.
type SiteMessageDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SiteMessageColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SiteMessageColumns defines and stores column names for the table site_message.
type SiteMessageColumns struct {
	Id         string //
	SenderId   string // 0=系统通知, 或管理员ID
	ReceiverId string //
	Content    string //
	IsRead     string //
	CreatedAt  string //
}

// siteMessageColumns holds the columns for the table site_message.
var siteMessageColumns = SiteMessageColumns{
	Id:         "id",
	SenderId:   "sender_id",
	ReceiverId: "receiver_id",
	Content:    "content",
	IsRead:     "is_read",
	CreatedAt:  "created_at",
}

// NewSiteMessageDao creates and returns a new DAO object for table data access.
func NewSiteMessageDao(handlers ...gdb.ModelHandler) *SiteMessageDao {
	return &SiteMessageDao{
		group:    "default",
		table:    "site_message",
		columns:  siteMessageColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SiteMessageDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SiteMessageDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SiteMessageDao) Columns() SiteMessageColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SiteMessageDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SiteMessageDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SiteMessageDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
