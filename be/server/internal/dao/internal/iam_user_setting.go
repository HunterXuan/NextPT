// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamUserSettingDao is the data access object for the table iam_user_setting.
type IamUserSettingDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  IamUserSettingColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// IamUserSettingColumns defines and stores column names for the table iam_user_setting.
type IamUserSettingColumns struct {
	Id           string //
	UserId       string //
	PrivacyLevel string // 0=宽松 1=普通 2=严格
	Extra        string // UI偏好/分页/语言/时区等前端设置 (JSON 扩展)
	CreatedAt    string //
	UpdatedAt    string //
}

// iamUserSettingColumns holds the columns for the table iam_user_setting.
var iamUserSettingColumns = IamUserSettingColumns{
	Id:           "id",
	UserId:       "user_id",
	PrivacyLevel: "privacy_level",
	Extra:        "extra",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewIamUserSettingDao creates and returns a new DAO object for table data access.
func NewIamUserSettingDao(handlers ...gdb.ModelHandler) *IamUserSettingDao {
	return &IamUserSettingDao{
		group:    "default",
		table:    "iam_user_setting",
		columns:  iamUserSettingColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamUserSettingDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamUserSettingDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamUserSettingDao) Columns() IamUserSettingColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamUserSettingDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamUserSettingDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamUserSettingDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
