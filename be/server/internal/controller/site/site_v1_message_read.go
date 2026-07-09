package site

import (
	"context"

	v1 "server/api/site/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) MessageRead(ctx context.Context, req *v1.MessageReadReq) (res *v1.MessageReadRes, err error) {
	if err := service.SiteMessageUsecase().MarkRead(ctx, contexts.GetActor(ctx), req.MessageReadInp); err != nil {
		return nil, err
	}
	return &v1.MessageReadRes{}, nil
}
