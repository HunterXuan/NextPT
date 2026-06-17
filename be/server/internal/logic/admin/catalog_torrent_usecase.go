package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/service"
)

type sAdminCatalogTorrentUsecase struct{}

func init() {
	service.RegisterAdminCatalogTorrentUsecase(NewAdminCatalogTorrentUsecase())
}

func NewAdminCatalogTorrentUsecase() *sAdminCatalogTorrentUsecase {
	return &sAdminCatalogTorrentUsecase{}
}

func (s *sAdminCatalogTorrentUsecase) Delete(ctx context.Context, actor *model.Actor, in adminin.CatalogTorrentDeleteInp) error {
	return service.CatalogTorrentUsecase().HardDeleteTorrent(ctx, in.Id)
}
