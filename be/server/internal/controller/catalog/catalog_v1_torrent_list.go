package catalog

import (
	"context"

	v1 "server/api/catalog/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentList(ctx context.Context, req *v1.TorrentListReq) (res *v1.TorrentListRes, err error) {
	out, err := service.CatalogTorrentUsecase().List(ctx, contexts.GetActor(ctx), req.TorrentListInp)
	if err != nil {
		return nil, err
	}

	return &v1.TorrentListRes{
		TorrentListOut: *out,
	}, nil
}
