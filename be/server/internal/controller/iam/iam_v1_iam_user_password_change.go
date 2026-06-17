package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserPasswordChange(ctx context.Context, req *v1.IamUserPasswordChangeReq) (res *v1.IamUserPasswordChangeRes, err error) {
	if err := service.IamUserUsecase().ChangePassword(ctx, contexts.GetActor(ctx), req.UserPasswordChangeInp); err != nil {
		return nil, err
	}
	return &v1.IamUserPasswordChangeRes{}, nil
}
