package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ModUserList(ctx context.Context, req *v1.ModUserListReq) (res *v1.ModUserListRes, err error) {
	out, err := service.AdminModUserUsecase().List(ctx, contexts.GetActor(ctx), req.ModUserListInp)
	if err != nil {
		return nil, err
	}
	return &v1.ModUserListRes{
		ListUserOut: *out,
	}, nil
}
