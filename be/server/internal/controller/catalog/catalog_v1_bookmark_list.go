package catalog

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/catalog/v1"
)

func (c *ControllerV1) BookmarkList(ctx context.Context, req *v1.BookmarkListReq) (res *v1.BookmarkListRes, err error) {
	out, err := service.CatalogTorrentUsecase().ListBookmarkedTorrents(ctx, contexts.GetActor(ctx), req.TorrentBookmarkListInp)
	if err == nil {
		res = &v1.BookmarkListRes{TorrentBookmarkListOut: *out}
	}
	return
}
