// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamUserProfileDao is the data access object for the table iam_user_profile.
type IamUserProfileDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  IamUserProfileColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// IamUserProfileColumns defines and stores column names for the table iam_user_profile.
type IamUserProfileColumns struct {
	Id        string //
	UserId    string //
	Avatar    string //
	Info      string // 个人简介 (Markdown)
	Signature string // 论坛签名
	CreatedAt string //
	UpdatedAt string //
}

// iamUserProfileColumns holds the columns for the table iam_user_profile.
var iamUserProfileColumns = IamUserProfileColumns{
	Id:        "id",
	UserId:    "user_id",
	Avatar:    "avatar",
	Info:      "info",
	Signature: "signature",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewIamUserProfileDao creates and returns a new DAO object for table data access.
func NewIamUserProfileDao(handlers ...gdb.ModelHandler) *IamUserProfileDao {
	return &IamUserProfileDao{
		group:    "default",
		table:    "iam_user_profile",
		columns:  iamUserProfileColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamUserProfileDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamUserProfileDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamUserProfileDao) Columns() IamUserProfileColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamUserProfileDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamUserProfileDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamUserProfileDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
