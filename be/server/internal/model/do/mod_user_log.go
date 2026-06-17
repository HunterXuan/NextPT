// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ModUserLog is the golang structure of table mod_user_log for DAO operations like Where/Data.
type ModUserLog struct {
	g.Meta     `orm:"table:mod_user_log, do:true"`
	Id         any         //
	UserId     any         //
	ModType    any         // 1=warned 2=banned 3=leech_warned 4=upload_banned 5=download_banned 6=forum_banned
	Reason     any         //
	ExpireAt   *gtime.Time // 过期时间，NULL=永久
	ModBy      any         // 操作人ID
	ModComment any         //
	IsActive   any         //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
}
