// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumReply is the golang structure for table forum_reply.
type ForumReply struct {
	Id          uint64      `json:"id"          orm:"id"           description:""`
	TopicId     uint64      `json:"topicId"     orm:"topic_id"     description:""`
	UserId      uint64      `json:"userId"      orm:"user_id"      description:""`
	Content     string      `json:"content"     orm:"content"      description:"回复内容"`
	RewardCount uint        `json:"rewardCount" orm:"reward_count" description:"获得的魔力值打赏"`
	LikeCount   uint        `json:"likeCount"   orm:"like_count"   description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""`
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:""`
}
