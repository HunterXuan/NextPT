// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogSubtitleDao is the data access object for the table catalog_subtitle.
type CatalogSubtitleDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  CatalogSubtitleColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// CatalogSubtitleColumns defines and stores column names for the table catalog_subtitle.
type CatalogSubtitleColumns struct {
	Id            string //
	TorrentId     string //
	UserId        string //
	Title         string //
	FileName      string // 原始文件名
	FileExt       string //
	FileSize      string //
	StoragePath   string // 存储路径/S3 key
	Language      string // 语言代码 (如 zh-CN, en-US)
	DownloadCount string //
	Anonymous     string // 匿名上传
	CreatedAt     string //
}

// catalogSubtitleColumns holds the columns for the table catalog_subtitle.
var catalogSubtitleColumns = CatalogSubtitleColumns{
	Id:            "id",
	TorrentId:     "torrent_id",
	UserId:        "user_id",
	Title:         "title",
	FileName:      "file_name",
	FileExt:       "file_ext",
	FileSize:      "file_size",
	StoragePath:   "storage_path",
	Language:      "language",
	DownloadCount: "download_count",
	Anonymous:     "anonymous",
	CreatedAt:     "created_at",
}

// NewCatalogSubtitleDao creates and returns a new DAO object for table data access.
func NewCatalogSubtitleDao(handlers ...gdb.ModelHandler) *CatalogSubtitleDao {
	return &CatalogSubtitleDao{
		group:    "default",
		table:    "catalog_subtitle",
		columns:  catalogSubtitleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogSubtitleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogSubtitleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogSubtitleDao) Columns() CatalogSubtitleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogSubtitleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogSubtitleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogSubtitleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
