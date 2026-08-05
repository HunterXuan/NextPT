package v1

import (
	"server/internal/model/in/iamin"
	"server/internal/model/out/iamout"

	"github.com/gogf/gf/v2/frame/g"
)

type IamSessionCreateReq struct {
	g.Meta `path:"/sessions" method:"post" tags:"IamSession" summary:"创建会话 (登录)" noAuth:"true"`
	iamin.SessionCreateInp
}

type IamSessionCreateRes struct {
	iamout.SessionCreateOut
}

type IamSessionTwoStepVerifyReq struct {
	g.Meta `path:"/sessions:verifyTwoStep" method:"post" tags:"IamSession" summary:"完成两步验证登录" noAuth:"true"`
	iamin.SessionTwoStepVerifyInp
}

type IamSessionTwoStepVerifyRes struct {
	iamout.SessionCreateOut
}

type IamSessionDeleteReq struct {
	g.Meta `path:"/sessions" method:"delete" tags:"IamSession" summary:"销毁会话 (登出)"`
}

type IamSessionDeleteRes struct{}
