package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumNodeDelete(ctx context.Context, req *v1.ForumNodeDeleteReq) (res *v1.ForumNodeDeleteRes, err error) {
	err = service.AdminForumNodeUsecase().Delete(ctx, contexts.GetActor(ctx), req.ForumNodeDeleteInp)
	if err == nil {
		res = &v1.ForumNodeDeleteRes{}
	}
	return
}
