// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TrackerAgentWhitelistDao is the data access object for the table tracker_agent_whitelist.
type TrackerAgentWhitelistDao struct {
	table    string                       // table is the underlying table name of the DAO.
	group    string                       // group is the database configuration group name of the current DAO.
	columns  TrackerAgentWhitelistColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler           // handlers for customized model modification.
}

// TrackerAgentWhitelistColumns defines and stores column names for the table tracker_agent_whitelist.
type TrackerAgentWhitelistColumns struct {
	Id           string //
	Family       string // 客户端家族名
	PeerIdPrefix string //
	AgentPattern string //
	MinVersion   string //
	MaxVersion   string //
	AllowHttps   string //
	Enabled      string //
	Comment      string //
	Hits         string //
	CreatedAt    string //
	UpdatedAt    string //
}

// trackerAgentWhitelistColumns holds the columns for the table tracker_agent_whitelist.
var trackerAgentWhitelistColumns = TrackerAgentWhitelistColumns{
	Id:           "id",
	Family:       "family",
	PeerIdPrefix: "peer_id_prefix",
	AgentPattern: "agent_pattern",
	MinVersion:   "min_version",
	MaxVersion:   "max_version",
	AllowHttps:   "allow_https",
	Enabled:      "enabled",
	Comment:      "comment",
	Hits:         "hits",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewTrackerAgentWhitelistDao creates and returns a new DAO object for table data access.
func NewTrackerAgentWhitelistDao(handlers ...gdb.ModelHandler) *TrackerAgentWhitelistDao {
	return &TrackerAgentWhitelistDao{
		group:    "default",
		table:    "tracker_agent_whitelist",
		columns:  trackerAgentWhitelistColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TrackerAgentWhitelistDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TrackerAgentWhitelistDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TrackerAgentWhitelistDao) Columns() TrackerAgentWhitelistColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TrackerAgentWhitelistDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TrackerAgentWhitelistDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TrackerAgentWhitelistDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
