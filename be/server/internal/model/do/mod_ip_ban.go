// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ModIpBan is the golang structure of table mod_ip_ban for DAO operations like Where/Data.
type ModIpBan struct {
	g.Meta    `orm:"table:mod_ip_ban, do:true"`
	Id        any         //
	Network   any         // IP或CIDR，如 192.168.1.0/24
	Reason    any         //
	BannedBy  any         //
	CreatedAt *gtime.Time //
}
