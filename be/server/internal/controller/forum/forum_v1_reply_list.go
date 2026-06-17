package forum

import (
	"context"

	v1 "server/api/forum/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) ReplyList(ctx context.Context, req *v1.ReplyListReq) (res *v1.ReplyListRes, err error) {
	res = &v1.ReplyListRes{}
	out, err := service.ForumReplyUsecase().List(ctx, contexts.GetActor(ctx), req.ReplyListInp)
	if err != nil {
		return nil, err
	}
	res.ReplyListOut = *out
	return
}
