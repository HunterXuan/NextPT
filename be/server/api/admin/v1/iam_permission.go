package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type IamPermissionListReq struct {
	g.Meta `path:"/iam/permissions" method:"get" tags:"AdminIAM" summary:"获取系统权限列表" perm:"admin:iam/role:*"`
	adminin.IamPermissionListInp
}

type IamPermissionListRes struct {
	adminout.IamPermissionListOut
}
