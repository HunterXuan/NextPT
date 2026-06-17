package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamUserMe(ctx context.Context, req *v1.IamUserMeReq) (res *v1.IamUserMeRes, err error) {
	out, err := service.IamUserUsecase().Me(ctx, contexts.GetActor(ctx))
	if err != nil {
		return nil, err
	}
	res = &v1.IamUserMeRes{UserMeOut: *out}
	return
}
