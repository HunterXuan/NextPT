// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogTorrentDao is the data access object for the table catalog_torrent.
type CatalogTorrentDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  CatalogTorrentColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// CatalogTorrentColumns defines and stores column names for the table catalog_torrent.
type CatalogTorrentColumns struct {
	Id             string //
	InfoHash       string // BitTorrent info_hash
	Name           string // 种子标题
	SubTitle       string // 副标题
	CategoryId     string //
	Description    string // 详情描述 (BBCode/Markdown)
	Nfo            string // NFO 文件内容
	FileName       string // 种子文件名
	Size           string // 总大小 (bytes)
	FileCount      string // 文件数量
	Type           string // 0=single 1=multi
	OwnerId        string // 上传者
	Anonymous      string // 匿名上传
	SpState        string // 0=normal 1=free 2=2x 3=2xfree 4=50%off 5=2x50% 6=30%off 7=custom
	SpExpireAt     string // 促销到期时间
	IsFeatured     string // 是否推荐
	IsPinned       string // 是否置顶
	PinWeight      string // 置顶权重 (越大越靠前)
	Seeders        string //
	Leechers       string //
	TimesCompleted string //
	CommentsCount  string //
	Views          string //
	LikeCount      string //
	RewardsCount   string // 收到的赞赏次数
	RewardsAmount  string // 收到的赞赏总金额(Bonus)
	Visible        string //
	Banned         string //
	LastAction     string // Tracker 最后活动时间
	LastReseed     string //
	CreatedAt      string //
	UpdatedAt      string //
}

// catalogTorrentColumns holds the columns for the table catalog_torrent.
var catalogTorrentColumns = CatalogTorrentColumns{
	Id:             "id",
	InfoHash:       "info_hash",
	Name:           "name",
	SubTitle:       "sub_title",
	CategoryId:     "category_id",
	Description:    "description",
	Nfo:            "nfo",
	FileName:       "file_name",
	Size:           "size",
	FileCount:      "file_count",
	Type:           "type",
	OwnerId:        "owner_id",
	Anonymous:      "anonymous",
	SpState:        "sp_state",
	SpExpireAt:     "sp_expire_at",
	IsFeatured:     "is_featured",
	IsPinned:       "is_pinned",
	PinWeight:      "pin_weight",
	Seeders:        "seeders",
	Leechers:       "leechers",
	TimesCompleted: "times_completed",
	CommentsCount:  "comments_count",
	Views:          "views",
	LikeCount:      "like_count",
	RewardsCount:   "rewards_count",
	RewardsAmount:  "rewards_amount",
	Visible:        "visible",
	Banned:         "banned",
	LastAction:     "last_action",
	LastReseed:     "last_reseed",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewCatalogTorrentDao creates and returns a new DAO object for table data access.
func NewCatalogTorrentDao(handlers ...gdb.ModelHandler) *CatalogTorrentDao {
	return &CatalogTorrentDao{
		group:    "default",
		table:    "catalog_torrent",
		columns:  catalogTorrentColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogTorrentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogTorrentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogTorrentDao) Columns() CatalogTorrentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogTorrentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogTorrentDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogTorrentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
