package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogCategoryDelete(ctx context.Context, req *v1.CatalogCategoryDeleteReq) (res *v1.CatalogCategoryDeleteRes, err error) {
	err = service.AdminCatalogCategoryUsecase().Delete(ctx, contexts.GetActor(ctx), req.CatalogCategoryDeleteInp)
	if err == nil {
		res = &v1.CatalogCategoryDeleteRes{}
	}
	return
}
