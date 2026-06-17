package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserBan(ctx context.Context, req *v1.IamUserBanReq) (res *v1.IamUserBanRes, err error) {
	err = service.AdminIamUserUsecase().Ban(ctx, contexts.GetActor(ctx), req.IamUserBanInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamUserBanRes{}, nil
}
