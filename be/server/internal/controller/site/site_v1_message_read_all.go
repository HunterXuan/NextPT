package site

import (
	"context"

	v1 "server/api/site/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) MessageReadAll(ctx context.Context, req *v1.MessageReadAllReq) (res *v1.MessageReadAllRes, err error) {
	if err := service.SiteMessageUsecase().MarkAllRead(ctx, contexts.GetActor(ctx), req.MessageReadAllInp); err != nil {
		return nil, err
	}
	return &v1.MessageReadAllRes{}, nil
}
