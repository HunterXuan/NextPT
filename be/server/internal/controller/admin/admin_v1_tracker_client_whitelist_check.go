package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TrackerClientWhitelistCheck(ctx context.Context, req *v1.TrackerClientWhitelistCheckReq) (res *v1.TrackerClientWhitelistCheckRes, err error) {
	out, err := service.AdminTrackerClientWhitelistUsecase().Check(ctx, contexts.GetActor(ctx), req.TrackerClientWhitelistCheckInp)
	if err != nil {
		return nil, err
	}
	return &v1.TrackerClientWhitelistCheckRes{TrackerClientWhitelistCheckOut: *out}, nil
}
