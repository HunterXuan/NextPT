package catalog

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/catalog/v1"
)

func (c *ControllerV1) CommentReport(ctx context.Context, req *v1.CommentReportReq) (res *v1.CommentReportRes, err error) {
	err = service.CatalogCommentUsecase().Report(ctx, contexts.GetActor(ctx), req.CommentReportInp)
	if err != nil {
		return nil, err
	}
	return &v1.CommentReportRes{}, nil
}
