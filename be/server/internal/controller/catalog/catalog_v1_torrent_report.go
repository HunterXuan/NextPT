package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentReport(ctx context.Context, req *v1.TorrentReportReq) (res *v1.TorrentReportRes, err error) {
	actor := contexts.GetActor(ctx)
	out, err := service.CatalogTorrentUsecase().Report(ctx, actor, req.TorrentReportInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentReportRes{TorrentReportOut: *out}, nil
}
