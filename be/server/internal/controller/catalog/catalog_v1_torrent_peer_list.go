package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentPeerList(ctx context.Context, req *v1.TorrentPeerListReq) (res *v1.TorrentPeerListRes, err error) {
	out, err := service.CatalogTorrentUsecase().ListPeers(ctx, contexts.GetActor(ctx), req.TorrentPeerListInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentPeerListRes{TorrentPeerListOut: *out}, nil
}
