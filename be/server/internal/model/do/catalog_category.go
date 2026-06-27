// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogCategory is the golang structure of table catalog_category for DAO operations like Where/Data.
type CatalogCategory struct {
	g.Meta    `orm:"table:catalog_category, do:true"`
	Id        any         //
	NameI18N  *gjson.Json // 多语言名称映射
	Slug      any         // URL-friendly
	SortOrder any         //
	Enabled   any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
