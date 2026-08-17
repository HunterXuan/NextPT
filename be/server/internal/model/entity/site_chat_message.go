// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteChatMessage is the golang structure for table site_chat_message.
type SiteChatMessage struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`
	Content   string      `json:"content"   orm:"content"    description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
