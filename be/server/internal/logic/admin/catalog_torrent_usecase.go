package admin

import (
	"context"
	"encoding/hex"
	"strings"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/adminin"
	"server/internal/model/in/catalogin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/adminout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
)

type sAdminCatalogTorrentUsecase struct{}

func init() {
	service.RegisterAdminCatalogTorrentUsecase(NewAdminCatalogTorrentUsecase())
}

func NewAdminCatalogTorrentUsecase() *sAdminCatalogTorrentUsecase {
	return &sAdminCatalogTorrentUsecase{}
}

func (s *sAdminCatalogTorrentUsecase) ReviewList(ctx context.Context, actor *model.Actor, in adminin.CatalogTorrentReviewListInp) (*adminout.CatalogTorrentReviewListOut, error) {
	torrents, total, err := service.CatalogTorrentDomain().AdminQueryReviewTorrents(ctx, model.CatalogTorrentReviewListOptions{
		Keyword:    in.Keyword,
		CategoryId: in.CategoryId,
		Status:     in.Status,
		Page:       in.Page,
		Size:       in.Size,
	})
	if err != nil {
		return nil, err
	}

	ownerMap := s.catalogTorrentReviewOwnerMap(ctx, torrents)
	items := make([]adminout.CatalogTorrentReviewItem, 0, len(torrents))
	for _, torrent := range torrents {
		items = append(items, adminout.CatalogTorrentReviewItem{
			Id:          torrent.Id,
			Name:        torrent.Name,
			SubTitle:    torrent.SubTitle,
			CategoryId:  torrent.CategoryId,
			Owner:       ownerMap[torrent.OwnerId],
			Size:        torrent.Size,
			FileCount:   torrent.FileCount,
			Status:      torrent.Status,
			SubmittedAt: torrent.SubmittedAt,
			PublishedAt: torrent.PublishedAt,
			ReviewedBy:  torrent.ReviewedBy,
			ReviewedAt:  torrent.ReviewedAt,
			Comment:     torrent.ReviewComment,
			CreatedAt:   torrent.CreatedAt,
		})
	}
	return &adminout.CatalogTorrentReviewListOut{List: items, Total: total, Page: in.Page, Size: in.Size}, nil
}

func (s *sAdminCatalogTorrentUsecase) Update(ctx context.Context, actor *model.Actor, in catalogin.TorrentUpdateInp) error {
	if _, err := service.CatalogTorrentUsecase().Update(ctx, actor, in); err != nil {
		return err
	}
	s.recordTorrentUpdateAudit(ctx, actor, in.Id, consts.SiteAuditOperationEdit, nil)
	return nil
}

func (s *sAdminCatalogTorrentUsecase) Approve(ctx context.Context, actor *model.Actor, in adminin.CatalogTorrentApproveInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	now := gtime.Now()
	var torrent *entity.CatalogTorrent
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		current, err := service.CatalogTorrentDomain().GetTorrentByIdForUpdate(ctx, in.Id)
		if err != nil {
			return err
		}
		if current.Status != consts.CatalogTorrentStatusPending {
			return gerror.New(gi18n.T(ctx, "admin.catalog.torrent.review_state_changed"))
		}
		torrent = current
		promotion := service.CatalogTorrentDomain().PickNewTorrentPromotion(ctx, current.Size, now)
		updated, err := service.CatalogTorrentDomain().AdminApproveTorrent(ctx, in.Id, actor.Id, in.Comment, now, promotion.SpState, promotion.SpExpireAt)
		if err != nil {
			return err
		}
		if !updated {
			return gerror.New(gi18n.T(ctx, "admin.catalog.torrent.review_state_changed"))
		}
		return nil
	})
	if err != nil {
		return err
	}

	s.invalidateReviewedTorrentCaches(ctx, torrent, true)
	s.notifyTorrentReview(ctx, torrent, true, strings.TrimSpace(in.Comment))
	s.recordTorrentUpdateAudit(ctx, actor, in.Id, consts.SiteAuditOperationApprove, map[string]any{"comment": strings.TrimSpace(in.Comment)})
	return nil
}

