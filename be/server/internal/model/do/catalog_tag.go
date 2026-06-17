// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTag is the golang structure of table catalog_tag for DAO operations like Where/Data.
type CatalogTag struct {
	g.Meta    `orm:"table:catalog_tag, do:true"`
	Id        any         //
	GroupId   any         // 所属分组
	NameI18N  *gjson.Json // 多语言名称映射
	SortOrder any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
