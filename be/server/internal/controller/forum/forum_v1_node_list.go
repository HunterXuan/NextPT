package forum

import (
	"context"

	v1 "server/api/forum/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) NodeList(ctx context.Context, req *v1.NodeListReq) (res *v1.NodeListRes, err error) {
	res = &v1.NodeListRes{}
	out, err := service.ForumNodeUsecase().List(ctx, contexts.GetActor(ctx))
	if err != nil {
		return nil, err
	}
	res.NodeListOut = *out
	return
}
