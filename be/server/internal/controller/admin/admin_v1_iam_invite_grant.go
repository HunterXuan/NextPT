package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamInviteGrant(ctx context.Context, req *v1.IamInviteGrantReq) (res *v1.IamInviteGrantRes, err error) {
	err = service.AdminIamInviteUsecase().Grant(ctx, contexts.GetActor(ctx), req.IamInviteGrantInp)
	if err == nil {
		res = &v1.IamInviteGrantRes{}
	}
	return
}
