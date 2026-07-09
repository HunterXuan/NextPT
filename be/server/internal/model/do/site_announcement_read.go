// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteAnnouncementRead is the golang structure of table site_announcement_read for DAO operations like Where/Data.
type SiteAnnouncementRead struct {
	g.Meta         `orm:"table:site_announcement_read, do:true"`
	Id             any         //
	AnnouncementId any         //
	UserId         any         //
	ReadAt         *gtime.Time //
}
