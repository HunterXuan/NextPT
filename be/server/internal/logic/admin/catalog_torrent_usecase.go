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

func (s *sAdminCatalogTorrentUsecase) Pin(ctx context.Context, actor *model.Actor, in adminin.CatalogTorrentPinInp) error {
	if _, err := service.CatalogTorrentDomain().GetTorrentById(ctx, in.Id); err != nil {
		return err
	}
	if err := service.CatalogTorrentDomain().AdminSetTorrentPinned(ctx, in.Id, true, in.PinWeight); err != nil {
		return err
	}
	s.recordTorrentUpdateAudit(ctx, actor, in.Id, consts.SiteAuditOperationPin, map[string]any{
		"pinWeight": in.PinWeight,
	})
	return nil
}

func (s *sAdminCatalogTorrentUsecase) Unpin(ctx context.Context, actor *model.Actor, in adminin.CatalogTorrentUnpinInp) error {
	if _, err := service.CatalogTorrentDomain().GetTorrentById(ctx, in.Id); err != nil {
		return err
	}
	if err := service.CatalogTorrentDomain().AdminSetTorrentPinned(ctx, in.Id, false, 0); err != nil {
		return err
	}
	s.recordTorrentUpdateAudit(ctx, actor, in.Id, consts.SiteAuditOperationUnpin, nil)
	return nil
}

func (s *sAdminCatalogTorrentUsecase) Feature(ctx context.Context, actor *model.Actor, in adminin.CatalogTorrentFeatureInp) error {
	if _, err := service.CatalogTorrentDomain().GetTorrentById(ctx, in.Id); err != nil {
		return err
	}
	if err := service.CatalogTorrentDomain().AdminSetTorrentFeatured(ctx, in.Id, true); err != nil {
		return err
	}
	s.recordTorrentUpdateAudit(ctx, actor, in.Id, consts.SiteAuditOperationFeature, nil)
	return nil
}

func (s *sAdminCatalogTorrentUsecase) Unfeature(ctx context.Context, actor *model.Actor, in adminin.CatalogTorrentUnfeatureInp) error {
	if _, err := service.CatalogTorrentDomain().GetTorrentById(ctx, in.Id); err != nil {
		return err
	}
	if err := service.CatalogTorrentDomain().AdminSetTorrentFeatured(ctx, in.Id, false); err != nil {
		return err
	}
	s.recordTorrentUpdateAudit(ctx, actor, in.Id, consts.SiteAuditOperationUnfeature, nil)
	return nil
}

func (s *sAdminCatalogTorrentUsecase) Promotion(ctx context.Context, actor *model.Actor, in adminin.CatalogTorrentPromotionInp) error {
	if _, err := service.CatalogTorrentDomain().GetTorrentById(ctx, in.Id); err != nil {
		return err
	}
	if err := service.CatalogTorrentDomain().AdminSetTorrentPromotion(ctx, in.Id, in.SpState, in.SpExpireAt); err != nil {
		return err
	}
	detail := map[string]any{
		"spState": in.SpState,
	}
	if in.SpExpireAt != nil {
		detail["spExpireAt"] = in.SpExpireAt.String()
	}
	s.recordTorrentUpdateAudit(ctx, actor, in.Id, consts.SiteAuditOperationSetPromotion, detail)
	return nil
}

func (s *sAdminCatalogTorrentUsecase) ClearPromotion(ctx context.Context, actor *model.Actor, in adminin.CatalogTorrentClearPromotionInp) error {
	if _, err := service.CatalogTorrentDomain().GetTorrentById(ctx, in.Id); err != nil {
		return err
	}
	if err := service.CatalogTorrentDomain().AdminSetTorrentPromotion(ctx, in.Id, consts.ResourceTorrentSpNormal, nil); err != nil {
		return err
	}
	s.recordTorrentUpdateAudit(ctx, actor, in.Id, consts.SiteAuditOperationClearPromotion, nil)
	return nil
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

func (s *sAdminCatalogTorrentUsecase) recordTorrentUpdateAudit(ctx context.Context, actor *model.Actor, torrentId uint64, operation string, detail map[string]any) {
	if detail == nil {
		detail = map[string]any{}
	}
	detail["operation"] = operation
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeCatalogTorrent,
		TargetId:   torrentId,
		Level:      consts.SiteAuditLevelImportant,
		Detail:     detail,
	})
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
