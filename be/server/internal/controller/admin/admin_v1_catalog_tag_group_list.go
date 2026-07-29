package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogTagGroupList(ctx context.Context, req *v1.CatalogTagGroupListReq) (res *v1.CatalogTagGroupListRes, err error) {
	out, err := service.AdminCatalogTagUsecase().List(ctx, contexts.GetActor(ctx), req.CatalogTagGroupListInp)
	if err != nil {
		return nil, err
	}
	return &v1.CatalogTagGroupListRes{CatalogTagGroupListOut: *out}, nil
}
