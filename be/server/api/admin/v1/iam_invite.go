package v1

import (
	"server/internal/model/in/adminin"

	"github.com/gogf/gf/v2/frame/g"
)

type IamInviteGrantReq struct {
	g.Meta `path:"/iam/invites:grant" method:"post" tags:"AdminIAM" summary:"发放邀请名额" perm:"admin:iam/invite:*"`
	adminin.IamInviteGrantInp
}

type IamInviteGrantRes struct{}
