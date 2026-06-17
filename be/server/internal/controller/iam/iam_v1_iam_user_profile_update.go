package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserProfileUpdate(ctx context.Context, req *v1.IamUserProfileUpdateReq) (res *v1.IamUserProfileUpdateRes, err error) {
	if err := service.IamUserUsecase().UpdateProfile(ctx, contexts.GetActor(ctx), req.UserProfileUpdateInp); err != nil {
		return nil, err
	}
	return &v1.IamUserProfileUpdateRes{}, nil
}
