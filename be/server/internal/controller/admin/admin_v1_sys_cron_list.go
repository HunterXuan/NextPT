package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SysCronList(ctx context.Context, req *v1.SysCronListReq) (res *v1.SysCronListRes, err error) {
	out, err := service.AdminSysCronUsecase().List(ctx, contexts.GetActor(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.SysCronListRes{SysCronListOut: *out}, nil
}
