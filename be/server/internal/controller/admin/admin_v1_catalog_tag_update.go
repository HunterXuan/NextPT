package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogTagUpdate(ctx context.Context, req *v1.CatalogTagUpdateReq) (res *v1.CatalogTagUpdateRes, err error) {
	err = service.AdminCatalogTagUsecase().UpdateTag(ctx, contexts.GetActor(ctx), req.CatalogTagUpdateInp)
	if err == nil {
		res = &v1.CatalogTagUpdateRes{}
	}
	return
}
