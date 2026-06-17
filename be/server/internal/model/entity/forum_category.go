// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumCategory is the golang structure for table forum_category.
type ForumCategory struct {
	Id          uint        `json:"id"          orm:"id"            description:""`
	NameI18N    *gjson.Json `json:"nameI18N"    orm:"name_i18n"     description:"多语言名称映射"`
	DescI18N    *gjson.Json `json:"descI18N"    orm:"desc_i18n"     description:"多语言描述映射"`
	SortOrder   int         `json:"sortOrder"   orm:"sort_order"    description:""`
	MinRoleView int         `json:"minRoleView" orm:"min_role_view" description:"最低可见等级"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"    description:""`
}
