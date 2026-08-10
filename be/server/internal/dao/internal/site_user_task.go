// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SiteUserTaskDao is the data access object for the table site_user_task.
type SiteUserTaskDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  SiteUserTaskColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// SiteUserTaskColumns defines and stores column names for the table site_user_task.
type SiteUserTaskColumns struct {
	Id             string //
	UserId         string //
	TaskKey        string //
	CycleKey       string // once、YYYY-Www 或 YYYY-MM
	CycleStartedAt string //
	CycleEndedAt   string //
	Status         string // 0=active 1=completed 2=rewarded 3=expired
	Progress       string //
	Target         string //
	TaskSnapshot   string //
	ClaimedAt      string //
	CompletedAt    string //
	RewardedAt     string //
	CreatedAt      string //
	UpdatedAt      string //
}

// siteUserTaskColumns holds the columns for the table site_user_task.
var siteUserTaskColumns = SiteUserTaskColumns{
	Id:             "id",
	UserId:         "user_id",
	TaskKey:        "task_key",
	CycleKey:       "cycle_key",
	CycleStartedAt: "cycle_started_at",
	CycleEndedAt:   "cycle_ended_at",
	Status:         "status",
	Progress:       "progress",
	Target:         "target",
	TaskSnapshot:   "task_snapshot",
	ClaimedAt:      "claimed_at",
	CompletedAt:    "completed_at",
	RewardedAt:     "rewarded_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewSiteUserTaskDao creates and returns a new DAO object for table data access.
func NewSiteUserTaskDao(handlers ...gdb.ModelHandler) *SiteUserTaskDao {
	return &SiteUserTaskDao{
		group:    "default",
		table:    "site_user_task",
		columns:  siteUserTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SiteUserTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SiteUserTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SiteUserTaskDao) Columns() SiteUserTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SiteUserTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SiteUserTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SiteUserTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
