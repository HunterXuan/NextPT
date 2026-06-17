package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/service"
)

func (c *ControllerV1) InviteCheck(ctx context.Context, req *v1.InviteCheckReq) (res *v1.InviteCheckRes, err error) {
	out, err := service.IamInviteUsecase().Check(ctx, req.InviteCheckInp)
	if err != nil {
		return nil, err
	}
	res = &v1.InviteCheckRes{
		InviteCheckOut: out,
	}
	return
}
