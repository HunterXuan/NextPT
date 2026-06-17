package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ModReportResolve(ctx context.Context, req *v1.ModReportResolveReq) (res *v1.ModReportResolveRes, err error) {
	err = service.AdminModReportUsecase().Resolve(ctx, contexts.GetActor(ctx), req.ModReportResolveInp)
	if err != nil {
		return nil, err
	}
	return &v1.ModReportResolveRes{}, nil
}
