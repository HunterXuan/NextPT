// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SiteAuditDao is the data access object for the table site_audit.
type SiteAuditDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SiteAuditColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SiteAuditColumns defines and stores column names for the table site_audit.
type SiteAuditColumns struct {
	Id         string //
	UserId     string //
	Action     string //
	TargetType string //
	TargetId   string //
	Detail     string //
	Ip         string //
	Level      string // 0=normal 1=mod 2=admin
	CreatedAt  string //
}

// siteAuditColumns holds the columns for the table site_audit.
var siteAuditColumns = SiteAuditColumns{
	Id:         "id",
	UserId:     "user_id",
	Action:     "action",
	TargetType: "target_type",
	TargetId:   "target_id",
	Detail:     "detail",
	Ip:         "ip",
	Level:      "level",
	CreatedAt:  "created_at",
}

// NewSiteAuditDao creates and returns a new DAO object for table data access.
func NewSiteAuditDao(handlers ...gdb.ModelHandler) *SiteAuditDao {
	return &SiteAuditDao{
		group:    "default",
		table:    "site_audit",
		columns:  siteAuditColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SiteAuditDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SiteAuditDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SiteAuditDao) Columns() SiteAuditColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SiteAuditDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SiteAuditDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SiteAuditDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
