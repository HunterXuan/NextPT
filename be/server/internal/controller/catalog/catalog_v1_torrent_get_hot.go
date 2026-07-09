package catalog

import (
	"context"

	v1 "server/api/catalog/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentGetHot(ctx context.Context, req *v1.TorrentGetHotReq) (res *v1.TorrentGetHotRes, err error) {
	out, err := service.CatalogTorrentUsecase().ListHot(ctx, contexts.GetActor(ctx), req.Size)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentGetHotRes{TorrentHotListOut: *out}, nil
}
