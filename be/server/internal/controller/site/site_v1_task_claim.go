package site

import (
	"context"

	v1 "server/api/site/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TaskClaim(ctx context.Context, req *v1.TaskClaimReq) (res *v1.TaskClaimRes, err error) {
	out, err := service.SiteTaskUsecase().Claim(ctx, contexts.GetActor(ctx), req.TaskClaimInp)
	if err != nil {
		return nil, err
	}
	return &v1.TaskClaimRes{TaskClaimOut: *out}, nil
}
