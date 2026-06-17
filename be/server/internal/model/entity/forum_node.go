// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumNode is the golang structure for table forum_node.
type ForumNode struct {
	Id            uint        `json:"id"            orm:"id"              description:""`
	CategoryId    uint        `json:"categoryId"    orm:"category_id"     description:"所属节点分类"`
	Slug          string      `json:"slug"          orm:"slug"            description:"节点英文标识符"`
	NameI18N      *gjson.Json `json:"nameI18N"      orm:"name_i18n"       description:"多语言名称映射"`
	DescI18N      *gjson.Json `json:"descI18N"      orm:"desc_i18n"       description:"多语言描述映射"`
	SortOrder     int         `json:"sortOrder"     orm:"sort_order"      description:""`
	MinRoleRead   int         `json:"minRoleRead"   orm:"min_role_read"   description:""`
	MinRoleWrite  int         `json:"minRoleWrite"  orm:"min_role_write"  description:""`
	MinRoleCreate int         `json:"minRoleCreate" orm:"min_role_create" description:""`
	TopicCount    uint        `json:"topicCount"    orm:"topic_count"     description:""`
	ReplyCount    uint        `json:"replyCount"    orm:"reply_count"     description:""`
	LastTopicId   uint64      `json:"lastTopicId"   orm:"last_topic_id"   description:""`
	Moderators    *gjson.Json `json:"moderators"    orm:"moderators"      description:"板块版主ID列表(前端展示用)"`
	LastReplyAt   *gtime.Time `json:"lastReplyAt"   orm:"last_reply_at"   description:""`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:""`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      description:""`
}
