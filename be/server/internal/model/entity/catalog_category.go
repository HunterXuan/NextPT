// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogCategory is the golang structure for table catalog_category.
type CatalogCategory struct {
	Id           uint        `json:"id"           orm:"id"            description:""`
	NameI18N     *gjson.Json `json:"nameI18N"     orm:"name_i18n"     description:"多语言名称映射"`
	Slug         string      `json:"slug"         orm:"slug"          description:"URL-friendly"`
	SortOrder    int         `json:"sortOrder"    orm:"sort_order"    description:""`
	Enabled      bool        `json:"enabled"      orm:"enabled"       description:""`
	UploadConfig *gjson.Json `json:"uploadConfig" orm:"upload_config" description:"分类发布表单配置"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:""`
}
