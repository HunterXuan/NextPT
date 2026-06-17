package forum

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/forum/v1"
)

func (c *ControllerV1) TopicReward(ctx context.Context, req *v1.TopicRewardReq) (res *v1.TopicRewardRes, err error) {
	err = service.ForumTopicUsecase().RewardTopic(ctx, contexts.GetActor(ctx), req.TopicRewardInp)
	return
}
