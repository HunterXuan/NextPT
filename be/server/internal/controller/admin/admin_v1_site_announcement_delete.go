package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SiteAnnouncementDelete(ctx context.Context, req *v1.SiteAnnouncementDeleteReq) (res *v1.SiteAnnouncementDeleteRes, err error) {
	if err := service.SiteAnnouncementUsecase().AdminDelete(ctx, contexts.GetActor(ctx), req.AdminAnnouncementDeleteInp); err != nil {
		return nil, err
	}
	return &v1.SiteAnnouncementDeleteRes{}, nil
}
