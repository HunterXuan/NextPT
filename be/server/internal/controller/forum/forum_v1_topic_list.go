package forum

import (
	"context"

	v1 "server/api/forum/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TopicList(ctx context.Context, req *v1.TopicListReq) (res *v1.TopicListRes, err error) {
	res = &v1.TopicListRes{}
	out, err := service.ForumTopicUsecase().List(ctx, contexts.GetActor(ctx), req.TopicListInp)
	if err != nil {
		return nil, err
	}
	res.TopicListOut = *out
	return
}
