// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ModReportDao is the data access object for the table mod_report.
type ModReportDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ModReportColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ModReportColumns defines and stores column names for the table mod_report.
type ModReportColumns struct {
	Id           string //
	ReporterId   string //
	TargetType   string // torrent/user/offer/request/forum_post/comment/subtitle
	TargetId     string //
	Reason       string //
	Status       string // 0=pending 1=resolved 2=rejected
	DealtBy      string //
	DealtComment string //
	DealtAt      string //
	CreatedAt    string //
}

// modReportColumns holds the columns for the table mod_report.
var modReportColumns = ModReportColumns{
	Id:           "id",
	ReporterId:   "reporter_id",
	TargetType:   "target_type",
	TargetId:     "target_id",
	Reason:       "reason",
	Status:       "status",
	DealtBy:      "dealt_by",
	DealtComment: "dealt_comment",
	DealtAt:      "dealt_at",
	CreatedAt:    "created_at",
}

// NewModReportDao creates and returns a new DAO object for table data access.
func NewModReportDao(handlers ...gdb.ModelHandler) *ModReportDao {
	return &ModReportDao{
		group:    "default",
		table:    "mod_report",
		columns:  modReportColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ModReportDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ModReportDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ModReportDao) Columns() ModReportColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ModReportDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ModReportDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ModReportDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
