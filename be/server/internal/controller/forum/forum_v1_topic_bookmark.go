package forum

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/forum/v1"
)

func (c *ControllerV1) TopicBookmark(ctx context.Context, req *v1.TopicBookmarkReq) (res *v1.TopicBookmarkRes, err error) {
	err = service.ForumTopicUsecase().BookmarkTopic(ctx, contexts.GetActor(ctx), req.TopicBookmarkInp)
	return
}
