package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogTagGroupDelete(ctx context.Context, req *v1.CatalogTagGroupDeleteReq) (res *v1.CatalogTagGroupDeleteRes, err error) {
	err = service.AdminCatalogTagUsecase().DeleteGroup(ctx, contexts.GetActor(ctx), req.CatalogTagGroupDeleteInp)
	if err == nil {
		res = &v1.CatalogTagGroupDeleteRes{}
	}
	return
}
