package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TrackerClientWhitelistList(ctx context.Context, req *v1.TrackerClientWhitelistListReq) (res *v1.TrackerClientWhitelistListRes, err error) {
	out, err := service.AdminTrackerClientWhitelistUsecase().List(ctx, contexts.GetActor(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.TrackerClientWhitelistListRes{TrackerClientWhitelistListOut: *out}, nil
}
