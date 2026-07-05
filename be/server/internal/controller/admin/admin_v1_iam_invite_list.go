package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamInviteList(ctx context.Context, req *v1.IamInviteListReq) (res *v1.IamInviteListRes, err error) {
	out, err := service.AdminIamInviteUsecase().List(ctx, contexts.GetActor(ctx), req.IamInviteListInp)
	if err == nil {
		res = &v1.IamInviteListRes{IamInviteListOut: out}
	}
	return
}
