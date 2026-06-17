package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumTopicUnlock(ctx context.Context, req *v1.ForumTopicUnlockReq) (res *v1.ForumTopicUnlockRes, err error) {
	err = service.AdminForumTopicUsecase().Unlock(ctx, contexts.GetActor(ctx), req.ForumTopicUnlockInp)
	if err == nil {
		res = &v1.ForumTopicUnlockRes{}
	}
	return
}
