// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserSetting is the golang structure for table iam_user_setting.
type IamUserSetting struct {
	Id           uint64      `json:"id"           orm:"id"            description:""`
	UserId       uint64      `json:"userId"       orm:"user_id"       description:""`
	PrivacyLevel int         `json:"privacyLevel" orm:"privacy_level" description:"0=宽松 1=普通 2=严格"`
	Extra        *gjson.Json `json:"extra"        orm:"extra"         description:"UI偏好/分页/语言/时区等前端设置 (JSON 扩展)"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:""`
}
