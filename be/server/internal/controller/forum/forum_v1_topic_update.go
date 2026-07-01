package forum

import (
	"context"

	"server/api/forum/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TopicUpdate(ctx context.Context, req *v1.TopicUpdateReq) (res *v1.TopicUpdateRes, err error) {
	err = service.ForumTopicUsecase().Update(ctx, contexts.GetActor(ctx), req.TopicUpdateInp)
	if err != nil {
		return nil, err
	}
	return &v1.TopicUpdateRes{}, nil
}
