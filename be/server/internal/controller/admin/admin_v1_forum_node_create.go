package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumNodeCreate(ctx context.Context, req *v1.ForumNodeCreateReq) (res *v1.ForumNodeCreateRes, err error) {
	err = service.AdminForumNodeUsecase().Create(ctx, contexts.GetActor(ctx), req.ForumNodeCreateInp)
	if err == nil {
		res = &v1.ForumNodeCreateRes{}
	}
	return
}
