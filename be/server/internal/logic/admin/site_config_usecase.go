package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/adminin"
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
	configs, err := service.SiteConfigDomain().AdminListConfigs(ctx, in.Group)
	if err != nil {
		return nil, err
	}
	return &adminout.SiteConfigListOut{Configs: configs}, nil
}

func (s *sAdminSiteConfigUsecase) Update(ctx context.Context, actor *model.Actor, in adminin.SiteConfigUpdateInp) error {
	return service.SiteConfigDomain().AdminUpdateConfig(ctx, in.Group, in.Key, in.Value)
}
