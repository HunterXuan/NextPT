// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamInvite is the golang structure of table iam_invite for DAO operations like Where/Data.
type IamInvite struct {
	g.Meta       `orm:"table:iam_invite, do:true"`
	Id           any         //
	InviterId    any         // 邀请人
	InviteeEmail any         // 被邀请人邮箱
	InviteeId    any         // 被邀请人ID（注册后回填）
	Hash         any         // 邀请码 (发放名额时就预生成唯一码)
	Status       any         // 0=未分配/待发送 1=已发送 2=已注册 3=已过期 4=已回收
	IsTemporary  any         // 是否限时邀请
	ExpireAt     *gtime.Time //
	UsedAt       *gtime.Time //
	CreatedAt    *gtime.Time //
}
