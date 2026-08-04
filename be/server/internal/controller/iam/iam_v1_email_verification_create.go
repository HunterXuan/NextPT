package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/service"
)

func (c *ControllerV1) EmailVerificationCreate(ctx context.Context, req *v1.EmailVerificationCreateReq) (res *v1.EmailVerificationCreateRes, err error) {
	if err := service.IamUserUsecase().CreateEmailVerification(ctx, req.EmailVerificationCreateInp); err != nil {
		return nil, err
	}
	return &v1.EmailVerificationCreateRes{}, nil
}
