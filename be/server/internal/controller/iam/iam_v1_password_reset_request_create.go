package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/service"
)

func (c *ControllerV1) PasswordResetRequestCreate(ctx context.Context, req *v1.PasswordResetRequestCreateReq) (res *v1.PasswordResetRequestCreateRes, err error) {
	if err := service.IamUserUsecase().CreatePasswordResetRequest(ctx, req.PasswordResetRequestCreateInp); err != nil {
		return nil, err
	}
	return &v1.PasswordResetRequestCreateRes{}, nil
}
