package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/service"
)

func (c *ControllerV1) PasswordResetCreate(ctx context.Context, req *v1.PasswordResetCreateReq) (res *v1.PasswordResetCreateRes, err error) {
	if err := service.IamUserUsecase().CreatePasswordReset(ctx, req.PasswordResetCreateInp); err != nil {
		return nil, err
	}
	return &v1.PasswordResetCreateRes{}, nil
}
