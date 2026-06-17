// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumTopic is the golang structure for table forum_topic.
type ForumTopic struct {
	Id          uint64      `json:"id"          orm:"id"            description:""`
	NodeId      uint        `json:"nodeId"      orm:"node_id"       description:""`
	UserId      uint64      `json:"userId"      orm:"user_id"       description:""`
	Subject     string      `json:"subject"     orm:"subject"       description:""`
	Content     string      `json:"content"     orm:"content"       description:"主题正文"`
	Appends     *gjson.Json `json:"appends"     orm:"appends"       description:"追加内容 (V2EX风格) [{\"content\": \"...\", \"created_at\": \"...\"}]"`
	IsLocked    bool        `json:"isLocked"    orm:"is_locked"     description:""`
	IsSticky    bool        `json:"isSticky"    orm:"is_sticky"     description:""`
	Views       uint        `json:"views"       orm:"views"         description:""`
	ReplyCount  uint        `json:"replyCount"  orm:"reply_count"   description:""`
	LikeCount   uint        `json:"likeCount"   orm:"like_count"    description:""`
	LastReplyId uint64      `json:"lastReplyId" orm:"last_reply_id" description:""`
	LastReplyAt *gtime.Time `json:"lastReplyAt" orm:"last_reply_at" description:""`
	LastReplyBy uint64      `json:"lastReplyBy" orm:"last_reply_by" description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"    description:""`
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"    description:""`
}
