package site

import (
	"context"

	v1 "server/api/site/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) MessageList(ctx context.Context, req *v1.MessageListReq) (res *v1.MessageListRes, err error) {
	out, err := service.SiteMessageUsecase().List(ctx, contexts.GetActor(ctx), req.MessageListInp)
	if err != nil {
		return nil, err
	}
	return &v1.MessageListRes{MessageListOut: *out}, nil
}
