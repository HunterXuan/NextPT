package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type IamUserPermissionDetailReq struct {
	g.Meta `path:"/iam/users/{userId}/permissions" method:"get" tags:"AdminIAM" summary:"获取用户权限详情" perm:"admin:iam/user:*"`
	adminin.IamUserPermissionDetailInp
}

type IamUserPermissionDetailRes struct {
	adminout.IamUserPermissionDetailOut
}

type IamUserPermissionGrantReq struct {
	g.Meta `path:"/iam/users/{userId}/permissions:grant" method:"post" tags:"AdminIAM" summary:"批量添加用户权限规则" perm:"admin:iam/user:*"`
	adminin.IamUserPermissionGrantInp
}

type IamUserPermissionGrantRes struct{}

type IamUserPermissionRevokeReq struct {
	g.Meta `path:"/iam/users/{userId}/permissions:revoke" method:"post" tags:"AdminIAM" summary:"批量撤销用户权限规则" perm:"admin:iam/user:*"`
	adminin.IamUserPermissionRevokeInp
}

type IamUserPermissionRevokeRes struct{}
