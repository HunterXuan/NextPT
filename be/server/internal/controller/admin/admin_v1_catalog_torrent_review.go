package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogTorrentReviewList(ctx context.Context, req *v1.CatalogTorrentReviewListReq) (res *v1.CatalogTorrentReviewListRes, err error) {
	out, err := service.AdminCatalogTorrentUsecase().ReviewList(ctx, contexts.GetActor(ctx), req.CatalogTorrentReviewListInp)
	if err != nil {
		return nil, err
	}
	return &v1.CatalogTorrentReviewListRes{CatalogTorrentReviewListOut: *out}, nil
}

func (c *ControllerV1) CatalogTorrentUpdate(ctx context.Context, req *v1.CatalogTorrentUpdateReq) (res *v1.CatalogTorrentUpdateRes, err error) {
	if err := service.AdminCatalogTorrentUsecase().Update(ctx, contexts.GetActor(ctx), req.TorrentUpdateInp); err != nil {
		return nil, err
	}
	return &v1.CatalogTorrentUpdateRes{}, nil
}

func (c *ControllerV1) CatalogTorrentApprove(ctx context.Context, req *v1.CatalogTorrentApproveReq) (res *v1.CatalogTorrentApproveRes, err error) {
	if err := service.AdminCatalogTorrentUsecase().Approve(ctx, contexts.GetActor(ctx), req.CatalogTorrentApproveInp); err != nil {
		return nil, err
	}
	return &v1.CatalogTorrentApproveRes{}, nil
}

func (c *ControllerV1) CatalogTorrentReject(ctx context.Context, req *v1.CatalogTorrentRejectReq) (res *v1.CatalogTorrentRejectRes, err error) {
	if err := service.AdminCatalogTorrentUsecase().Reject(ctx, contexts.GetActor(ctx), req.CatalogTorrentRejectInp); err != nil {
		return nil, err
	}
	return &v1.CatalogTorrentRejectRes{}, nil
}
