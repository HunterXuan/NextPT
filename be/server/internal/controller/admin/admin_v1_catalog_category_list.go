package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogCategoryList(ctx context.Context, req *v1.CatalogCategoryListReq) (res *v1.CatalogCategoryListRes, err error) {
	out, err := service.AdminCatalogCategoryUsecase().List(ctx, contexts.GetActor(ctx), req.CatalogCategoryListInp)
	if err != nil {
		return nil, err
	}
	res = &v1.CatalogCategoryListRes{CatalogCategoryListOut: *out}
	return
}
