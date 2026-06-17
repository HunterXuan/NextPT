package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SiteAuditList(ctx context.Context, req *v1.SiteAuditListReq) (res *v1.SiteAuditListRes, err error) {
	out, err := service.AdminSiteAuditUsecase().List(ctx, contexts.GetActor(ctx), req.SiteAuditListInp)
	if err != nil {
		return nil, err
	}
	return &v1.SiteAuditListRes{
		SiteAuditListOut: *out,
	}, nil
}
