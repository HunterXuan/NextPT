// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogTorrentLikeDao is the data access object for the table catalog_torrent_like.
type CatalogTorrentLikeDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  CatalogTorrentLikeColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// CatalogTorrentLikeColumns defines and stores column names for the table catalog_torrent_like.
type CatalogTorrentLikeColumns struct {
	Id        string //
	UserId    string //
	TorrentId string //
	CreatedAt string //
}

// catalogTorrentLikeColumns holds the columns for the table catalog_torrent_like.
var catalogTorrentLikeColumns = CatalogTorrentLikeColumns{
	Id:        "id",
	UserId:    "user_id",
	TorrentId: "torrent_id",
	CreatedAt: "created_at",
}

// NewCatalogTorrentLikeDao creates and returns a new DAO object for table data access.
func NewCatalogTorrentLikeDao(handlers ...gdb.ModelHandler) *CatalogTorrentLikeDao {
	return &CatalogTorrentLikeDao{
		group:    "default",
		table:    "catalog_torrent_like",
		columns:  catalogTorrentLikeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogTorrentLikeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogTorrentLikeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogTorrentLikeDao) Columns() CatalogTorrentLikeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogTorrentLikeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogTorrentLikeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogTorrentLikeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
