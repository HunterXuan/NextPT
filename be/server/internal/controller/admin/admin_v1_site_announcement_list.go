package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SiteAnnouncementList(ctx context.Context, req *v1.SiteAnnouncementListReq) (res *v1.SiteAnnouncementListRes, err error) {
	out, err := service.SiteAnnouncementUsecase().AdminList(ctx, contexts.GetActor(ctx), req.AdminAnnouncementListInp)
	if err != nil {
		return nil, err
	}
	return &v1.SiteAnnouncementListRes{AnnouncementListOut: *out}, nil
}
