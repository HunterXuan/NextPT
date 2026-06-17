// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTagGroup is the golang structure for table catalog_tag_group.
type CatalogTagGroup struct {
	Id          uint        `json:"id"          orm:"id"           description:""`
	NameI18N    *gjson.Json `json:"nameI18N"    orm:"name_i18n"    description:"多语言名称映射"`
	Slug        string      `json:"slug"        orm:"slug"         description:"英文标识如 resolution"`
	CategoryIds *gjson.Json `json:"categoryIds" orm:"category_ids" description:"适用的分类 ID 数组 (例: [1,2])，为空则全站通用"`
	SortOrder   int         `json:"sortOrder"   orm:"sort_order"   description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""`
}
