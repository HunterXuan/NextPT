// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteChatMessage is the golang structure of table site_chat_message for DAO operations like Where/Data.
type SiteChatMessage struct {
	g.Meta    `orm:"table:site_chat_message, do:true"`
	Id        any         //
	UserId    any         //
	Content   any         //
	CreatedAt *gtime.Time //
}
