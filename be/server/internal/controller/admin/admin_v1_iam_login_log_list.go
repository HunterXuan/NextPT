package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamLoginLogList(ctx context.Context, req *v1.IamLoginLogListReq) (res *v1.IamLoginLogListRes, err error) {
	out, err := service.AdminIamUserUsecase().LoginLogs(ctx, contexts.GetActor(ctx), req.IamLoginLogListInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamLoginLogListRes{IamLoginLogListOut: *out}, nil
}
