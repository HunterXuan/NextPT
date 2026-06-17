package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SysCronLogList(ctx context.Context, req *v1.SysCronLogListReq) (res *v1.SysCronLogListRes, err error) {
	out, err := service.AdminSysCronUsecase().LogList(ctx, contexts.GetActor(ctx), req.Name, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &v1.SysCronLogListRes{SysCronLogListOut: *out}, nil
}
