package site

import (
	"context"

	v1 "server/api/site/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) UserTaskRewardClaim(ctx context.Context, req *v1.UserTaskRewardClaimReq) (res *v1.UserTaskRewardClaimRes, err error) {
	if err := service.SiteTaskUsecase().ClaimReward(ctx, contexts.GetActor(ctx), req.UserTaskRewardClaimInp); err != nil {
		return nil, err
	}
	return &v1.UserTaskRewardClaimRes{}, nil
}
