package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) InviteList(ctx context.Context, req *v1.InviteListReq) (res *v1.InviteListRes, err error) {
	out, err := service.IamInviteUsecase().List(ctx, contexts.GetActor(ctx), req.InviteListInp)
	if err != nil {
		return nil, err
	}
	res = &v1.InviteListRes{
		InviteListOut: out,
	}
	return
}
