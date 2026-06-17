package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/service"
)

func (c *ControllerV1) IamSessionCreate(ctx context.Context, req *v1.IamSessionCreateReq) (res *v1.IamSessionCreateRes, err error) {
	out, err := service.IamSessionUsecase().Create(ctx, req.SessionCreateInp)
	if err != nil {
		return nil, err
	}
	res = &v1.IamSessionCreateRes{SessionCreateOut: *out}
	return
}
