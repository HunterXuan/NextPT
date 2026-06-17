// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumTopicLike is the golang structure of table forum_topic_like for DAO operations like Where/Data.
type ForumTopicLike struct {
	g.Meta    `orm:"table:forum_topic_like, do:true"`
	Id        any         //
	UserId    any         //
	TopicId   any         //
	CreatedAt *gtime.Time //
}
