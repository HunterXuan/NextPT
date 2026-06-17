package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type IamRoleListReq struct {
	g.Meta `path:"/iam/roles" method:"get" tags:"AdminIAM" summary:"获取角色列表" perm:"admin:iam/role:*"`
	adminin.IamRoleListInp
}

type IamRoleListRes struct {
	adminout.IamRoleListOut
}

type IamRoleCreateReq struct {
	g.Meta `path:"/iam/roles" method:"post" tags:"AdminIAM" summary:"创建角色" perm:"admin:iam/role:*"`
	adminin.IamRoleCreateInp
}

type IamRoleCreateRes struct {
	Id uint `json:"id"`
}

type IamRoleUpdateReq struct {
	g.Meta `path:"/iam/roles/{id}" method:"patch" tags:"AdminIAM" summary:"部分更新角色" perm:"admin:iam/role:*"`
	adminin.IamRoleUpdateInp
}

type IamRoleUpdateRes struct{}

type IamRoleDeleteReq struct {
	g.Meta `path:"/iam/roles/{id}" method:"delete" tags:"AdminIAM" summary:"删除角色" perm:"admin:iam/role:*"`
	adminin.IamRoleDeleteInp
}

type IamRoleDeleteRes struct{}
