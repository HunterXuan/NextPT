// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUser is the golang structure for table iam_user.
type IamUser struct {
	Id            uint64      `json:"id"            orm:"id"              description:""`
	Username      string      `json:"username"      orm:"username"        description:""`
	Email         string      `json:"email"         orm:"email"           description:""`
	PasswordHash  string      `json:"passwordHash"  orm:"password_hash"   description:"bcrypt hash"`
	Passkey       string      `json:"passkey"       orm:"passkey"         description:"Tracker passkey"`
	Status        int         `json:"status"        orm:"status"          description:"0=pending 1=confirmed 2=disabled"`
	Role          uint        `json:"role"          orm:"role"            description:"当前角色ID (关联 user_role 表)"`
	VipUntil      *gtime.Time `json:"vipUntil"      orm:"vip_until"       description:"VIP 过期时间"`
	VipRemark     string      `json:"vipRemark"     orm:"vip_remark"      description:"VIP 身份获取备注/来源"`
	TwoStepType   int         `json:"twoStepType"   orm:"two_step_type"   description:"两步验证方式: 0=关闭 1=TOTP(Authenticator) 2=邮件验证码"`
	TwoStepSecret string      `json:"twoStepSecret" orm:"two_step_secret" description:"TOTP 密钥 (two_step_type=1 时使用；邮件验证码走 Redis 临时存储)"`
	InvitedBy     uint64      `json:"invitedBy"     orm:"invited_by"      description:""`
	LastLogin     *gtime.Time `json:"lastLogin"     orm:"last_login"      description:"最后登录时间"`
	LastIp        string      `json:"lastIp"        orm:"last_ip"         description:"最后登录 IP"`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:""`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      description:""`
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"      description:""`
}
