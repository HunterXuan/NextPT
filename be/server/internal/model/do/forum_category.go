// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ForumCategory is the golang structure of table forum_category for DAO operations like Where/Data.
type ForumCategory struct {
	g.Meta      `orm:"table:forum_category, do:true"`
	Id          any         //
	NameI18N    *gjson.Json // 多语言名称映射
	DescI18N    *gjson.Json // 多语言描述映射
	SortOrder   any         //
	MinRoleView any         // 最低可见等级
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
