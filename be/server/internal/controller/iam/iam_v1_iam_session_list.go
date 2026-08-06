package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamSessionList(ctx context.Context, req *v1.IamSessionListReq) (res *v1.IamSessionListRes, err error) {
	out, err := service.IamSessionUsecase().List(ctx, contexts.GetActor(ctx), contexts.GetSessionId(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.IamSessionListRes{SessionListOut: *out}, nil
}
