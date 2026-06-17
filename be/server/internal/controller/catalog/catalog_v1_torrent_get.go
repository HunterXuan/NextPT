package catalog

import (
	"context"

	v1 "server/api/catalog/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentGet(ctx context.Context, req *v1.TorrentGetReq) (res *v1.TorrentGetRes, err error) {
	out, err := service.CatalogTorrentUsecase().GetTorrent(ctx, contexts.GetActor(ctx), req.TorrentGetInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentGetRes{TorrentDetailOut: *out}, nil
}
