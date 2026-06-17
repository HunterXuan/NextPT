package catalog

import (
	"context"

	v1 "server/api/catalog/v1"

	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) SubtitleUpload(ctx context.Context, req *v1.SubtitleUploadReq) (res *v1.SubtitleUploadRes, err error) {
	actor := contexts.GetActor(ctx)
	id, err := service.CatalogSubtitleUsecase().Upload(ctx, actor, req.SubtitleUploadInp)
	if err != nil {
		return nil, err
	}
	return &v1.SubtitleUploadRes{Id: id}, nil
}
