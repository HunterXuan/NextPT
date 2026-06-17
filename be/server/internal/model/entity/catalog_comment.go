// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogComment is the golang structure for table catalog_comment.
type CatalogComment struct {
	Id          uint64      `json:"id"          orm:"id"           description:""`
	TargetType  string      `json:"targetType"  orm:"target_type"  description:"torrent/offer/request"`
	TargetId    uint64      `json:"targetId"    orm:"target_id"    description:""`
	UserId      uint64      `json:"userId"      orm:"user_id"      description:""`
	Content     string      `json:"content"     orm:"content"      description:"评论内容 (Markdown/BBCode)"`
	RewardCount uint        `json:"rewardCount" orm:"reward_count" description:"获得的魔力值打赏"`
	LikeCount   uint        `json:"likeCount"   orm:"like_count"   description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""`
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:""`
}
