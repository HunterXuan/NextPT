// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogRequestDao is the data access object for the table catalog_request.
type CatalogRequestDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  CatalogRequestColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// CatalogRequestColumns defines and stores column names for the table catalog_request.
type CatalogRequestColumns struct {
	Id              string //
	RequestType     string // 1=求种 2=续种
	RequesterId     string //
	CategoryId      string //
	TargetTorrentId string // 续种目标种子
	ResultTorrentId string // 求种完成后关联的种子
	Title           string //
	Description     string // 请求说明 (Markdown/BBCode)
	RewardAmount    string // 已托管的魔力奖励
	Status          string // 0=开放 1=已认领 2=待确认 3=已完成 4=已取消
	ClaimedBy       string //
	ClaimedAt       string //
	ClaimExpiresAt  string //
	SubmittedAt     string //
	CompletedAt     string //
	CancelledBy     string //
	CancelledAt     string //
	CancelReason    string //
	CreatedAt       string //
	UpdatedAt       string //
}

// catalogRequestColumns holds the columns for the table catalog_request.
var catalogRequestColumns = CatalogRequestColumns{
	Id:              "id",
	RequestType:     "request_type",
	RequesterId:     "requester_id",
	CategoryId:      "category_id",
	TargetTorrentId: "target_torrent_id",
	ResultTorrentId: "result_torrent_id",
	Title:           "title",
	Description:     "description",
	RewardAmount:    "reward_amount",
	Status:          "status",
	ClaimedBy:       "claimed_by",
	ClaimedAt:       "claimed_at",
	ClaimExpiresAt:  "claim_expires_at",
	SubmittedAt:     "submitted_at",
	CompletedAt:     "completed_at",
	CancelledBy:     "cancelled_by",
	CancelledAt:     "cancelled_at",
	CancelReason:    "cancel_reason",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewCatalogRequestDao creates and returns a new DAO object for table data access.
func NewCatalogRequestDao(handlers ...gdb.ModelHandler) *CatalogRequestDao {
	return &CatalogRequestDao{
		group:    "default",
		table:    "catalog_request",
		columns:  catalogRequestColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogRequestDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogRequestDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogRequestDao) Columns() CatalogRequestColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogRequestDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogRequestDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogRequestDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
