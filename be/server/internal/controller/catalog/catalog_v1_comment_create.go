package catalog

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/catalog/v1"
)

func (c *ControllerV1) CommentCreate(ctx context.Context, req *v1.CommentCreateReq) (res *v1.CommentCreateRes, err error) {
	out, err := service.CatalogCommentUsecase().Create(ctx, contexts.GetActor(ctx), req.CommentCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.CommentCreateRes{Id: out.Id}, nil
}
