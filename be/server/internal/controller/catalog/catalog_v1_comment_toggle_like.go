package catalog

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/catalog/v1"
)

func (c *ControllerV1) CommentToggleLike(ctx context.Context, req *v1.CommentToggleLikeReq) (res *v1.CommentToggleLikeRes, err error) {
	out, err := service.CatalogCommentUsecase().ToggleLike(ctx, contexts.GetActor(ctx), req.CommentToggleLikeInp)
	if err != nil {
		return nil, err
	}
	return &v1.CommentToggleLikeRes{CommentToggleLikeOut: *out}, nil
}
