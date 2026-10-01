package catalog

import (
	"context"

	v1 "server/api/catalog/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentCompletionList(ctx context.Context, req *v1.TorrentCompletionListReq) (res *v1.TorrentCompletionListRes, err error) {
	out, err := service.CatalogTorrentUsecase().ListCompletions(ctx, contexts.GetActor(ctx), req.TorrentCompletionListInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentCompletionListRes{TorrentCompletionListOut: *out}, nil
}
