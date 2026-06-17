package v1

import (
	"server/internal/model/in/adminin"

	"github.com/gogf/gf/v2/frame/g"
)

type IamSessionDeleteReq struct {
	g.Meta `path:"/iam/users/{userId}/sessions" method:"delete" tags:"AdminIAM" summary:"删除用户所有会话(踢下线)" perm:"admin:iam/user:*"`
	adminin.IamSessionDeleteInp
}

type IamSessionDeleteRes struct{}
