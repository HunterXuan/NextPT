package v1

import (
	"server/internal/model/in/iamin"

	"github.com/gogf/gf/v2/frame/g"
)

type PasswordResetRequestCreateReq struct {
	g.Meta `path:"/password-reset-requests" method:"post" tags:"IamUser" summary:"申请重置密码" noAuth:"true"`
	iamin.PasswordResetRequestCreateInp
}

type PasswordResetRequestCreateRes struct{}

type PasswordResetCreateReq struct {
	g.Meta `path:"/password-resets" method:"post" tags:"IamUser" summary:"确认重置密码" noAuth:"true"`
	iamin.PasswordResetCreateInp
}

type PasswordResetCreateRes struct{}
