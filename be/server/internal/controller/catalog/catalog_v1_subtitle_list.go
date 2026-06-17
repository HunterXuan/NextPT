package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SubtitleList(ctx context.Context, req *v1.SubtitleListReq) (res *v1.SubtitleListRes, err error) {
	actor := contexts.GetActor(ctx)
	out, err := service.CatalogSubtitleUsecase().List(ctx, actor, req.SubtitleListInp)
	if err != nil {
		return nil, err
	}
	return &v1.SubtitleListRes{SubtitleListOut: *out}, nil
}
