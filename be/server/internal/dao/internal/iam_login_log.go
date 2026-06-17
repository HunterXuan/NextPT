// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamLoginLogDao is the data access object for the table iam_login_log.
type IamLoginLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  IamLoginLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// IamLoginLogColumns defines and stores column names for the table iam_login_log.
type IamLoginLogColumns struct {
	Id         string //
	UserId     string //
	Ip         string //
	UserAgent  string //
	Result     string // 1=success 0=fail
	FailReason string //
	CreatedAt  string //
}

// iamLoginLogColumns holds the columns for the table iam_login_log.
var iamLoginLogColumns = IamLoginLogColumns{
	Id:         "id",
	UserId:     "user_id",
	Ip:         "ip",
	UserAgent:  "user_agent",
	Result:     "result",
	FailReason: "fail_reason",
	CreatedAt:  "created_at",
}

// NewIamLoginLogDao creates and returns a new DAO object for table data access.
func NewIamLoginLogDao(handlers ...gdb.ModelHandler) *IamLoginLogDao {
	return &IamLoginLogDao{
		group:    "default",
		table:    "iam_login_log",
		columns:  iamLoginLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamLoginLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamLoginLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamLoginLogDao) Columns() IamLoginLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamLoginLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamLoginLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamLoginLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
