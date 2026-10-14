package admin

import (
	"context"

	"server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TrackerClientWhitelistCreate(ctx context.Context, req *v1.TrackerClientWhitelistCreateReq) (res *v1.TrackerClientWhitelistCreateRes, err error) {
	out, err := service.AdminTrackerClientWhitelistUsecase().Create(ctx, contexts.GetActor(ctx), req.TrackerClientWhitelistCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.TrackerClientWhitelistCreateRes{TrackerClientWhitelistCreateOut: *out}, nil
}
