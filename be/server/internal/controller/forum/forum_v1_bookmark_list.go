package forum

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/forum/v1"
)

func (c *ControllerV1) BookmarkList(ctx context.Context, req *v1.BookmarkListReq) (res *v1.BookmarkListRes, err error) {
	out, err := service.ForumTopicUsecase().ListBookmarkedTopics(ctx, contexts.GetActor(ctx), req.TopicBookmarkListInp)
	if err == nil {
		res = &v1.BookmarkListRes{TopicBookmarkListOut: *out}
	}
	return
}
