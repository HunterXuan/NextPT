package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentSubtitleList(ctx context.Context, req *v1.TorrentSubtitleListReq) (res *v1.TorrentSubtitleListRes, err error) {
	actor := contexts.GetActor(ctx)
	out, err := service.CatalogSubtitleUsecase().ListByTorrent(ctx, actor, req.TorrentSubtitleListInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentSubtitleListRes{SubtitleListOut: *out}, nil
}
