package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ModStaffMessageList(ctx context.Context, req *v1.ModStaffMessageListReq) (res *v1.ModStaffMessageListRes, err error) {
	out, err := service.ModStaffMessageUsecase().AdminList(ctx, contexts.GetActor(ctx), req.AdminStaffMessageListInp)
	if err != nil {
		return nil, err
	}
	return &v1.ModStaffMessageListRes{StaffMessageListOut: *out}, nil
}
