package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/adminout"
	"server/internal/service"
)

type sAdminSiteConfigUsecase struct{}

func init() {
	service.RegisterAdminSiteConfigUsecase(NewAdminSiteConfigUsecase())
}

func NewAdminSiteConfigUsecase() *sAdminSiteConfigUsecase {
	return &sAdminSiteConfigUsecase{}
}

func (s *sAdminSiteConfigUsecase) List(ctx context.Context, actor *model.Actor, in adminin.SiteConfigListInp) (*adminout.SiteConfigListOut, error) {
	out, err := service.SiteConfigDomain().AdminListConfigs(ctx, sitein.SiteConfigListInp{Group: in.Group})
	if err != nil {
		return nil, err
	}
	return &adminout.SiteConfigListOut{Configs: out.Configs}, nil
}

func (s *sAdminSiteConfigUsecase) Update(ctx context.Context, actor *model.Actor, in adminin.SiteConfigUpdateInp) error {
	return service.SiteConfigDomain().AdminUpdateConfig(ctx, sitein.SiteConfigUpdateInp{
		Group: in.Group,
		Key:   in.Key,
		Value: in.Value,
	})
}
