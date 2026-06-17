// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ModIpBan is the golang structure for table mod_ip_ban.
type ModIpBan struct {
	Id        uint        `json:"id"        orm:"id"         description:""`
	Network   string      `json:"network"   orm:"network"    description:"IP或CIDR，如 192.168.1.0/24"`
	Reason    string      `json:"reason"    orm:"reason"     description:""`
	BannedBy  uint64      `json:"bannedBy"  orm:"banned_by"  description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
