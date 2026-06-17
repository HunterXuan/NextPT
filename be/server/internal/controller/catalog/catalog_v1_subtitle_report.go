package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SubtitleReport(ctx context.Context, req *v1.SubtitleReportReq) (res *v1.SubtitleReportRes, err error) {
	actor := contexts.GetActor(ctx)
	err = service.CatalogSubtitleUsecase().Report(ctx, actor, req.SubtitleReportInp)
	return nil, err
}
