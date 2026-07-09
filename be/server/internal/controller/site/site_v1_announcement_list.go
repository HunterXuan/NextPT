package site

import (
	"context"

	v1 "server/api/site/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) AnnouncementList(ctx context.Context, req *v1.AnnouncementListReq) (res *v1.AnnouncementListRes, err error) {
	out, err := service.SiteAnnouncementUsecase().List(ctx, contexts.GetActor(ctx), req.AnnouncementListInp)
	if err != nil {
		return nil, err
	}
	return &v1.AnnouncementListRes{AnnouncementListOut: *out}, nil
}
