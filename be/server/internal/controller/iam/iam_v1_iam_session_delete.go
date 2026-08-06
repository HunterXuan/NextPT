package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamSessionDelete(ctx context.Context, req *v1.IamSessionDeleteReq) (res *v1.IamSessionDeleteRes, err error) {
	err = service.IamSessionUsecase().Delete(ctx, contexts.GetActor(ctx), contexts.GetSessionId(ctx))
	if err == nil {
		res = &v1.IamSessionDeleteRes{}
	}
	return
}
