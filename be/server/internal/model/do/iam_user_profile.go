// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserProfile is the golang structure of table iam_user_profile for DAO operations like Where/Data.
type IamUserProfile struct {
	g.Meta    `orm:"table:iam_user_profile, do:true"`
	Id        any         //
	UserId    any         //
	Avatar    any         //
	Info      any         // 个人简介 (Markdown)
	Signature any         // 论坛签名
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
