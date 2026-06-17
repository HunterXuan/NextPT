// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumNode is the golang structure of table forum_node for DAO operations like Where/Data.
type ForumNode struct {
	g.Meta        `orm:"table:forum_node, do:true"`
	Id            any         //
	CategoryId    any         // 所属节点分类
	Slug          any         // 节点英文标识符
	NameI18N      *gjson.Json // 多语言名称映射
	DescI18N      *gjson.Json // 多语言描述映射
	SortOrder     any         //
	MinRoleRead   any         //
	MinRoleWrite  any         //
	MinRoleCreate any         //
	TopicCount    any         //
	ReplyCount    any         //
	LastTopicId   any         //
	Moderators    *gjson.Json // 板块版主ID列表(前端展示用)
	LastReplyAt   *gtime.Time //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
