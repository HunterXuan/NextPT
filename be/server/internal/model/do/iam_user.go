// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUser is the golang structure of table iam_user for DAO operations like Where/Data.
type IamUser struct {
	g.Meta        `orm:"table:iam_user, do:true"`
	Id            any         //
	Username      any         //
	Email         any         //
	PasswordHash  any         // bcrypt hash
	Passkey       any         // Tracker passkey
	Status        any         // 0=pending 1=confirmed 2=disabled
	Role          any         // 当前角色ID (关联 user_role 表)
	VipUntil      *gtime.Time // VIP 过期时间
	VipRemark     any         // VIP 身份获取备注/来源
	TwoStepType   any         // 两步验证方式: 0=关闭 1=TOTP(Authenticator) 2=邮件验证码
	TwoStepSecret any         // TOTP 密钥 (two_step_type=1 时使用；邮件验证码走 Redis 临时存储)
	InvitedBy     any         //
	LastLogin     *gtime.Time // 最后登录时间
	LastIp        any         // 最后登录 IP
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	DeletedAt     *gtime.Time //
}
