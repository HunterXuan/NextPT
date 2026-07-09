// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteAnnouncement is the golang structure for table site_announcement.
type SiteAnnouncement struct {
	Id          uint64      `json:"id"          orm:"id"           description:""`
	Title       string      `json:"title"       orm:"title"        description:""`
	Content     string      `json:"content"     orm:"content"      description:""`
	Status      int         `json:"status"      orm:"status"       description:"0=draft 1=published 2=archived"`
	CreatedBy   uint64      `json:"createdBy"   orm:"created_by"   description:""`
	UpdatedBy   uint64      `json:"updatedBy"   orm:"updated_by"   description:""`
	PublishedAt *gtime.Time `json:"publishedAt" orm:"published_at" description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""`
}
