package site

import (
	"context"

	v1 "server/api/site/v1"
	"server/internal/service"
)

func (c *ControllerV1) AdvertisementList(ctx context.Context, req *v1.AdvertisementListReq) (res *v1.AdvertisementListRes, err error) {
	out, err := service.SiteAdvertisementUsecase().List(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.AdvertisementListRes{AdvertisementListOut: *out}, nil
}
