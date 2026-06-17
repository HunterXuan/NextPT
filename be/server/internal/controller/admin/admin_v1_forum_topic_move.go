package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumTopicMove(ctx context.Context, req *v1.ForumTopicMoveReq) (res *v1.ForumTopicMoveRes, err error) {
	err = service.AdminForumTopicUsecase().Move(ctx, contexts.GetActor(ctx), req.ForumTopicMoveInp)
	if err == nil {
		res = &v1.ForumTopicMoveRes{}
	}
	return
}
