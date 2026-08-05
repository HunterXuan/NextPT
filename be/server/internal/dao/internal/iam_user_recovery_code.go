// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamUserRecoveryCodeDao is the data access object for the table iam_user_recovery_code.
type IamUserRecoveryCodeDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  IamUserRecoveryCodeColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// IamUserRecoveryCodeColumns defines and stores column names for the table iam_user_recovery_code.
type IamUserRecoveryCodeColumns struct {
	Id        string //
	UserId    string //
	CodeHash  string //
	UsedAt    string //
	CreatedAt string //
}

// iamUserRecoveryCodeColumns holds the columns for the table iam_user_recovery_code.
var iamUserRecoveryCodeColumns = IamUserRecoveryCodeColumns{
	Id:        "id",
	UserId:    "user_id",
	CodeHash:  "code_hash",
	UsedAt:    "used_at",
	CreatedAt: "created_at",
}

// NewIamUserRecoveryCodeDao creates and returns a new DAO object for table data access.
func NewIamUserRecoveryCodeDao(handlers ...gdb.ModelHandler) *IamUserRecoveryCodeDao {
	return &IamUserRecoveryCodeDao{
		group:    "default",
		table:    "iam_user_recovery_code",
		columns:  iamUserRecoveryCodeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamUserRecoveryCodeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamUserRecoveryCodeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamUserRecoveryCodeDao) Columns() IamUserRecoveryCodeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamUserRecoveryCodeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamUserRecoveryCodeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamUserRecoveryCodeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
