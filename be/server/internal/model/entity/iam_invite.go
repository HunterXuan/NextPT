// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamInvite is the golang structure for table iam_invite.
type IamInvite struct {
	Id           uint64      `json:"id"           orm:"id"            description:""`
	InviterId    uint64      `json:"inviterId"    orm:"inviter_id"    description:"邀请人"`
	InviteeEmail string      `json:"inviteeEmail" orm:"invitee_email" description:"被邀请人邮箱"`
	InviteeId    uint64      `json:"inviteeId"    orm:"invitee_id"    description:"被邀请人ID（注册后回填）"`
	Hash         string      `json:"hash"         orm:"hash"          description:"邀请码 (发放名额时就预生成唯一码)"`
	Status       int         `json:"status"       orm:"status"        description:"0=未分配/待发送 1=已发送 2=已注册 3=已过期 4=已回收"`
	IsTemporary  bool        `json:"isTemporary"  orm:"is_temporary"  description:"是否限时邀请"`
	ExpireAt     *gtime.Time `json:"expireAt"     orm:"expire_at"     description:""`
	UsedAt       *gtime.Time `json:"usedAt"       orm:"used_at"       description:""`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
}
