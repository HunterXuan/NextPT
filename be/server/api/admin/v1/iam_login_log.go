package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type IamLoginLogListReq struct {
	g.Meta `path:"/iam/login-logs" method:"get" tags:"AdminIAM" summary:"获取用户登录记录" perm:"admin:iam/user:*"`
	adminin.IamLoginLogListInp
}

type IamLoginLogListRes struct {
	adminout.IamLoginLogListOut
}
