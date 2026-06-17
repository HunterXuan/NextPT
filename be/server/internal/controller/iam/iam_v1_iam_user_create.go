package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/service"
)

func (c *ControllerV1) IamUserCreate(ctx context.Context, req *v1.IamUserCreateReq) (res *v1.IamUserCreateRes, err error) {
	id, err := service.IamUserUsecase().Create(ctx, req.UserCreateInp)
	if err == nil {
		res = &v1.IamUserCreateRes{Id: id}
	}
	return
}
