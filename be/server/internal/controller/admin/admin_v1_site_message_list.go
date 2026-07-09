package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SiteMessageList(ctx context.Context, req *v1.SiteMessageListReq) (res *v1.SiteMessageListRes, err error) {
	out, err := service.SiteMessageUsecase().AdminList(ctx, contexts.GetActor(ctx), req.AdminMessageListInp)
	if err != nil {
		return nil, err
	}
	return &v1.SiteMessageListRes{MessageListOut: *out}, nil
}
