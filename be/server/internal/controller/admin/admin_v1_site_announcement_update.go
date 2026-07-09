package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SiteAnnouncementUpdate(ctx context.Context, req *v1.SiteAnnouncementUpdateReq) (res *v1.SiteAnnouncementUpdateRes, err error) {
	if err := service.SiteAnnouncementUsecase().AdminUpdate(ctx, contexts.GetActor(ctx), req.AdminAnnouncementUpdateInp); err != nil {
		return nil, err
	}
	return &v1.SiteAnnouncementUpdateRes{}, nil
}
