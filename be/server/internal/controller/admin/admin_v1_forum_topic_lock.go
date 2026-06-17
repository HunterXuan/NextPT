package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumTopicLock(ctx context.Context, req *v1.ForumTopicLockReq) (res *v1.ForumTopicLockRes, err error) {
	err = service.AdminForumTopicUsecase().Lock(ctx, contexts.GetActor(ctx), req.ForumTopicLockInp)
	if err == nil {
		res = &v1.ForumTopicLockRes{}
	}
	return
}
