package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentUpdate(ctx context.Context, req *v1.TorrentUpdateReq) (res *v1.TorrentUpdateRes, err error) {
	actor := contexts.GetActor(ctx)
	out, err := service.CatalogTorrentUsecase().Update(ctx, actor, req.TorrentUpdateInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentUpdateRes{TorrentUpdateOut: *out}, nil
}
