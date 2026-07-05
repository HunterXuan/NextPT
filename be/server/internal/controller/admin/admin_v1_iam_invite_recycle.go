package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamInviteRecycle(ctx context.Context, req *v1.IamInviteRecycleReq) (res *v1.IamInviteRecycleRes, err error) {
	err = service.AdminIamInviteUsecase().Recycle(ctx, contexts.GetActor(ctx), req.IamInviteRecycleInp)
	if err == nil {
		res = &v1.IamInviteRecycleRes{}
	}
	return
}
