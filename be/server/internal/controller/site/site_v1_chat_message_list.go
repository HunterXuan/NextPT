package site

import (
	"context"

	v1 "server/api/site/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ChatMessageList(ctx context.Context, req *v1.ChatMessageListReq) (res *v1.ChatMessageListRes, err error) {
	out, err := service.SiteChatMessageUsecase().List(ctx, contexts.GetActor(ctx), req.ChatMessageListInp)
	if err != nil {
		return nil, err
	}
	return &v1.ChatMessageListRes{ChatMessageListOut: *out}, nil
}
