// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SiteAnnouncementDao is the data access object for the table site_announcement.
type SiteAnnouncementDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  SiteAnnouncementColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// SiteAnnouncementColumns defines and stores column names for the table site_announcement.
type SiteAnnouncementColumns struct {
	Id          string //
	Title       string //
	Content     string //
	Status      string // 0=draft 1=published 2=archived
	CreatedBy   string //
	UpdatedBy   string //
	PublishedAt string //
	CreatedAt   string //
	UpdatedAt   string //
}

// siteAnnouncementColumns holds the columns for the table site_announcement.
var siteAnnouncementColumns = SiteAnnouncementColumns{
	Id:          "id",
	Title:       "title",
	Content:     "content",
	Status:      "status",
	CreatedBy:   "created_by",
	UpdatedBy:   "updated_by",
	PublishedAt: "published_at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewSiteAnnouncementDao creates and returns a new DAO object for table data access.
func NewSiteAnnouncementDao(handlers ...gdb.ModelHandler) *SiteAnnouncementDao {
	return &SiteAnnouncementDao{
		group:    "default",
		table:    "site_announcement",
		columns:  siteAnnouncementColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SiteAnnouncementDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SiteAnnouncementDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SiteAnnouncementDao) Columns() SiteAnnouncementColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SiteAnnouncementDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SiteAnnouncementDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SiteAnnouncementDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
