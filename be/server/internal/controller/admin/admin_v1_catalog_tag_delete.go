package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogTagDelete(ctx context.Context, req *v1.CatalogTagDeleteReq) (res *v1.CatalogTagDeleteRes, err error) {
	err = service.AdminCatalogTagUsecase().DeleteTag(ctx, contexts.GetActor(ctx), req.CatalogTagDeleteInp)
	if err == nil {
		res = &v1.CatalogTagDeleteRes{}
	}
	return
}
