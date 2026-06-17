package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogTorrentDelete(ctx context.Context, req *v1.CatalogTorrentDeleteReq) (res *v1.CatalogTorrentDeleteRes, err error) {
	err = service.AdminCatalogTorrentUsecase().Delete(ctx, contexts.GetActor(ctx), req.CatalogTorrentDeleteInp)
	if err != nil {
		return nil, err
	}
	return &v1.CatalogTorrentDeleteRes{}, nil
}
