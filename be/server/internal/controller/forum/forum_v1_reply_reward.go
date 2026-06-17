package forum

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/forum/v1"
)

func (c *ControllerV1) ReplyReward(ctx context.Context, req *v1.ReplyRewardReq) (res *v1.ReplyRewardRes, err error) {
	err = service.ForumReplyUsecase().RewardReply(ctx, contexts.GetActor(ctx), req.ReplyRewardInp)
	return
}
