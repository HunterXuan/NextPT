package admin

import (
	"context"
	"encoding/hex"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
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
	torrent, err := service.CatalogTorrentDomain().GetTorrentById(ctx, in.Id)
	if err != nil {
		return err
	}
	if err := service.CatalogTorrentUsecase().HardDeleteTorrent(ctx, in.Id); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionDelete,
		TargetType: consts.SiteAuditTargetTypeCatalogTorrent,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"snapshot": s.catalogTorrentAuditSnapshot(torrent),
		},
	})
	return nil
}

func (s *sAdminCatalogTorrentUsecase) catalogTorrentAuditSnapshot(torrent *entity.CatalogTorrent) map[string]any {
	if torrent == nil {
		return nil
	}
	return map[string]any{
		"id":         torrent.Id,
		"name":       torrent.Name,
		"categoryId": torrent.CategoryId,
		"ownerId":    torrent.OwnerId,
		"size":       torrent.Size,
		"infoHash":   hex.EncodeToString(torrent.InfoHash),
	}
}
