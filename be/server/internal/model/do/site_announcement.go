// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteAnnouncement is the golang structure of table site_announcement for DAO operations like Where/Data.
type SiteAnnouncement struct {
	g.Meta      `orm:"table:site_announcement, do:true"`
	Id          any         //
	Title       any         //
	Content     any         //
	Status      any         // 0=draft 1=published 2=archived
	CreatedBy   any         //
	UpdatedBy   any         //
	PublishedAt *gtime.Time //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
