package forum

import (
	"context"

	v1 "server/api/forum/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TopicDetail(ctx context.Context, req *v1.TopicDetailReq) (res *v1.TopicDetailRes, err error) {
	res = &v1.TopicDetailRes{}
	out, err := service.ForumTopicUsecase().Detail(ctx, contexts.GetActor(ctx), req.TopicDetailInp)
	if err != nil {
		return nil, err
	}
	res.TopicDetailOut = *out
	return
}
