package site

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/sitein"
	"server/internal/model/out/siteout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sSiteAnnouncementUsecase struct{}

func init() {
	service.RegisterSiteAnnouncementUsecase(NewSiteAnnouncementUsecase())
}

func NewSiteAnnouncementUsecase() *sSiteAnnouncementUsecase {
	return &sSiteAnnouncementUsecase{}
}

func (s *sSiteAnnouncementUsecase) List(ctx context.Context, actor *model.Actor, in sitein.AnnouncementListInp) (*siteout.AnnouncementListOut, error) {
	actorId := s.actorId(actor)
	list, total, err := service.SiteAnnouncementDomain().ListPublished(ctx, actorId, in)
	if err != nil {
		return nil, err
	}
	readMap, err := service.SiteAnnouncementDomain().QueryReadIds(ctx, actorId, s.announcementIds(list))
	if err != nil {
		return nil, err
	}
	return &siteout.AnnouncementListOut{
		List:  s.buildAnnouncementItems(list, readMap),
		Total: total,
		Page:  in.Page,
		Size:  in.Size,
	}, nil
}

func (s *sSiteAnnouncementUsecase) MarkRead(ctx context.Context, actor *model.Actor, in sitein.AnnouncementReadInp) error {
	return service.SiteAnnouncementDomain().MarkRead(ctx, s.actorId(actor), in.Id)
}

func (s *sSiteAnnouncementUsecase) AdminList(ctx context.Context, actor *model.Actor, in sitein.AdminAnnouncementListInp) (*siteout.AnnouncementListOut, error) {
	if err := s.validateStatusPtr(in.Status); err != nil {
		return nil, err
	}
	list, total, err := service.SiteAnnouncementDomain().AdminList(ctx, in)
	if err != nil {
		return nil, err
	}
	return &siteout.AnnouncementListOut{
		List:  s.buildAnnouncementItems(list, nil),
		Total: total,
		Page:  in.Page,
		Size:  in.Size,
	}, nil
}

func (s *sSiteAnnouncementUsecase) AdminCreate(ctx context.Context, actor *model.Actor, in sitein.AdminAnnouncementCreateInp) (*siteout.AnnouncementCreateOut, error) {
	if err := s.validateStatus(in.Status); err != nil {
		return nil, err
	}
	id, err := service.SiteAnnouncementDomain().AdminCreate(ctx, in, s.actorId(actor))
	if err != nil {
		return nil, err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionCreate,
		TargetType: consts.SiteAuditTargetTypeSiteAnnouncement,
		TargetId:   id,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"title":  in.Title,
			"status": in.Status,
		},
	})
	return &siteout.AnnouncementCreateOut{Id: id}, nil
}

func (s *sSiteAnnouncementUsecase) AdminUpdate(ctx context.Context, actor *model.Actor, in sitein.AdminAnnouncementUpdateInp) error {
	if err := s.validateStatus(in.Status); err != nil {
		return err
	}
	if err := service.SiteAnnouncementDomain().AdminUpdate(ctx, in, s.actorId(actor)); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeSiteAnnouncement,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"title":  in.Title,
			"status": in.Status,
		},
	})
	return nil
}

func (s *sSiteAnnouncementUsecase) AdminDelete(ctx context.Context, actor *model.Actor, in sitein.AdminAnnouncementDeleteInp) error {
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := service.SiteAnnouncementDomain().DeleteReadRecordsByAnnouncementId(ctx, in.Id); err != nil {
			return err
		}
		return service.SiteAnnouncementDomain().AdminDelete(ctx, in.Id)
	})
	if err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionDelete,
		TargetType: consts.SiteAuditTargetTypeSiteAnnouncement,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
	})
	return nil
}

func (s *sSiteAnnouncementUsecase) buildAnnouncementItems(list []entity.SiteAnnouncement, readMap map[uint64]bool) []*siteout.AnnouncementItem {
	items := make([]*siteout.AnnouncementItem, 0, len(list))
	for _, item := range list {
		items = append(items, &siteout.AnnouncementItem{
			Id:          item.Id,
			Title:       item.Title,
			Content:     item.Content,
			Status:      item.Status,
			IsRead:      readMap != nil && readMap[item.Id],
			CreatedBy:   item.CreatedBy,
			UpdatedBy:   item.UpdatedBy,
			PublishedAt: item.PublishedAt,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}
	return items
}

func (s *sSiteAnnouncementUsecase) announcementIds(list []entity.SiteAnnouncement) []uint64 {
	ids := make([]uint64, 0, len(list))
	for _, item := range list {
		if item.Id > 0 {
			ids = append(ids, item.Id)
		}
	}
	return ids
}

func (s *sSiteAnnouncementUsecase) actorId(actor *model.Actor) uint64 {
	if actor == nil {
		return 0
	}
	return actor.Id
}

func (s *sSiteAnnouncementUsecase) validateStatusPtr(status *int) error {
	if status == nil {
		return nil
	}
	return s.validateStatus(*status)
}

func (s *sSiteAnnouncementUsecase) validateStatus(status int) error {
	switch status {
	case consts.SiteAnnouncementStatusDraft,
		consts.SiteAnnouncementStatusPublished,
		consts.SiteAnnouncementStatusArchived:
		return nil
	default:
		return gerror.New("invalid announcement status")
	}
}
