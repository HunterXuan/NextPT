package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TrackerClientWhitelistDelete(ctx context.Context, req *v1.TrackerClientWhitelistDeleteReq) (res *v1.TrackerClientWhitelistDeleteRes, err error) {
	if err := service.AdminTrackerClientWhitelistUsecase().Delete(ctx, contexts.GetActor(ctx), req.TrackerClientWhitelistDeleteInp); err != nil {
		return nil, err
	}
	return &v1.TrackerClientWhitelistDeleteRes{}, nil
}
