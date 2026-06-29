package catalog

import (
	"context"

	v1 "server/api/catalog/v1"
	"server/internal/library/contexts"
	"server/internal/model/in/catalogin"
	"server/internal/service"
)

func (c *ControllerV1) TorrentList(ctx context.Context, req *v1.TorrentListReq) (res *v1.TorrentListRes, err error) {
	in := catalogin.TorrentListInp{
		Page:        req.Page,
		Size:        req.Size,
		Keyword:     req.Keyword,
		CategoryIds: req.CategoryIds,
	}

	out, err := service.CatalogTorrentUsecase().List(ctx, contexts.GetActor(ctx), in)
	if err != nil {
		return nil, err
	}

	return &v1.TorrentListRes{
		TorrentListOut: *out,
	}, nil
}
