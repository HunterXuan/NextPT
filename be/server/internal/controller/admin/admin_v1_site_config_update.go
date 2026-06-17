package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SiteConfigUpdate(ctx context.Context, req *v1.SiteConfigUpdateReq) (res *v1.SiteConfigUpdateRes, err error) {
	err = service.AdminSiteConfigUsecase().Update(ctx, contexts.GetActor(ctx), req.SiteConfigUpdateInp)
	if err != nil {
		return nil, err
	}
	return &v1.SiteConfigUpdateRes{}, nil
}
