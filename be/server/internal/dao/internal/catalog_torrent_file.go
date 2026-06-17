// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogTorrentFileDao is the data access object for the table catalog_torrent_file.
type CatalogTorrentFileDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  CatalogTorrentFileColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// CatalogTorrentFileColumns defines and stores column names for the table catalog_torrent_file.
type CatalogTorrentFileColumns struct {
	Id        string //
	TorrentId string //
	FilePath  string // 文件路径
	Size      string //
}

// catalogTorrentFileColumns holds the columns for the table catalog_torrent_file.
var catalogTorrentFileColumns = CatalogTorrentFileColumns{
	Id:        "id",
	TorrentId: "torrent_id",
	FilePath:  "file_path",
	Size:      "size",
}

// NewCatalogTorrentFileDao creates and returns a new DAO object for table data access.
func NewCatalogTorrentFileDao(handlers ...gdb.ModelHandler) *CatalogTorrentFileDao {
	return &CatalogTorrentFileDao{
		group:    "default",
		table:    "catalog_torrent_file",
		columns:  catalogTorrentFileColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogTorrentFileDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogTorrentFileDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogTorrentFileDao) Columns() CatalogTorrentFileColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogTorrentFileDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogTorrentFileDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogTorrentFileDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
