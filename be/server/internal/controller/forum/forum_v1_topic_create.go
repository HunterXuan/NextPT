package forum

import (
	"context"

	v1 "server/api/forum/v1"
	"server/internal/library/contexts"
	"server/internal/model/out/forumout"
	"server/internal/service"
)

func (c *ControllerV1) TopicCreate(ctx context.Context, req *v1.TopicCreateReq) (res *v1.TopicCreateRes, err error) {
	id, err := service.ForumTopicUsecase().Create(ctx, contexts.GetActor(ctx), req.TopicCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.TopicCreateRes{
		TopicCreateOut: forumout.TopicCreateOut{Id: id},
	}, nil
}
