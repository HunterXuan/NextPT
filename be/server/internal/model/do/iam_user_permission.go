// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserPermission is the golang structure of table iam_user_permission for DAO operations like Where/Data.
type IamUserPermission struct {
	g.Meta     `orm:"table:iam_user_permission, do:true"`
	Id         any         //
	UserId     any         //
	PermKey    any         // 权限标识符，如 update:torrent:123 或 update:torrent:*
	SourceType any         // 来源类型: 1=manual 2=user_mod
	SourceId   any         // 来源记录ID，manual=0，user_mod=mod_user_log.id
	ExpireAt   *gtime.Time // 权限过期时间，NULL=永久
	IsActive   any         // 是否生效
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
}
