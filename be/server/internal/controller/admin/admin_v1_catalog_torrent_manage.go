package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogTorrentPin(ctx context.Context, req *v1.CatalogTorrentPinReq) (res *v1.CatalogTorrentPinRes, err error) {
	err = service.AdminCatalogTorrentUsecase().Pin(ctx, contexts.GetActor(ctx), req.CatalogTorrentPinInp)
	if err == nil {
		res = &v1.CatalogTorrentPinRes{}
	}
	return
}

func (c *ControllerV1) CatalogTorrentUnpin(ctx context.Context, req *v1.CatalogTorrentUnpinReq) (res *v1.CatalogTorrentUnpinRes, err error) {
	err = service.AdminCatalogTorrentUsecase().Unpin(ctx, contexts.GetActor(ctx), req.CatalogTorrentUnpinInp)
	if err == nil {
		res = &v1.CatalogTorrentUnpinRes{}
	}
	return
}

func (c *ControllerV1) CatalogTorrentFeature(ctx context.Context, req *v1.CatalogTorrentFeatureReq) (res *v1.CatalogTorrentFeatureRes, err error) {
	err = service.AdminCatalogTorrentUsecase().Feature(ctx, contexts.GetActor(ctx), req.CatalogTorrentFeatureInp)
	if err == nil {
		res = &v1.CatalogTorrentFeatureRes{}
	}
	return
}

func (c *ControllerV1) CatalogTorrentUnfeature(ctx context.Context, req *v1.CatalogTorrentUnfeatureReq) (res *v1.CatalogTorrentUnfeatureRes, err error) {
	err = service.AdminCatalogTorrentUsecase().Unfeature(ctx, contexts.GetActor(ctx), req.CatalogTorrentUnfeatureInp)
	if err == nil {
		res = &v1.CatalogTorrentUnfeatureRes{}
	}
	return
}

func (c *ControllerV1) CatalogTorrentPromotion(ctx context.Context, req *v1.CatalogTorrentPromotionReq) (res *v1.CatalogTorrentPromotionRes, err error) {
	err = service.AdminCatalogTorrentUsecase().Promotion(ctx, contexts.GetActor(ctx), req.CatalogTorrentPromotionInp)
	if err == nil {
		res = &v1.CatalogTorrentPromotionRes{}
	}
	return
}

func (c *ControllerV1) CatalogTorrentClearPromotion(ctx context.Context, req *v1.CatalogTorrentClearPromotionReq) (res *v1.CatalogTorrentClearPromotionRes, err error) {
	err = service.AdminCatalogTorrentUsecase().ClearPromotion(ctx, contexts.GetActor(ctx), req.CatalogTorrentClearPromotionInp)
	if err == nil {
		res = &v1.CatalogTorrentClearPromotionRes{}
	}
	return
}
