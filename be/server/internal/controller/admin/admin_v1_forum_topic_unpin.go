package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumTopicUnpin(ctx context.Context, req *v1.ForumTopicUnpinReq) (res *v1.ForumTopicUnpinRes, err error) {
	err = service.AdminForumTopicUsecase().Unpin(ctx, contexts.GetActor(ctx), req.ForumTopicUnpinInp)
	if err == nil {
		res = &v1.ForumTopicUnpinRes{}
	}
	return
}
