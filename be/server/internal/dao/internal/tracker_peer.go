// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TrackerPeerDao is the data access object for the table tracker_peer.
type TrackerPeerDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TrackerPeerColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TrackerPeerColumns defines and stores column names for the table tracker_peer.
type TrackerPeerColumns struct {
	Id             string //
	TorrentId      string //
	UserId         string //
	PeerId         string //
	Ipv4           string //
	Ipv6           string //
	Port           string //
	Uploaded       string //
	Downloaded     string //
	Remaining      string // 剩余大小
	IsSeeder       string //
	IsConnectable  string //
	Agent          string // BT 客户端
	Passkey        string //
	UploadOffset   string //
	DownloadOffset string //
	StartedAt      string //
	LastAction     string //
	FinishedAt     string //
}

// trackerPeerColumns holds the columns for the table tracker_peer.
var trackerPeerColumns = TrackerPeerColumns{
	Id:             "id",
	TorrentId:      "torrent_id",
	UserId:         "user_id",
	PeerId:         "peer_id",
	Ipv4:           "ipv4",
	Ipv6:           "ipv6",
	Port:           "port",
	Uploaded:       "uploaded",
	Downloaded:     "downloaded",
	Remaining:      "remaining",
	IsSeeder:       "is_seeder",
	IsConnectable:  "is_connectable",
	Agent:          "agent",
	Passkey:        "passkey",
	UploadOffset:   "upload_offset",
	DownloadOffset: "download_offset",
	StartedAt:      "started_at",
	LastAction:     "last_action",
	FinishedAt:     "finished_at",
}

// NewTrackerPeerDao creates and returns a new DAO object for table data access.
func NewTrackerPeerDao(handlers ...gdb.ModelHandler) *TrackerPeerDao {
	return &TrackerPeerDao{
		group:    "default",
		table:    "tracker_peer",
		columns:  trackerPeerColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TrackerPeerDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TrackerPeerDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TrackerPeerDao) Columns() TrackerPeerColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TrackerPeerDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TrackerPeerDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TrackerPeerDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
