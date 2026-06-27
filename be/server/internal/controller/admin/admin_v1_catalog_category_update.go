package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogCategoryUpdate(ctx context.Context, req *v1.CatalogCategoryUpdateReq) (res *v1.CatalogCategoryUpdateRes, err error) {
	err = service.AdminCatalogCategoryUsecase().Update(ctx, contexts.GetActor(ctx), req.CatalogCategoryUpdateInp)
	if err == nil {
		res = &v1.CatalogCategoryUpdateRes{}
	}
	return
}
