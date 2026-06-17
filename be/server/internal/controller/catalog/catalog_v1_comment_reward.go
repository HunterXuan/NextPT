package catalog

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/catalog/v1"
)

func (c *ControllerV1) CommentReward(ctx context.Context, req *v1.CommentRewardReq) (res *v1.CommentRewardRes, err error) {
	err = service.CatalogCommentUsecase().Reward(ctx, contexts.GetActor(ctx), req.CommentRewardInp)
	if err != nil {
		return nil, err
	}
	return &v1.CommentRewardRes{}, nil
}
