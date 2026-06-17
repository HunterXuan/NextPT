package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumNodeUpdate(ctx context.Context, req *v1.ForumNodeUpdateReq) (res *v1.ForumNodeUpdateRes, err error) {
	err = service.AdminForumNodeUsecase().Update(ctx, contexts.GetActor(ctx), req.ForumNodeUpdateInp)
	if err == nil {
		res = &v1.ForumNodeUpdateRes{}
	}
	return
}
