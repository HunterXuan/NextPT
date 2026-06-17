// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumTopicBookmark is the golang structure of table forum_topic_bookmark for DAO operations like Where/Data.
type ForumTopicBookmark struct {
	g.Meta    `orm:"table:forum_topic_bookmark, do:true"`
	Id        any         //
	UserId    any         //
	TopicId   any         //
	CreatedAt *gtime.Time //
}
