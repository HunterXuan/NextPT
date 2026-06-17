// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamRole is the golang structure of table iam_role for DAO operations like Where/Data.
type IamRole struct {
	g.Meta      `orm:"table:iam_role, do:true"`
	Id          any         //
	Level       any         // 等级权重(用于权限比对，值越大权限越高，如普通用户10，管理员100)
	NameI18N    *gjson.Json // 角色名称多语言映射字典
	Rules       *gjson.Json // 角色规则(JSON: 包含 upgrade 升级条件, keep 保级条件等)
	Permissions *gjson.Json // 角色关联的权限标识符列表
	IsStaff     any         // 是否为管理组成员
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
