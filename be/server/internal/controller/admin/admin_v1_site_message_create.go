package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SiteMessageCreate(ctx context.Context, req *v1.SiteMessageCreateReq) (res *v1.SiteMessageCreateRes, err error) {
	out, err := service.SiteMessageUsecase().AdminCreate(ctx, contexts.GetActor(ctx), req.AdminMessageCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.SiteMessageCreateRes{MessageCreateOut: *out}, nil
}
