package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SubtitleUpdate(ctx context.Context, req *v1.SubtitleUpdateReq) (res *v1.SubtitleUpdateRes, err error) {
	actor := contexts.GetActor(ctx)
	err = service.CatalogSubtitleUsecase().Update(ctx, actor, req.SubtitleUpdateInp)
	return nil, err
}
