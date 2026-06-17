// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysCronLogDao is the data access object for the table sys_cron_log.
type SysCronLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysCronLogColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysCronLogColumns defines and stores column names for the table sys_cron_log.
type SysCronLogColumns struct {
	Id           string //
	JobName      string // 任务名称
	NodeIp       string // 执行节点IP
	Status       string // 状态 0: 执行中 1: 成功 2: 失败
	DurationMs   string // 执行耗时(毫秒)
	ErrorMessage string // 错误信息
	CreatedAt    string // 开始时间
	UpdatedAt    string // 更新时间
}

// sysCronLogColumns holds the columns for the table sys_cron_log.
var sysCronLogColumns = SysCronLogColumns{
	Id:           "id",
	JobName:      "job_name",
	NodeIp:       "node_ip",
	Status:       "status",
	DurationMs:   "duration_ms",
	ErrorMessage: "error_message",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewSysCronLogDao creates and returns a new DAO object for table data access.
func NewSysCronLogDao(handlers ...gdb.ModelHandler) *SysCronLogDao {
	return &SysCronLogDao{
		group:    "default",
		table:    "sys_cron_log",
		columns:  sysCronLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysCronLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysCronLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysCronLogDao) Columns() SysCronLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysCronLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysCronLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysCronLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
