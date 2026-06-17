// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserProfile is the golang structure for table iam_user_profile.
type IamUserProfile struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`
	Avatar    string      `json:"avatar"    orm:"avatar"     description:""`
	Info      string      `json:"info"      orm:"info"       description:"个人简介 (Markdown)"`
	Signature string      `json:"signature" orm:"signature"  description:"论坛签名"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
