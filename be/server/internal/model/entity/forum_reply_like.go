// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumReplyLike is the golang structure for table forum_reply_like.
type ForumReplyLike struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`
	ReplyId   uint64      `json:"replyId"   orm:"reply_id"   description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
