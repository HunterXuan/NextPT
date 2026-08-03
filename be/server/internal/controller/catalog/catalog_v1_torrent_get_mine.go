package catalog

import (
	"context"

	v1 "server/api/catalog/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentGetMine(ctx context.Context, req *v1.TorrentGetMineReq) (res *v1.TorrentGetMineRes, err error) {
	out, err := service.CatalogTorrentUsecase().ListMine(ctx, contexts.GetActor(ctx), req.TorrentGetMineInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentGetMineRes{TorrentMineListOut: *out}, nil
}
