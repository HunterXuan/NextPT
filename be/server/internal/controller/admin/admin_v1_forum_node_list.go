package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ForumNodeList(ctx context.Context, req *v1.ForumNodeListReq) (res *v1.ForumNodeListRes, err error) {
	out, err := service.AdminForumNodeUsecase().List(ctx, contexts.GetActor(ctx), req.ForumNodeListInp)
	if err != nil {
		return nil, err
	}
	res = &v1.ForumNodeListRes{ForumNodeListOut: *out}
	return
}
