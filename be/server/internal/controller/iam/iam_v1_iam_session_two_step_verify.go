package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/service"
)

func (c *ControllerV1) IamSessionTwoStepVerify(ctx context.Context, req *v1.IamSessionTwoStepVerifyReq) (res *v1.IamSessionTwoStepVerifyRes, err error) {
	out, err := service.IamSessionUsecase().VerifyTwoStep(ctx, req.SessionTwoStepVerifyInp)
	if err != nil {
		return nil, err
	}
	return &v1.IamSessionTwoStepVerifyRes{SessionCreateOut: *out}, nil
}
