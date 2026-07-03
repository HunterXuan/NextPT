package accounting

import (
	"context"
	"server/internal/library/contexts"
	"server/internal/service"

	v1 "server/api/accounting/v1"
)

func (c *ControllerV1) PeerList(ctx context.Context, req *v1.PeerListReq) (res *v1.PeerListRes, err error) {
	actor := contexts.GetActor(ctx)
	out, err := service.AccountingPeerUsecase().ListMyPeers(ctx, actor, req.PeerListInp)
	if err != nil {
		return nil, err
	}
	return &v1.PeerListRes{PeerListOut: *out}, nil
}
