package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentToggleLike(ctx context.Context, req *v1.TorrentToggleLikeReq) (res *v1.TorrentToggleLikeRes, err error) {
	actor := contexts.GetActor(ctx)
	out, err := service.CatalogTorrentUsecase().ToggleLike(ctx, actor, req.TorrentToggleLikeInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentToggleLikeRes{TorrentToggleLikeOut: *out}, nil
}
