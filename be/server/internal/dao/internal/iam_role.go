// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamRoleDao is the data access object for the table iam_role.
type IamRoleDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  IamRoleColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// IamRoleColumns defines and stores column names for the table iam_role.
type IamRoleColumns struct {
	Id          string //
	Level       string // 等级权重(用于权限比对，值越大权限越高，如普通用户10，管理员100)
	NameI18N    string // 角色名称多语言映射字典
	Rules       string // 角色规则(JSON: 包含 upgrade 升级条件, keep 保级条件等)
	Permissions string // 角色关联的权限标识符列表
	IsStaff     string // 是否为管理组成员
	CreatedAt   string //
	UpdatedAt   string //
}

// iamRoleColumns holds the columns for the table iam_role.
var iamRoleColumns = IamRoleColumns{
	Id:          "id",
	Level:       "level",
	NameI18N:    "name_i18n",
	Rules:       "rules",
	Permissions: "permissions",
	IsStaff:     "is_staff",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewIamRoleDao creates and returns a new DAO object for table data access.
func NewIamRoleDao(handlers ...gdb.ModelHandler) *IamRoleDao {
	return &IamRoleDao{
		group:    "default",
		table:    "iam_role",
		columns:  iamRoleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamRoleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamRoleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamRoleDao) Columns() IamRoleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamRoleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamRoleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamRoleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
