// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteConfig is the golang structure of table site_config for DAO operations like Where/Data.
type SiteConfig struct {
	g.Meta      `orm:"table:site_config, do:true"`
	Id          any         //
	Group       any         //
	Key         any         //
	Value       *gjson.Json // 支持存 boolean/number/array
	Description any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
