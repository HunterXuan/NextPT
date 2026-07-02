package forum

import (
	"context"

	v1 "server/api/forum/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TopicRewardList(ctx context.Context, req *v1.TopicRewardListReq) (res *v1.TopicRewardListRes, err error) {
	out, err := service.ForumTopicUsecase().RewardList(ctx, contexts.GetActor(ctx), req.TopicRewardListInp)
	if err != nil {
		return nil, err
	}
	return &v1.TopicRewardListRes{
		TopicRewardListOut: *out,
	}, nil
}