func (s *sAdminCatalogTorrentUsecase) Reject(ctx context.Context, actor *model.Actor, in adminin.CatalogTorrentRejectInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	comment := strings.TrimSpace(in.Comment)
	if comment == "" {
		return gerror.New(gi18n.T(ctx, "admin.catalog.torrent.review_comment_required"))
	}
	now := gtime.Now()
	var torrent *entity.CatalogTorrent
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		current, err := service.CatalogTorrentDomain().GetTorrentByIdForUpdate(ctx, in.Id)
		if err != nil {
			return err
		}
		if current.Status != consts.CatalogTorrentStatusPending {
			return gerror.New(gi18n.T(ctx, "admin.catalog.torrent.review_state_changed"))
		}
		torrent = current
		updated, err := service.CatalogTorrentDomain().AdminRejectTorrent(ctx, in.Id, actor.Id, comment, now)
		if err != nil {
			return err
		}
		if !updated {
			return gerror.New(gi18n.T(ctx, "admin.catalog.torrent.review_state_changed"))
		}
		return nil
	})
	if err != nil {
		return err
	}

	s.invalidateReviewedTorrentCaches(ctx, torrent, false)
	s.notifyTorrentReview(ctx, torrent, false, comment)
	s.recordTorrentUpdateAudit(ctx, actor, in.Id, consts.SiteAuditOperationReject, map[string]any{"comment": comment})
	return nil
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
		"status":     torrent.Status,
		"infoHash":   hex.EncodeToString(torrent.InfoHash),
	}
}

func (s *sAdminCatalogTorrentUsecase) catalogTorrentReviewOwnerMap(ctx context.Context, torrents []entity.CatalogTorrent) map[uint64]model.IamUserSummary {
	ownerMap := make(map[uint64]model.IamUserSummary)
	ownerIds := make([]uint64, 0, len(torrents))
	for _, torrent := range torrents {
		if torrent.OwnerId == 0 {
			continue
		}
		if _, exists := ownerMap[torrent.OwnerId]; exists {
			continue
		}
		ownerMap[torrent.OwnerId] = model.IamUserSummary{Id: torrent.OwnerId}
		ownerIds = append(ownerIds, torrent.OwnerId)
	}
	if len(ownerIds) == 0 {
		return ownerMap
	}

	users, err := service.IamUserDomain().GetUsersByIds(ctx, ownerIds)
	if err == nil {
		for _, user := range users {
			summary := ownerMap[user.Id]
			summary.Username = user.Username
			ownerMap[user.Id] = summary
		}
	}
	profiles, err := service.IamUserDomain().GetUserProfilesByUserIds(ctx, ownerIds)
	if err == nil {
		for _, profile := range profiles {
			summary := ownerMap[profile.UserId]
			summary.Avatar = profile.Avatar
			ownerMap[profile.UserId] = summary
		}
	}
	return ownerMap
}

func (s *sAdminCatalogTorrentUsecase) invalidateReviewedTorrentCaches(ctx context.Context, torrent *entity.CatalogTorrent, published bool) {
	if torrent == nil {
		return
	}
	if len(torrent.InfoHash) > 0 {
		cacheKey := service.SysCache().KeyCatalogTorrentInfoHash(ctx, hex.EncodeToString(torrent.InfoHash))
		_, _ = gcache.Remove(ctx, cacheKey)
		_ = service.SysCache().PublishInvalidate(ctx, cacheKey)
	}
	if published {
		cacheKey := service.SysCache().KeyCatalogHotTorrents(ctx)
		_, _ = gcache.Remove(ctx, cacheKey)
		_ = service.SysCache().PublishInvalidate(ctx, cacheKey)
	}
}

func (s *sAdminCatalogTorrentUsecase) notifyTorrentReview(ctx context.Context, torrent *entity.CatalogTorrent, approved bool, comment string) {
	if torrent == nil {
		return
	}
	notify := sitein.MessageNotifyInp{
		ReceiverId: torrent.OwnerId,
		TargetType: consts.SiteMessageTargetTypeCatalogTorrent,
		TargetId:   torrent.Id,
	}
	if approved {
		notify.TitleKey = "site.message.catalog_torrent_review.approved.title"
		notify.ContentKey = "site.message.catalog_torrent_review.approved.content"
		notify.ContentArgs = []any{torrent.Name}
	} else {
		notify.TitleKey = "site.message.catalog_torrent_review.rejected.title"
		notify.ContentKey = "site.message.catalog_torrent_review.rejected.content"
		notify.ContentArgs = []any{torrent.Name, comment}
	}
	service.SiteMessageUsecase().Notify(ctx, notify)
}
