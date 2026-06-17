// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumTopic is the golang structure of table forum_topic for DAO operations like Where/Data.
type ForumTopic struct {
	g.Meta      `orm:"table:forum_topic, do:true"`
	Id          any         //
	NodeId      any         //
	UserId      any         //
	Subject     any         //
	Content     any         // 主题正文
	Appends     *gjson.Json // 追加内容 (V2EX风格) [{"content": "...", "created_at": "..."}]
	IsLocked    any         //
	IsSticky    any         //
	Views       any         //
	ReplyCount  any         //
	LikeCount   any         //
	LastReplyId any         //
	LastReplyAt *gtime.Time //
	LastReplyBy any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
