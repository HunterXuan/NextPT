package admin

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
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
	if err := service.CatalogTorrentUsecase().HardDeleteTorrent(ctx, in.Id); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionDelete,
		TargetType: consts.SiteAuditTargetTypeCatalogTorrent,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
	})
	return nil
}
