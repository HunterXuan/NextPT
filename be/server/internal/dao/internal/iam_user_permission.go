// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamUserPermissionDao is the data access object for the table iam_user_permission.
type IamUserPermissionDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  IamUserPermissionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// IamUserPermissionColumns defines and stores column names for the table iam_user_permission.
type IamUserPermissionColumns struct {
	Id         string //
	UserId     string //
	PermKey    string // 权限标识符，如 update:torrent:123 或 update:torrent:*
	SourceType string // 来源类型: 1=manual 2=user_mod
	SourceId   string // 来源记录ID，manual=0，user_mod=mod_user_log.id
	ExpireAt   string // 权限过期时间，NULL=永久
	IsActive   string // 是否生效
	CreatedAt  string //
	UpdatedAt  string //
}

// iamUserPermissionColumns holds the columns for the table iam_user_permission.
var iamUserPermissionColumns = IamUserPermissionColumns{
	Id:         "id",
	UserId:     "user_id",
	PermKey:    "perm_key",
	SourceType: "source_type",
	SourceId:   "source_id",
	ExpireAt:   "expire_at",
	IsActive:   "is_active",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewIamUserPermissionDao creates and returns a new DAO object for table data access.
func NewIamUserPermissionDao(handlers ...gdb.ModelHandler) *IamUserPermissionDao {
	return &IamUserPermissionDao{
		group:    "default",
		table:    "iam_user_permission",
		columns:  iamUserPermissionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamUserPermissionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamUserPermissionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamUserPermissionDao) Columns() IamUserPermissionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamUserPermissionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamUserPermissionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamUserPermissionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
