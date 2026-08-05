// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamUserDao is the data access object for the table iam_user.
type IamUserDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  IamUserColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// IamUserColumns defines and stores column names for the table iam_user.
type IamUserColumns struct {
	Id            string //
	Username      string //
	Email         string //
	PasswordHash  string // bcrypt hash
	Passkey       string // Tracker passkey
	Status        string // 0=pending 1=confirmed 2=disabled
	Role          string // 当前角色ID (关联 user_role 表)
	VipUntil      string // VIP 过期时间
	VipRemark     string // VIP 身份获取备注/来源
	TwoStepType   string // 两步验证方式: 0=关闭 1=TOTP(Authenticator) 2=邮件验证码
	TwoStepSecret string // 加密后的 TOTP 密钥 (two_step_type=1 时使用)
	InvitedBy     string //
	LastLogin     string // 最后登录时间
	LastIp        string // 最后登录 IP
	CreatedAt     string //
	UpdatedAt     string //
	DeletedAt     string //
}

// iamUserColumns holds the columns for the table iam_user.
var iamUserColumns = IamUserColumns{
	Id:            "id",
	Username:      "username",
	Email:         "email",
	PasswordHash:  "password_hash",
	Passkey:       "passkey",
	Status:        "status",
	Role:          "role",
	VipUntil:      "vip_until",
	VipRemark:     "vip_remark",
	TwoStepType:   "two_step_type",
	TwoStepSecret: "two_step_secret",
	InvitedBy:     "invited_by",
	LastLogin:     "last_login",
	LastIp:        "last_ip",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewIamUserDao creates and returns a new DAO object for table data access.
func NewIamUserDao(handlers ...gdb.ModelHandler) *IamUserDao {
	return &IamUserDao{
		group:    "default",
		table:    "iam_user",
		columns:  iamUserColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamUserDao) Columns() IamUserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamUserDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
