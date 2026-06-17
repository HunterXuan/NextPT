// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumReplyLike is the golang structure of table forum_reply_like for DAO operations like Where/Data.
type ForumReplyLike struct {
	g.Meta    `orm:"table:forum_reply_like, do:true"`
	Id        any         //
	UserId    any         //
	ReplyId   any         //
	CreatedAt *gtime.Time //
}
