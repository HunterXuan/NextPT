package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogCategoryCreate(ctx context.Context, req *v1.CatalogCategoryCreateReq) (res *v1.CatalogCategoryCreateRes, err error) {
	err = service.AdminCatalogCategoryUsecase().Create(ctx, contexts.GetActor(ctx), req.CatalogCategoryCreateInp)
	if err == nil {
		res = &v1.CatalogCategoryCreateRes{}
	}
	return
}
