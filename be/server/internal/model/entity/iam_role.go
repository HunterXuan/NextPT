// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamRole is the golang structure for table iam_role.
type IamRole struct {
	Id          uint        `json:"id"          orm:"id"          description:""`
	Level       int         `json:"level"       orm:"level"       description:"等级权重(用于权限比对，值越大权限越高，如普通用户10，管理员100)"`
	NameI18N    *gjson.Json `json:"nameI18N"    orm:"name_i18n"   description:"角色名称多语言映射字典"`
	Rules       *gjson.Json `json:"rules"       orm:"rules"       description:"角色规则(JSON: 包含 upgrade 升级条件, keep 保级条件等)"`
	Permissions *gjson.Json `json:"permissions" orm:"permissions" description:"角色关联的权限标识符列表"`
	IsStaff     bool        `json:"isStaff"     orm:"is_staff"    description:"是否为管理组成员"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  description:""`
}
