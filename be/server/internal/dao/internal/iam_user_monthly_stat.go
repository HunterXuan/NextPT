// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamUserMonthlyStatDao is the data access object for the table iam_user_monthly_stat.
type IamUserMonthlyStatDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  IamUserMonthlyStatColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// IamUserMonthlyStatColumns defines and stores column names for the table iam_user_monthly_stat.
type IamUserMonthlyStatColumns struct {
	Id         string //
	UserId     string //
	YearMonth  string // YYYY-MM
	Uploaded   string // 当月新增上传量 (bytes)
	Downloaded string // 当月新增下载量 (bytes)
	SeedTime   string // 当月新增做种时间 (秒)
	LeechTime  string // 当月新增下载时间 (秒)
	Bonus      string // 当月获得魔力值
	CreatedAt  string //
}

// iamUserMonthlyStatColumns holds the columns for the table iam_user_monthly_stat.
var iamUserMonthlyStatColumns = IamUserMonthlyStatColumns{
	Id:         "id",
	UserId:     "user_id",
	YearMonth:  "year_month",
	Uploaded:   "uploaded",
	Downloaded: "downloaded",
	SeedTime:   "seed_time",
	LeechTime:  "leech_time",
	Bonus:      "bonus",
	CreatedAt:  "created_at",
}

// NewIamUserMonthlyStatDao creates and returns a new DAO object for table data access.
func NewIamUserMonthlyStatDao(handlers ...gdb.ModelHandler) *IamUserMonthlyStatDao {
	return &IamUserMonthlyStatDao{
		group:    "default",
		table:    "iam_user_monthly_stat",
		columns:  iamUserMonthlyStatColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamUserMonthlyStatDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamUserMonthlyStatDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamUserMonthlyStatDao) Columns() IamUserMonthlyStatColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamUserMonthlyStatDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamUserMonthlyStatDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamUserMonthlyStatDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
