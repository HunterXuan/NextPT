// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTag is the golang structure for table catalog_tag.
type CatalogTag struct {
	Id        uint        `json:"id"        orm:"id"         description:""`
	GroupId   uint        `json:"groupId"   orm:"group_id"   description:"所属分组"`
	NameI18N  *gjson.Json `json:"nameI18N"  orm:"name_i18n"  description:"多语言名称映射"`
	Value     string      `json:"value"     orm:"value"      description:"稳定值"`
	SortOrder int         `json:"sortOrder" orm:"sort_order" description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
