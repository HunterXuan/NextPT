package v1

import (
	"server/internal/model/in/adminin"

	"github.com/gogf/gf/v2/frame/g"
)

type IamUserPermissionGrantReq struct {
	g.Meta `path:"/iam/users/{userId}/permissions:grant" method:"post" tags:"AdminIAM" summary:"赋予/收回用户权限规则" perm:"admin:iam/user:*"`
	adminin.IamUserPermissionGrantInp
}

type IamUserPermissionGrantRes struct{}

type IamUserPermissionRevokeReq struct {
	g.Meta `path:"/iam/users/{userId}/permissions:revoke" method:"post" tags:"AdminIAM" summary:"撤销用户特定的权限规则" perm:"admin:iam/user:*"`
	adminin.IamUserPermissionRevokeInp
}

type IamUserPermissionRevokeRes struct{}
