package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogTagCreate(ctx context.Context, req *v1.CatalogTagCreateReq) (res *v1.CatalogTagCreateRes, err error) {
	err = service.AdminCatalogTagUsecase().CreateTag(ctx, contexts.GetActor(ctx), req.CatalogTagCreateInp)
	if err == nil {
		res = &v1.CatalogTagCreateRes{}
	}
	return
}
