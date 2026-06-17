package tracker

import (
	"context"

	v1 "server/api/tracker/v1"
	"server/internal/library/contexts"
	"server/internal/library/httpx"
	"server/internal/model/in/catalogin"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) Download(ctx context.Context, req *v1.DownloadReq) (res *v1.DownloadRes, err error) {
	in := catalogin.TorrentDownloadInp{
		Id: req.Id,
	}

	out, err := service.CatalogTorrentUsecase().Download(ctx, contexts.GetActor(ctx), in)
	if err != nil {
		return nil, err
	}

	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "application/x-bittorrent")
	r.Response.Header().Set("Content-Disposition", httpx.HeaderContentDispositionAttachment(out.FileName))
	r.Response.Write(out.Bytes)

	r.ExitAll()
	return nil, nil
}
