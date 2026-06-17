// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteMessage is the golang structure of table site_message for DAO operations like Where/Data.
type SiteMessage struct {
	g.Meta     `orm:"table:site_message, do:true"`
	Id         any         //
	SenderId   any         // 0=系统通知, 或管理员ID
	ReceiverId any         //
	Content    any         //
	IsRead     any         //
	CreatedAt  *gtime.Time //
}
