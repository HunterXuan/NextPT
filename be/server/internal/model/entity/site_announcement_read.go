// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteAnnouncementRead is the golang structure for table site_announcement_read.
type SiteAnnouncementRead struct {
	Id             uint64      `json:"id"             orm:"id"              description:""`
	AnnouncementId uint64      `json:"announcementId" orm:"announcement_id" description:""`
	UserId         uint64      `json:"userId"         orm:"user_id"         description:""`
	ReadAt         *gtime.Time `json:"readAt"         orm:"read_at"         description:""`
}
