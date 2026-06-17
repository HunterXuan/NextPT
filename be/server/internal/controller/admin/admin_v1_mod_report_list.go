package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ModReportList(ctx context.Context, req *v1.ModReportListReq) (res *v1.ModReportListRes, err error) {
	out, err := service.AdminModReportUsecase().List(ctx, contexts.GetActor(ctx), req.ModReportListInp)
	if err != nil {
		return nil, err
	}
	return &v1.ModReportListRes{
		ListReportsOut: *out,
	}, nil
}
