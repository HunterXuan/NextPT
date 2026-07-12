package admin

import (
	"context"

	v1 "server/api/admin/v1"
	"server/internal/library/contexts"
	"server/internal/service"
)

func (c *ControllerV1) CatalogRequestComplete(ctx context.Context, req *v1.CatalogRequestCompleteReq) (res *v1.CatalogRequestCompleteRes, err error) {
	err = service.AdminCatalogRequestUsecase().Complete(ctx, contexts.GetActor(ctx), req.RequestCompleteInp)
	if err != nil {
		return nil, err
	}
	return &v1.CatalogRequestCompleteRes{}, nil
}

func (c *ControllerV1) CatalogRequestCancel(ctx context.Context, req *v1.CatalogRequestCancelReq) (res *v1.CatalogRequestCancelRes, err error) {
	err = service.AdminCatalogRequestUsecase().Cancel(ctx, contexts.GetActor(ctx), req.RequestCancelInp)
	if err != nil {
		return nil, err
	}
	return &v1.CatalogRequestCancelRes{}, nil
}
