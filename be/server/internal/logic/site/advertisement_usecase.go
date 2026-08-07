package site

import (
	"context"
	"time"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/out/siteout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
)

type sSiteAdvertisementUsecase struct{}

func init() {
	service.RegisterSiteAdvertisementUsecase(NewSiteAdvertisementUsecase())
}

func NewSiteAdvertisementUsecase() *sSiteAdvertisementUsecase {
	return &sSiteAdvertisementUsecase{}
}

func (s *sSiteAdvertisementUsecase) List(ctx context.Context) (*siteout.AdvertisementListOut, error) {
	advertisements, err := s.loadAdvertisementsCache(ctx)
	if err != nil {
		return nil, err
	}
	return &siteout.AdvertisementListOut{Placements: advertisements.Enabled()}, nil
}

func (s *sSiteAdvertisementUsecase) loadAdvertisementsCache(ctx context.Context) (model.SiteAdvertisements, error) {
	cacheKey := service.SysCache().KeySiteConfigFullPath(ctx, consts.SiteConfigSiteAdvertisements)
	value, err := gcache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		return s.loadAdvertisements(ctx)
	}, 10*time.Minute)
	if err != nil {
		return nil, err
	}
	if advertisements, ok := value.Val().(model.SiteAdvertisements); ok {
		return advertisements, nil
	}
	return s.loadAdvertisements(ctx)
}

func (s *sSiteAdvertisementUsecase) loadAdvertisements(ctx context.Context) (model.SiteAdvertisements, error) {
	var advertisements model.SiteAdvertisements
	if err := service.SiteConfigDomain().GetByPath(ctx, consts.SiteConfigSiteAdvertisements).Scan(&advertisements); err != nil {
		return nil, gerror.New("invalid site advertisements config")
	}
	if err := advertisements.Validate(); err != nil {
		return nil, err
	}
	advertisements = advertisements.Normalized()
	if err := advertisements.Validate(); err != nil {
		return nil, err
	}
	return advertisements, nil
}
