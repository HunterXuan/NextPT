package mod

import (
	"context"

	v1 "server/api/mod/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) StaffMessageCreate(ctx context.Context, req *v1.StaffMessageCreateReq) (res *v1.StaffMessageCreateRes, err error) {
	out, err := service.ModStaffMessageUsecase().Create(ctx, contexts.GetActor(ctx), req.StaffMessageCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.StaffMessageCreateRes{StaffMessageCreateOut: *out}, nil
}
