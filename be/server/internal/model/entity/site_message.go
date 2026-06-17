// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteMessage is the golang structure for table site_message.
type SiteMessage struct {
	Id         uint64      `json:"id"         orm:"id"          description:""`
	SenderId   uint64      `json:"senderId"   orm:"sender_id"   description:"0=系统通知, 或管理员ID"`
	ReceiverId uint64      `json:"receiverId" orm:"receiver_id" description:""`
	Content    string      `json:"content"    orm:"content"     description:""`
	IsRead     bool        `json:"isRead"     orm:"is_read"     description:""`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`
}
