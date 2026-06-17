package catalog

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/catalog/v1"
)

func (c *ControllerV1) TorrentBookmark(ctx context.Context, req *v1.TorrentBookmarkReq) (res *v1.TorrentBookmarkRes, err error) {
	err = service.CatalogTorrentUsecase().Bookmark(ctx, contexts.GetActor(ctx), req.TorrentBookmarkInp)
	return
}
