package site

import (
	"context"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/sitein"
	"server/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

const siteAnnouncementListMaxSize = 100

type sSiteAnnouncementDomain struct{}

func init() {
	service.RegisterSiteAnnouncementDomain(NewSiteAnnouncementDomain())
}

func NewSiteAnnouncementDomain() *sSiteAnnouncementDomain {
	return &sSiteAnnouncementDomain{}
}

func (s *sSiteAnnouncementDomain) ListPublished(ctx context.Context, userId uint64, in sitein.AnnouncementListInp) ([]entity.SiteAnnouncement, int, error) {
	page, size := s.normalizeListPage(in.Page, in.Size)
	columns := dao.SiteAnnouncement.Columns()
	m := dao.SiteAnnouncement.Ctx(ctx).
		Where(columns.Status, consts.SiteAnnouncementStatusPublished)
	if in.IsRead != nil {
		readIds, err := s.QueryReadAnnouncementIdsByUser(ctx, userId)
		if err != nil {
			return nil, 0, err
		}
		if *in.IsRead {
			if len(readIds) == 0 {
				return []entity.SiteAnnouncement{}, 0, nil
			}
			m = m.WhereIn(columns.Id, readIds)
		} else if len(readIds) > 0 {
			m = m.WhereNotIn(columns.Id, readIds)
		}
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []entity.SiteAnnouncement
	err = m.Page(page, size).Order(columns.PublishedAt + " DESC, " + columns.Id + " DESC").Scan(&list)
	return list, total, err
}

func (s *sSiteAnnouncementDomain) QueryReadAnnouncementIdsByUser(ctx context.Context, userId uint64) ([]uint64, error) {
	if userId == 0 {
		return []uint64{}, nil
	}
	columns := dao.SiteAnnouncementRead.Columns()
	var reads []entity.SiteAnnouncementRead
	if err := dao.SiteAnnouncementRead.Ctx(ctx).
		Where(columns.UserId, userId).
		Scan(&reads); err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(reads))
	for _, read := range reads {
		if read.AnnouncementId > 0 {
			ids = append(ids, read.AnnouncementId)
		}
	}
	return ids, nil
}

func (s *sSiteAnnouncementDomain) QueryReadIds(ctx context.Context, userId uint64, announcementIds []uint64) (map[uint64]bool, error) {
	readMap := make(map[uint64]bool, len(announcementIds))
	if userId == 0 || len(announcementIds) == 0 {
		return readMap, nil
	}
	columns := dao.SiteAnnouncementRead.Columns()
	var reads []entity.SiteAnnouncementRead
	err := dao.SiteAnnouncementRead.Ctx(ctx).
		Where(columns.UserId, userId).
		WhereIn(columns.AnnouncementId, announcementIds).
		Scan(&reads)
	if err != nil {
		return nil, err
	}
	for _, read := range reads {
		readMap[read.AnnouncementId] = true
	}
	return readMap, nil
}

func (s *sSiteAnnouncementDomain) MarkRead(ctx context.Context, userId uint64, announcementId uint64) error {
	if userId == 0 || announcementId == 0 {
		return nil
	}
	_, err := dao.SiteAnnouncementRead.Ctx(ctx).Data(do.SiteAnnouncementRead{
		AnnouncementId: announcementId,
		UserId:         userId,
		ReadAt:         gtime.Now(),
	}).Save()
	return err
}

func (s *sSiteAnnouncementDomain) AdminList(ctx context.Context, in sitein.AdminAnnouncementListInp) ([]entity.SiteAnnouncement, int, error) {
	page, size := s.normalizeListPage(in.Page, in.Size)
	columns := dao.SiteAnnouncement.Columns()
	m := dao.SiteAnnouncement.Ctx(ctx)
	if in.Status != nil {
		m = m.Where(columns.Status, *in.Status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []entity.SiteAnnouncement
	err = m.Page(page, size).OrderDesc(columns.Id).Scan(&list)
	return list, total, err
}

func (s *sSiteAnnouncementDomain) AdminCreate(ctx context.Context, in sitein.AdminAnnouncementCreateInp, actorId uint64) (uint64, error) {
	publishedAt := s.normalizePublishedAt(in.Status, in.PublishedAt)
	id, err := dao.SiteAnnouncement.Ctx(ctx).Data(do.SiteAnnouncement{
		Title:       in.Title,
		Content:     in.Content,
		Status:      in.Status,
		CreatedBy:   actorId,
		UpdatedBy:   actorId,
		PublishedAt: publishedAt,
		CreatedAt:   gtime.Now(),
		UpdatedAt:   gtime.Now(),
	}).InsertAndGetId()
	return uint64(id), err
}

func (s *sSiteAnnouncementDomain) AdminUpdate(ctx context.Context, in sitein.AdminAnnouncementUpdateInp, actorId uint64) error {
	columns := dao.SiteAnnouncement.Columns()
	_, err := dao.SiteAnnouncement.Ctx(ctx).
		Where(columns.Id, in.Id).
		Data(do.SiteAnnouncement{
			Title:       in.Title,
			Content:     in.Content,
			Status:      in.Status,
			UpdatedBy:   actorId,
			PublishedAt: s.normalizePublishedAt(in.Status, in.PublishedAt),
			UpdatedAt:   gtime.Now(),
		}).
		Update()
	return err
}

func (s *sSiteAnnouncementDomain) AdminDelete(ctx context.Context, id uint64) error {
	announcementColumns := dao.SiteAnnouncement.Columns()
	_, err := dao.SiteAnnouncement.Ctx(ctx).
		Where(announcementColumns.Id, id).
		Delete()
	return err
}

func (s *sSiteAnnouncementDomain) DeleteReadRecordsByAnnouncementId(ctx context.Context, announcementId uint64) error {
	readColumns := dao.SiteAnnouncementRead.Columns()
	_, err := dao.SiteAnnouncementRead.Ctx(ctx).
		Where(readColumns.AnnouncementId, announcementId).
		Delete()
	return err
}

func (s *sSiteAnnouncementDomain) normalizeListPage(page int, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > siteAnnouncementListMaxSize {
		size = siteAnnouncementListMaxSize
	}
	return page, size
}

func (s *sSiteAnnouncementDomain) normalizePublishedAt(status int, publishedAt *gtime.Time) *gtime.Time {
	if status == consts.SiteAnnouncementStatusPublished && publishedAt == nil {
		return gtime.Now()
	}
	return publishedAt
}
