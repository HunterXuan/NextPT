// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogRequest is the golang structure for table catalog_request.
type CatalogRequest struct {
	Id              uint64      `json:"id"              orm:"id"                description:""`
	RequestType     uint        `json:"requestType"     orm:"request_type"      description:"1=求种 2=续种"`
	RequesterId     uint64      `json:"requesterId"     orm:"requester_id"      description:""`
	CategoryId      uint        `json:"categoryId"      orm:"category_id"       description:""`
	TargetTorrentId uint64      `json:"targetTorrentId" orm:"target_torrent_id" description:"续种目标种子"`
	ResultTorrentId uint64      `json:"resultTorrentId" orm:"result_torrent_id" description:"求种完成后关联的种子"`
	Title           string      `json:"title"           orm:"title"             description:""`
	Description     string      `json:"description"     orm:"description"       description:"请求说明 (Markdown/BBCode)"`
	RewardAmount    float64     `json:"rewardAmount"    orm:"reward_amount"     description:"已托管的魔力奖励"`
	Status          uint        `json:"status"          orm:"status"            description:"0=开放 1=已认领 2=待确认 3=已完成 4=已取消"`
	ClaimedBy       uint64      `json:"claimedBy"       orm:"claimed_by"        description:""`
	ClaimedAt       *gtime.Time `json:"claimedAt"       orm:"claimed_at"        description:""`
	ClaimExpiresAt  *gtime.Time `json:"claimExpiresAt"  orm:"claim_expires_at"  description:""`
	SubmittedAt     *gtime.Time `json:"submittedAt"     orm:"submitted_at"      description:""`
	CompletedAt     *gtime.Time `json:"completedAt"     orm:"completed_at"      description:""`
	CancelledBy     uint64      `json:"cancelledBy"     orm:"cancelled_by"      description:""`
	CancelledAt     *gtime.Time `json:"cancelledAt"     orm:"cancelled_at"      description:""`
	CancelReason    string      `json:"cancelReason"    orm:"cancel_reason"     description:""`
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"        description:""`
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"        description:""`
}
