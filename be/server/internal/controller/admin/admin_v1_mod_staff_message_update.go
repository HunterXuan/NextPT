package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ModStaffMessageUpdate(ctx context.Context, req *v1.ModStaffMessageUpdateReq) (res *v1.ModStaffMessageUpdateRes, err error) {
	if err := service.ModStaffMessageUsecase().AdminProcess(ctx, contexts.GetActor(ctx), req.StaffMessageProcessInp); err != nil {
		return nil, err
	}
	return &v1.ModStaffMessageUpdateRes{}, nil
}
