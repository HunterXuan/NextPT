package forum

import (
	"context"

	v1 "server/api/forum/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TopicGetHot(ctx context.Context, req *v1.TopicGetHotReq) (res *v1.TopicGetHotRes, err error) {
	out, err := service.ForumTopicUsecase().ListHot(ctx, contexts.GetActor(ctx), req.Size)
	if err != nil {
		return nil, err
	}
	return &v1.TopicGetHotRes{TopicHotListOut: *out}, nil
}
