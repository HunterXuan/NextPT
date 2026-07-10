package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserLoginLogList(ctx context.Context, req *v1.IamUserLoginLogListReq) (res *v1.IamUserLoginLogListRes, err error) {
	out, err := service.IamUserUsecase().LoginLogs(ctx, contexts.GetActor(ctx), req.UserLoginLogListInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamUserLoginLogListRes{UserLoginLogListOut: *out}, nil
}
