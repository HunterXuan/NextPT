package forum

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/forum/v1"
)

func (c *ControllerV1) TopicToggleLike(ctx context.Context, req *v1.TopicToggleLikeReq) (res *v1.TopicToggleLikeRes, err error) {
	err = service.ForumTopicUsecase().ToggleTopicLike(ctx, contexts.GetActor(ctx), req.TopicToggleLikeInp)
	return
}
