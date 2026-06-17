package forum

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/forum/v1"
)

func (c *ControllerV1) TopicUnbookmark(ctx context.Context, req *v1.TopicUnbookmarkReq) (res *v1.TopicUnbookmarkRes, err error) {
	err = service.ForumTopicUsecase().UnbookmarkTopic(ctx, contexts.GetActor(ctx), req.TopicUnbookmarkInp)
	return
}
