package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogTagGroupUpdate(ctx context.Context, req *v1.CatalogTagGroupUpdateReq) (res *v1.CatalogTagGroupUpdateRes, err error) {
	err = service.AdminCatalogTagUsecase().UpdateGroup(ctx, contexts.GetActor(ctx), req.CatalogTagGroupUpdateInp)
	if err == nil {
		res = &v1.CatalogTagGroupUpdateRes{}
	}
	return
}
