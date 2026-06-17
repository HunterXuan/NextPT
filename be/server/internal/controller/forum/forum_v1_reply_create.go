package forum

import (
	"context"

	v1 "server/api/forum/v1"
	"server/internal/library/contexts"
	"server/internal/model/out/forumout"
	"server/internal/service"
)

func (c *ControllerV1) ReplyCreate(ctx context.Context, req *v1.ReplyCreateReq) (res *v1.ReplyCreateRes, err error) {
	id, err := service.ForumReplyUsecase().Create(ctx, contexts.GetActor(ctx), req.ReplyCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.ReplyCreateRes{
		ReplyCreateOut: forumout.ReplyCreateOut{Id: id},
	}, nil
}
