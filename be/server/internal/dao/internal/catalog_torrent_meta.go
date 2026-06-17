// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogTorrentMetaDao is the data access object for the table catalog_torrent_meta.
type CatalogTorrentMetaDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  CatalogTorrentMetaColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// CatalogTorrentMetaColumns defines and stores column names for the table catalog_torrent_meta.
type CatalogTorrentMetaColumns struct {
	Id            string //
	TorrentId     string //
	ImdbId        string //
	ImdbRating    string //
	DoubanId      string //
	DoubanRating  string //
	BangumiId     string //
	BangumiRating string //
	TmdbId        string //
	TmdbRating    string //
	Extra         string // 其他元数据 (JSON)
	CreatedAt     string //
	UpdatedAt     string //
}

// catalogTorrentMetaColumns holds the columns for the table catalog_torrent_meta.
var catalogTorrentMetaColumns = CatalogTorrentMetaColumns{
	Id:            "id",
	TorrentId:     "torrent_id",
	ImdbId:        "imdb_id",
	ImdbRating:    "imdb_rating",
	DoubanId:      "douban_id",
	DoubanRating:  "douban_rating",
	BangumiId:     "bangumi_id",
	BangumiRating: "bangumi_rating",
	TmdbId:        "tmdb_id",
	TmdbRating:    "tmdb_rating",
	Extra:         "extra",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewCatalogTorrentMetaDao creates and returns a new DAO object for table data access.
func NewCatalogTorrentMetaDao(handlers ...gdb.ModelHandler) *CatalogTorrentMetaDao {
	return &CatalogTorrentMetaDao{
		group:    "default",
		table:    "catalog_torrent_meta",
		columns:  catalogTorrentMetaColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogTorrentMetaDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogTorrentMetaDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogTorrentMetaDao) Columns() CatalogTorrentMetaColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogTorrentMetaDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogTorrentMetaDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogTorrentMetaDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
