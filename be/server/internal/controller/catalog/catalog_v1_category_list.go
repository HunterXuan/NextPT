package catalog

import (
	"context"
	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/catalog/v1"
)

func (c *ControllerV1) CategoryList(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	out, err := service.CatalogCategoryUsecase().ListCategories(ctx, contexts.GetActor(ctx), req.CategoryListInp)
	if err != nil {
		return nil, err
	}
	return &v1.CategoryListRes{CategoryListOut: *out}, nil
}
