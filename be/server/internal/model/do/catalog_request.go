// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogRequest is the golang structure of table catalog_request for DAO operations like Where/Data.
type CatalogRequest struct {
	g.Meta          `orm:"table:catalog_request, do:true"`
	Id              any         //
	RequestType     any         // 1=求种 2=续种
	RequesterId     any         //
	CategoryId      any         //
	TargetTorrentId any         // 续种目标种子
	ResultTorrentId any         // 求种完成后关联的种子
	Title           any         //
	Description     any         // 请求说明 (Markdown/BBCode)
	RewardAmount    any         // 已托管的魔力奖励
	Status          any         // 0=开放 1=已认领 2=待确认 3=已完成 4=已取消
	ClaimedBy       any         //
	ClaimedAt       *gtime.Time //
	ClaimExpiresAt  *gtime.Time //
	SubmittedAt     *gtime.Time //
	CompletedAt     *gtime.Time //
	CancelledBy     any         //
	CancelledAt     *gtime.Time //
	CancelReason    any         //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
