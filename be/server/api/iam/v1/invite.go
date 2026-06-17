package v1

import (
	"server/internal/model/in/iamin"
	"server/internal/model/out/iamout"

	"github.com/gogf/gf/v2/frame/g"
)

type InviteListReq struct {
	g.Meta `path:"/invites" method:"get" tags:"IamInvite" summary:"获取我的邀请码列表"`
	iamin.InviteListInp
}

type InviteListRes struct {
	*iamout.InviteListOut
}

type InviteSendReq struct {
	g.Meta `path:"/invites:send" method:"post" tags:"IamInvite" summary:"发送邀请码"`
	iamin.InviteSendInp
}

type InviteSendRes struct{}

type InviteCheckReq struct {
	g.Meta `path:"/invites:check" method:"get" tags:"IamInvite" summary:"校验邀请码有效性" noAuth:"true"`
	iamin.InviteCheckInp
}

type InviteCheckRes struct {
	*iamout.InviteCheckOut
}
