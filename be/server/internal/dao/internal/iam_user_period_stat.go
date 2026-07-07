// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamUserPeriodStatDao is the data access object for the table iam_user_period_stat.
type IamUserPeriodStatDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  IamUserPeriodStatColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// IamUserPeriodStatColumns defines and stores column names for the table iam_user_period_stat.
type IamUserPeriodStatColumns struct {
	Id            string //
	UserId        string //
	PeriodType    string // 1=每日 2=每月
	PeriodKey     string // YYYY-MM-DD 或 YYYY-MM
	Uploaded      string // 周期新增入账上传量 (bytes)
	Downloaded    string // 周期新增入账下载量 (bytes)
	RawUploaded   string // 周期新增真实上传量 (bytes)
	RawDownloaded string // 周期新增真实下载量 (bytes)
	SeedTime      string // 周期新增做种时间 (秒)
	LeechTime     string // 周期新增下载时间 (秒)
	Bonus         string // 周期获得魔力值
	CreatedAt     string //
}

// iamUserPeriodStatColumns holds the columns for the table iam_user_period_stat.
var iamUserPeriodStatColumns = IamUserPeriodStatColumns{
	Id:            "id",
	UserId:        "user_id",
	PeriodType:    "period_type",
	PeriodKey:     "period_key",
	Uploaded:      "uploaded",
	Downloaded:    "downloaded",
	RawUploaded:   "raw_uploaded",
	RawDownloaded: "raw_downloaded",
	SeedTime:      "seed_time",
	LeechTime:     "leech_time",
	Bonus:         "bonus",
	CreatedAt:     "created_at",
}

// NewIamUserPeriodStatDao creates and returns a new DAO object for table data access.
func NewIamUserPeriodStatDao(handlers ...gdb.ModelHandler) *IamUserPeriodStatDao {
	return &IamUserPeriodStatDao{
		group:    "default",
		table:    "iam_user_period_stat",
		columns:  iamUserPeriodStatColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamUserPeriodStatDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamUserPeriodStatDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamUserPeriodStatDao) Columns() IamUserPeriodStatColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamUserPeriodStatDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamUserPeriodStatDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamUserPeriodStatDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
