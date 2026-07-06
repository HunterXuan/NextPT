// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteAudit is the golang structure of table site_audit for DAO operations like Where/Data.
type SiteAudit struct {
	g.Meta     `orm:"table:site_audit, do:true"`
	Id         any         //
	UserId     any         //
	Action     any         //
	TargetType any         //
	TargetId   any         //
	Detail     any         //
	Ip         any         //
	Level      any         // 0=normal 1=important 2=critical
	CreatedAt  *gtime.Time //
}
