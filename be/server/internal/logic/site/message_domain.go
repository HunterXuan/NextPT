package site

import (
	"context"

	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/sitein"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

const siteMessageListMaxSize = 100

type sSiteMessageDomain struct{}

func init() {
	service.RegisterSiteMessageDomain(NewSiteMessageDomain())
}

func NewSiteMessageDomain() *sSiteMessageDomain {
	return &sSiteMessageDomain{}
}

func (s *sSiteMessageDomain) ListByReceiver(ctx context.Context, receiverId uint64, in sitein.MessageListInp) ([]entity.SiteMessage, int, error) {
	page, size := s.normalizeListPage(in.Page, in.Size)
	columns := dao.SiteMessage.Columns()
	m := dao.SiteMessage.Ctx(ctx).
		Where(columns.ReceiverId, receiverId)
	if in.IsRead != nil {
		m = m.Where(columns.IsRead, *in.IsRead)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []entity.SiteMessage
	err = m.Page(page, size).OrderDesc(columns.Id).Scan(&list)
	return list, total, err
}

func (s *sSiteMessageDomain) MarkRead(ctx context.Context, receiverId uint64, id uint64) error {
	if receiverId == 0 || id == 0 {
		return nil
	}
	columns := dao.SiteMessage.Columns()
	_, err := dao.SiteMessage.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.ReceiverId, receiverId).
		Where(columns.IsRead, false).
		Data(g.Map{
			columns.IsRead: true,
			columns.ReadAt: gtime.Now(),
		}).
		Update()
	return err
}

func (s *sSiteMessageDomain) MarkAllRead(ctx context.Context, receiverId uint64) error {
	if receiverId == 0 {
		return nil
	}
	columns := dao.SiteMessage.Columns()
	_, err := dao.SiteMessage.Ctx(ctx).
		Where(columns.ReceiverId, receiverId).
		Where(columns.IsRead, false).
		Data(g.Map{
			columns.IsRead: true,
			columns.ReadAt: gtime.Now(),
		}).
		Update()
	return err
}

func (s *sSiteMessageDomain) Create(ctx context.Context, in sitein.MessageCreateInp) (uint64, error) {
	id, err := dao.SiteMessage.Ctx(ctx).Data(s.messageDo(in)).InsertAndGetId()
	return uint64(id), err
}

func (s *sSiteMessageDomain) BatchCreate(ctx context.Context, items []sitein.MessageCreateInp) error {
	if len(items) == 0 {
		return nil
	}
	for _, item := range items {
		if _, err := s.Create(ctx, item); err != nil {
			return err
		}
	}
	return nil
}

func (s *sSiteMessageDomain) AdminList(ctx context.Context, in sitein.AdminMessageListInp) ([]entity.SiteMessage, int, error) {
	page, size := s.normalizeListPage(in.Page, in.Size)
	columns := dao.SiteMessage.Columns()
	m := dao.SiteMessage.Ctx(ctx)
	if in.ReceiverId > 0 {
		m = m.Where(columns.ReceiverId, in.ReceiverId)
	}
	if in.IsRead != nil {
		m = m.Where(columns.IsRead, *in.IsRead)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []entity.SiteMessage
	err = m.Page(page, size).OrderDesc(columns.Id).Scan(&list)
	return list, total, err
}

func (s *sSiteMessageDomain) messageDo(in sitein.MessageCreateInp) do.SiteMessage {
	return do.SiteMessage{
		SenderId:   in.SenderId,
		ReceiverId: in.ReceiverId,
		Title:      in.Title,
		Content:    in.Content,
		TargetType: in.TargetType,
		TargetId:   in.TargetId,
		IsRead:     false,
		CreatedAt:  gtime.Now(),
	}
}

func (s *sSiteMessageDomain) normalizeListPage(page int, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > siteMessageListMaxSize {
		size = siteMessageListMaxSize
	}
	return page, size
}
