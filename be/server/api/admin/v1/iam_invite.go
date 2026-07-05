package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type IamInviteListReq struct {
	g.Meta `path:"/iam/invites" method:"get" tags:"AdminIAM" summary:"获取站点级邀请列表" perm:"admin:iam/invite:*"`
	adminin.IamInviteListInp
}

type IamInviteListRes struct {
	*adminout.IamInviteListOut
}

type IamInviteGrantReq struct {
	g.Meta `path:"/iam/invites:grant" method:"post" tags:"AdminIAM" summary:"发放邀请名额" perm:"admin:iam/invite:*"`
	adminin.IamInviteGrantInp
}

type IamInviteGrantRes struct{}

type IamInviteRecycleReq struct {
	g.Meta `path:"/iam/invites/{id}:recycle" method:"post" tags:"AdminIAM" summary:"回收站点级邀请" perm:"admin:iam/invite:*"`
	adminin.IamInviteRecycleInp
}

type IamInviteRecycleRes struct{}
