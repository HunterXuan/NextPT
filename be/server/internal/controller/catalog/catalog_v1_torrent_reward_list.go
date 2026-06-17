package catalog

import (
	"context"

	v1 "server/api/catalog/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) TorrentRewardList(ctx context.Context, req *v1.TorrentRewardListReq) (res *v1.TorrentRewardListRes, err error) {
	out, err := service.CatalogTorrentUsecase().RewardList(ctx, contexts.GetActor(ctx), req.TorrentRewardListInp)
	if err != nil {
		return nil, err
	}
	return &v1.TorrentRewardListRes{
		TorrentRewardListOut: *out,
	}, nil
}
