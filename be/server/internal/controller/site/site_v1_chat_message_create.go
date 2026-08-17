package site

import (
	"context"

	v1 "server/api/site/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ChatMessageCreate(ctx context.Context, req *v1.ChatMessageCreateReq) (res *v1.ChatMessageCreateRes, err error) {
	item, err := service.SiteChatMessageUsecase().Create(ctx, contexts.GetActor(ctx), req.ChatMessageCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.ChatMessageCreateRes{ChatMessageItem: *item}, nil
}
