// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EconomyRewardRecord is the golang structure for table economy_reward_record.
type EconomyRewardRecord struct {
	Id         uint64      `json:"id"         orm:"id"           description:""`
	TargetType string      `json:"targetType" orm:"target_type"  description:"catalog_torrent/catalog_comment/forum_topic/forum_reply"`
	TargetId   uint64      `json:"targetId"   orm:"target_id"    description:""`
	FromUserId uint64      `json:"fromUserId" orm:"from_user_id" description:"赞赏者ID"`
	ToUserId   uint64      `json:"toUserId"   orm:"to_user_id"   description:"接收者ID"`
	Amount     float64     `json:"amount"     orm:"amount"       description:"赞赏金额"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"   description:""`
}
