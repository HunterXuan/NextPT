package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/service"
)

func (c *ControllerV1) EmailVerificationRequestCreate(ctx context.Context, req *v1.EmailVerificationRequestCreateReq) (res *v1.EmailVerificationRequestCreateRes, err error) {
	if err := service.IamUserUsecase().CreateEmailVerificationRequest(ctx, req.EmailVerificationRequestCreateInp); err != nil {
		return nil, err
	}
	return &v1.EmailVerificationRequestCreateRes{}, nil
}
