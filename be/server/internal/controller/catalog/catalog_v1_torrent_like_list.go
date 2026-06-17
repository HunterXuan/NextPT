package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentLikeList(ctx context.Context, req *v1.TorrentLikeListReq) (res *v1.TorrentLikeListRes, err error) {
	actor := contexts.GetActor(ctx)
	out, err := service.CatalogTorrentUsecase().ListLikes(ctx, actor, req.TorrentLikeListInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentLikeListRes{TorrentLikeListOut: *out}, nil
}
