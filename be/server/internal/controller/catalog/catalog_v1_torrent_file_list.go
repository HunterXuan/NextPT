package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentFileList(ctx context.Context, req *v1.TorrentFileListReq) (res *v1.TorrentFileListRes, err error) {
	out, err := service.CatalogTorrentUsecase().ListFiles(ctx, contexts.GetActor(ctx), req.TorrentFileListInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentFileListRes{TorrentFileListOut: *out}, nil
}
