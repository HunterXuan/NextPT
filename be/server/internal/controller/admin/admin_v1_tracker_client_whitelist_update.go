package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TrackerClientWhitelistUpdate(ctx context.Context, req *v1.TrackerClientWhitelistUpdateReq) (res *v1.TrackerClientWhitelistUpdateRes, err error) {
	if err := service.AdminTrackerClientWhitelistUsecase().Update(ctx, contexts.GetActor(ctx), req.TrackerClientWhitelistUpdateInp); err != nil {
		return nil, err
	}
	return &v1.TrackerClientWhitelistUpdateRes{}, nil
}
