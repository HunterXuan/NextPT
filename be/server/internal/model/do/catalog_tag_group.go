// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTagGroup is the golang structure of table catalog_tag_group for DAO operations like Where/Data.
type CatalogTagGroup struct {
	g.Meta      `orm:"table:catalog_tag_group, do:true"`
	Id          any         //
	NameI18N    *gjson.Json // 多语言名称映射
	Slug        any         // 英文标识如 resolution
	CategoryIds *gjson.Json // 适用的分类 ID 数组 (例: [1,2])，为空则全站通用
	SortOrder   any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
