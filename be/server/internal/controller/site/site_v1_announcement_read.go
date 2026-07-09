package site

import (
	"context"

	v1 "server/api/site/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) AnnouncementRead(ctx context.Context, req *v1.AnnouncementReadReq) (res *v1.AnnouncementReadRes, err error) {
	if err := service.SiteAnnouncementUsecase().MarkRead(ctx, contexts.GetActor(ctx), req.AnnouncementReadInp); err != nil {
		return nil, err
	}
	return &v1.AnnouncementReadRes{}, nil
}
