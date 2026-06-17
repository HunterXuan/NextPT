package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumTopicPin(ctx context.Context, req *v1.ForumTopicPinReq) (res *v1.ForumTopicPinRes, err error) {
	err = service.AdminForumTopicUsecase().Pin(ctx, contexts.GetActor(ctx), req.ForumTopicPinInp)
	if err == nil {
		res = &v1.ForumTopicPinRes{}
	}
	return
}
