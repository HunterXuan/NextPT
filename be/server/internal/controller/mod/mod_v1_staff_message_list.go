package mod

import (
	"context"

	v1 "server/api/mod/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) StaffMessageList(ctx context.Context, req *v1.StaffMessageListReq) (res *v1.StaffMessageListRes, err error) {
	out, err := service.ModStaffMessageUsecase().List(ctx, contexts.GetActor(ctx), req.StaffMessageListInp)
	if err != nil {
		return nil, err
	}
	return &v1.StaffMessageListRes{StaffMessageListOut: *out}, nil
}
