package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumCategoryUpdate(ctx context.Context, req *v1.ForumCategoryUpdateReq) (res *v1.ForumCategoryUpdateRes, err error) {
	err = service.AdminForumCategoryUsecase().Update(ctx, contexts.GetActor(ctx), req.ForumCategoryUpdateInp)
	if err == nil {
		res = &v1.ForumCategoryUpdateRes{}
	}
	return
}
