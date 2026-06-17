package forum

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/forum/v1"
)

func (c *ControllerV1) ReplyReport(ctx context.Context, req *v1.ReplyReportReq) (res *v1.ReplyReportRes, err error) {
	err = service.ForumReplyUsecase().ReportReply(ctx, contexts.GetActor(ctx), req.ReplyReportInp)
	return
}
