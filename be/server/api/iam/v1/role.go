package v1

import (
	"server/internal/model/out/iamout"

	"github.com/gogf/gf/v2/frame/g"
)

type RoleListReq struct {
	g.Meta `path:"/roles" method:"get" tags:"IamRole" summary:"获取角色等级标准"`
}

type RoleListRes struct {
	*iamout.RoleListOut
}
