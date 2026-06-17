package catalog

import (
	"context"

	v1 "server/api/catalog/v1"
	"server/internal/library/contexts"
	"server/internal/library/httpx"
	"server/internal/model/in/catalogin"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) TorrentDownload(ctx context.Context, req *v1.TorrentDownloadReq) (res *v1.TorrentDownloadRes, err error) {
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

	// 因为已经直接写入响应流，所以我们不再返回正常的 JSON res
	r.ExitAll()
	return nil, nil
}
