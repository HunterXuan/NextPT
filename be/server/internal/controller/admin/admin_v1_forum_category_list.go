package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumCategoryList(ctx context.Context, req *v1.ForumCategoryListReq) (res *v1.ForumCategoryListRes, err error) {
	out, err := service.AdminForumCategoryUsecase().List(ctx, contexts.GetActor(ctx), req.ForumCategoryListInp)
	if err != nil {
		return nil, err
	}
	res = &v1.ForumCategoryListRes{ForumCategoryListOut: *out}
	return
}
