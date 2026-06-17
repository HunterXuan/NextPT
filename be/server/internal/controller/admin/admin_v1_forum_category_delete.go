package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumCategoryDelete(ctx context.Context, req *v1.ForumCategoryDeleteReq) (res *v1.ForumCategoryDeleteRes, err error) {
	err = service.AdminForumCategoryUsecase().Delete(ctx, contexts.GetActor(ctx), req.ForumCategoryDeleteInp)
	if err == nil {
		res = &v1.ForumCategoryDeleteRes{}
	}
	return
}
