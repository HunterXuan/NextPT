package forum

import (
	"context"

	v1 "server/api/forum/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TopicAppend(ctx context.Context, req *v1.TopicAppendReq) (res *v1.TopicAppendRes, err error) {
	err = service.ForumTopicUsecase().Append(ctx, contexts.GetActor(ctx), req.TopicAppendInp)
	if err != nil {
		return nil, err
	}
	return &v1.TopicAppendRes{}, nil
}
