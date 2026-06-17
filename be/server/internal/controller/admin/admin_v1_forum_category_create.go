package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumCategoryCreate(ctx context.Context, req *v1.ForumCategoryCreateReq) (res *v1.ForumCategoryCreateRes, err error) {
	err = service.AdminForumCategoryUsecase().Create(ctx, contexts.GetActor(ctx), req.ForumCategoryCreateInp)
	if err == nil {
		res = &v1.ForumCategoryCreateRes{}
	}
	return
}
