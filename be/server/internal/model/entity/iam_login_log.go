// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamLoginLog is the golang structure for table iam_login_log.
type IamLoginLog struct {
	Id         uint64      `json:"id"         orm:"id"          description:""`
	UserId     uint64      `json:"userId"     orm:"user_id"     description:""`
	Ip         string      `json:"ip"         orm:"ip"          description:""`
	UserAgent  string      `json:"userAgent"  orm:"user_agent"  description:""`
	Result     int         `json:"result"     orm:"result"      description:"1=success 0=fail"`
	FailReason string      `json:"failReason" orm:"fail_reason" description:""`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`
}
