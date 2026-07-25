package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentMetadataSearch(ctx context.Context, req *v1.TorrentMetadataSearchReq) (res *v1.TorrentMetadataSearchRes, err error) {
	out, err := service.CatalogMetadataUsecase().Search(ctx, contexts.GetActor(ctx), req.TorrentMetadataSearchInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentMetadataSearchRes{TorrentMetadataSearchOut: *out}, nil
}
