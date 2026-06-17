// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumReply is the golang structure of table forum_reply for DAO operations like Where/Data.
type ForumReply struct {
	g.Meta      `orm:"table:forum_reply, do:true"`
	Id          any         //
	TopicId     any         //
	UserId      any         //
	Content     any         // 回复内容
	RewardCount any         // 获得的魔力值打赏
	LikeCount   any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
