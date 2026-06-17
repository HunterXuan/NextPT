// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ModStaffMessageDao is the data access object for the table mod_staff_message.
type ModStaffMessageDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  ModStaffMessageColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// ModStaffMessageColumns defines and stores column names for the table mod_staff_message.
type ModStaffMessageColumns struct {
	Id         string //
	SenderId   string //
	Subject    string //
	Content    string //
	Status     string // 0=pending 1=answered 2=closed
	AnsweredBy string //
	Answer     string //
	AnsweredAt string //
	CreatedAt  string //
	UpdatedAt  string //
}

// modStaffMessageColumns holds the columns for the table mod_staff_message.
var modStaffMessageColumns = ModStaffMessageColumns{
	Id:         "id",
	SenderId:   "sender_id",
	Subject:    "subject",
	Content:    "content",
	Status:     "status",
	AnsweredBy: "answered_by",
	Answer:     "answer",
	AnsweredAt: "answered_at",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewModStaffMessageDao creates and returns a new DAO object for table data access.
func NewModStaffMessageDao(handlers ...gdb.ModelHandler) *ModStaffMessageDao {
	return &ModStaffMessageDao{
		group:    "default",
		table:    "mod_staff_message",
		columns:  modStaffMessageColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ModStaffMessageDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ModStaffMessageDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ModStaffMessageDao) Columns() ModStaffMessageColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ModStaffMessageDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ModStaffMessageDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ModStaffMessageDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
