package iam

import (
	"context"

	v1 "server/api/iam/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) IamSessionDeleteById(ctx context.Context, req *v1.IamSessionDeleteByIdReq) (res *v1.IamSessionDeleteByIdRes, err error) {
	if err = service.IamSessionUsecase().DeleteById(ctx, contexts.GetActor(ctx), req.Id); err != nil {
		return nil, err
	}
	return &v1.IamSessionDeleteByIdRes{}, nil
}
