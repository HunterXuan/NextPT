package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumTopicDelete(ctx context.Context, req *v1.ForumTopicDeleteReq) (res *v1.ForumTopicDeleteRes, err error) {
	err = service.AdminForumTopicUsecase().Delete(ctx, contexts.GetActor(ctx), req.ForumTopicDeleteInp)
	if err == nil {
		res = &v1.ForumTopicDeleteRes{}
	}
	return
}
