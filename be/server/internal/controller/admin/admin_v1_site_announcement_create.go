package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SiteAnnouncementCreate(ctx context.Context, req *v1.SiteAnnouncementCreateReq) (res *v1.SiteAnnouncementCreateRes, err error) {
	out, err := service.SiteAnnouncementUsecase().AdminCreate(ctx, contexts.GetActor(ctx), req.AdminAnnouncementCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.SiteAnnouncementCreateRes{AnnouncementCreateOut: *out}, nil
}
