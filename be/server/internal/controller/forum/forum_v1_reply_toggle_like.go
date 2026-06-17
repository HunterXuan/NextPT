package forum

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/forum/v1"
)

func (c *ControllerV1) ReplyToggleLike(ctx context.Context, req *v1.ReplyToggleLikeReq) (res *v1.ReplyToggleLikeRes, err error) {
	err = service.ForumReplyUsecase().ToggleReplyLike(ctx, contexts.GetActor(ctx), req.ReplyToggleLikeInp)
	return
}
