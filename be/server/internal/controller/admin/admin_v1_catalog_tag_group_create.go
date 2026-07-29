package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogTagGroupCreate(ctx context.Context, req *v1.CatalogTagGroupCreateReq) (res *v1.CatalogTagGroupCreateRes, err error) {
	err = service.AdminCatalogTagUsecase().CreateGroup(ctx, contexts.GetActor(ctx), req.CatalogTagGroupCreateInp)
	if err == nil {
		res = &v1.CatalogTagGroupCreateRes{}
	}
	return
}
