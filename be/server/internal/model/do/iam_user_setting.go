// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserSetting is the golang structure of table iam_user_setting for DAO operations like Where/Data.
type IamUserSetting struct {
	g.Meta       `orm:"table:iam_user_setting, do:true"`
	Id           any         //
	UserId       any         //
	PrivacyLevel any         // 0=宽松 1=普通 2=严格
	Extra        *gjson.Json // UI偏好/分页/语言/时区等前端设置 (JSON 扩展)
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
