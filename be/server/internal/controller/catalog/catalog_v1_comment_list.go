package catalog

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/catalog/v1"
)

func (c *ControllerV1) CommentList(ctx context.Context, req *v1.CommentListReq) (res *v1.CommentListRes, err error) {
	out, err := service.CatalogCommentUsecase().List(ctx, contexts.GetActor(ctx), req.CommentListInp)
	if err != nil {
		return nil, err
	}
	return &v1.CommentListRes{CommentListOut: *out}, nil
}
