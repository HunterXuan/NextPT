package v1

import (
	"server/internal/model/in/iamin"

	"github.com/gogf/gf/v2/frame/g"
)

type EmailVerificationRequestCreateReq struct {
	g.Meta `path:"/email-verification-requests" method:"post" tags:"IamUser" summary:"申请邮箱验证" noAuth:"true"`
	iamin.EmailVerificationRequestCreateInp
}

type EmailVerificationRequestCreateRes struct{}

type EmailVerificationCreateReq struct {
	g.Meta `path:"/email-verifications" method:"post" tags:"IamUser" summary:"确认邮箱验证" noAuth:"true"`
	iamin.EmailVerificationCreateInp
}

type EmailVerificationCreateRes struct{}
