// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ModReport is the golang structure for table mod_report.
type ModReport struct {
	Id           uint64      `json:"id"           orm:"id"            description:""`
	ReporterId   uint64      `json:"reporterId"   orm:"reporter_id"   description:""`
	TargetType   string      `json:"targetType"   orm:"target_type"   description:"torrent/user/offer/request/forum_post/comment/subtitle"`
	TargetId     uint64      `json:"targetId"     orm:"target_id"     description:""`
	Reason       string      `json:"reason"       orm:"reason"        description:""`
	Status       int         `json:"status"       orm:"status"        description:"0=pending 1=resolved 2=rejected"`
	DealtBy      uint64      `json:"dealtBy"      orm:"dealt_by"      description:""`
	DealtComment string      `json:"dealtComment" orm:"dealt_comment" description:""`
	DealtAt      *gtime.Time `json:"dealtAt"      orm:"dealt_at"      description:""`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
}
