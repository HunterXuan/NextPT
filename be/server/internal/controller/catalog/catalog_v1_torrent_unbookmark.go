package catalog

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/catalog/v1"
)

func (c *ControllerV1) TorrentUnbookmark(ctx context.Context, req *v1.TorrentUnbookmarkReq) (res *v1.TorrentUnbookmarkRes, err error) {
	err = service.CatalogTorrentUsecase().Unbookmark(ctx, contexts.GetActor(ctx), req.TorrentUnbookmarkInp)
	return
}
