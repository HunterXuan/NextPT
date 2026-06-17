package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/library/httpx"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) SubtitleDownload(ctx context.Context, req *v1.SubtitleDownloadReq) (res *v1.SubtitleDownloadRes, err error) {

	actor := contexts.GetActor(ctx)
	out, err := service.CatalogSubtitleUsecase().Download(ctx, actor, req.SubtitleDownloadInp)
	if err != nil {
		return nil, err
	}
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "application/octet-stream")
	r.Response.Header().Set("Content-Disposition", httpx.HeaderContentDispositionAttachment(out.FileName))
	r.Response.Write(out.Bytes)
	r.Exit()
	return nil, nil
}
