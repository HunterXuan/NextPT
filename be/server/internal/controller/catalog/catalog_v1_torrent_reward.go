package catalog

import (
	"context"

	v1 "server/api/catalog/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentReward(ctx context.Context, req *v1.TorrentRewardReq) (res *v1.TorrentRewardRes, err error) {
	out, err := service.CatalogTorrentUsecase().Reward(ctx, contexts.GetActor(ctx), req.TorrentRewardInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentRewardRes{
		TorrentRewardOut: *out,
	}, nil
}
