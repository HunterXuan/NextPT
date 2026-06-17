// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumTopicLike is the golang structure for table forum_topic_like.
type ForumTopicLike struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`
	TopicId   uint64      `json:"topicId"   orm:"topic_id"   description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
