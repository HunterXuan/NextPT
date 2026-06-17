// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamInviteDao is the data access object for the table iam_invite.
type IamInviteDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  IamInviteColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// IamInviteColumns defines and stores column names for the table iam_invite.
type IamInviteColumns struct {
	Id           string //
	InviterId    string // 邀请人
	InviteeEmail string // 被邀请人邮箱
	InviteeId    string // 被邀请人ID（注册后回填）
	Hash         string // 邀请码 (发放名额时就预生成唯一码)
	Status       string // 0=未分配/待发送 1=已发送 2=已注册 3=已过期 4=已回收
	IsTemporary  string // 是否限时邀请
	ExpireAt     string //
	UsedAt       string //
	CreatedAt    string //
}

// iamInviteColumns holds the columns for the table iam_invite.
var iamInviteColumns = IamInviteColumns{
	Id:           "id",
	InviterId:    "inviter_id",
	InviteeEmail: "invitee_email",
	InviteeId:    "invitee_id",
	Hash:         "hash",
	Status:       "status",
	IsTemporary:  "is_temporary",
	ExpireAt:     "expire_at",
	UsedAt:       "used_at",
	CreatedAt:    "created_at",
}

// NewIamInviteDao creates and returns a new DAO object for table data access.
func NewIamInviteDao(handlers ...gdb.ModelHandler) *IamInviteDao {
	return &IamInviteDao{
		group:    "default",
		table:    "iam_invite",
		columns:  iamInviteColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamInviteDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamInviteDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamInviteDao) Columns() IamInviteColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamInviteDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamInviteDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamInviteDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
