package forum

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/forum/v1"
)

func (c *ControllerV1) TopicReport(ctx context.Context, req *v1.TopicReportReq) (res *v1.TopicReportRes, err error) {
	err = service.ForumTopicUsecase().ReportTopic(ctx, contexts.GetActor(ctx), req.TopicReportInp)
	return
}
