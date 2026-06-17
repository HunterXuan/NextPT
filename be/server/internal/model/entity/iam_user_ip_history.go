// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserIpHistory is the golang structure for table iam_user_ip_history.
type IamUserIpHistory struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`
	Ip        string      `json:"ip"        orm:"ip"         description:""`
	Type      string      `json:"type"      orm:"type"       description:"login/tracker/register"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
