package catalog

import (
	"context"
	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/catalog/v1"
)

func (c *ControllerV1) TagGroupList(ctx context.Context, req *v1.TagGroupListReq) (res *v1.TagGroupListRes, err error) {
	out, err := service.CatalogTagUsecase().ListTagGroups(ctx, contexts.GetActor(ctx), req.TagGroupListInp)
	if err != nil {
		return nil, err
	}
	return &v1.TagGroupListRes{TagGroupListOut: *out}, nil
}
