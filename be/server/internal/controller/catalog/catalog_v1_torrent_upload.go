package catalog

import (
	"context"

	v1 "server/api/catalog/v1"
	"server/internal/library/contexts"
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"
	"server/internal/service"
)

func (c *ControllerV1) TorrentUpload(ctx context.Context, req *v1.TorrentUploadReq) (res *v1.TorrentUploadRes, err error) {
	in := catalogin.TorrentUploadInp{
		File:        req.File,
		Name:        req.Name,
		SubTitle:    req.SubTitle,
		CategoryId:  req.CategoryId,
		Description: req.Description,
		Anonymous:   req.Anonymous,
	}

	out, err := service.CatalogTorrentUsecase().Upload(ctx, contexts.GetActor(ctx), in)
	if err != nil {
		return nil, err
	}

	return &v1.TorrentUploadRes{
		TorrentUploadOut: catalogout.TorrentUploadOut{
			TorrentId: out.TorrentId,
			InfoHash:  out.InfoHash,
		},
	}, nil
}
