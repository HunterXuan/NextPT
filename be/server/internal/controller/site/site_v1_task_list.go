package site

import (
	"context"

	v1 "server/api/site/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TaskList(ctx context.Context, req *v1.TaskListReq) (res *v1.TaskListRes, err error) {
	out, err := service.SiteTaskUsecase().List(ctx, contexts.GetActor(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.TaskListRes{TaskListOut: *out}, nil
}
