package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SiteConfigList(ctx context.Context, req *v1.SiteConfigListReq) (res *v1.SiteConfigListRes, err error) {
	out, err := service.AdminSiteConfigUsecase().List(ctx, contexts.GetActor(ctx), req.SiteConfigListInp)
	if err != nil {
		return nil, err
	}
	res = &v1.SiteConfigListRes{SiteConfigListOut: *out}
	return res, nil
}
