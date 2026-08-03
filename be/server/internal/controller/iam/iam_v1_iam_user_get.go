package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserGet(ctx context.Context, req *v1.IamUserGetReq) (res *v1.IamUserGetRes, err error) {
	out, err := service.IamUserUsecase().Get(ctx, contexts.GetActor(ctx), req.UserGetInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamUserGetRes{UserGetOut: *out}, nil
}
