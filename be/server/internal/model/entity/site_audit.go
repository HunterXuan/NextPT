// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteAudit is the golang structure for table site_audit.
type SiteAudit struct {
	Id         uint64      `json:"id"         orm:"id"          description:""`
	UserId     uint64      `json:"userId"     orm:"user_id"     description:""`
	Action     string      `json:"action"     orm:"action"      description:""`
	TargetType string      `json:"targetType" orm:"target_type" description:""`
	TargetId   uint64      `json:"targetId"   orm:"target_id"   description:""`
	Detail     string      `json:"detail"     orm:"detail"      description:""`
	Ip         string      `json:"ip"         orm:"ip"          description:""`
	Level      int         `json:"level"      orm:"level"       description:"0=normal 1=important 2=critical"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`
}
