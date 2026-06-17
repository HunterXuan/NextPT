package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) InviteSend(ctx context.Context, req *v1.InviteSendReq) (res *v1.InviteSendRes, err error) {
	err = service.IamInviteUsecase().Send(ctx, contexts.GetActor(ctx), req.InviteSendInp)
	if err == nil {
		res = &v1.InviteSendRes{}
	}
	return
}
