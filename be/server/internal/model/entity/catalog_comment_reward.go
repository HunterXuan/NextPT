// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogCommentReward is the golang structure for table catalog_comment_reward.
type CatalogCommentReward struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`
	CommentId uint64      `json:"commentId" orm:"comment_id" description:""`
	Amount    float64     `json:"amount"    orm:"amount"     description:"打赏金额"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
