package site

import (
	"context"

	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/sitein"
	"server/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

type sSiteAuditDomain struct{}

const siteAuditListMaxSize = 100

func init() {
	service.RegisterSiteAuditDomain(NewSiteAuditDomain())
}

func NewSiteAuditDomain() *sSiteAuditDomain {
	return &sSiteAuditDomain{}
}

func (s *sSiteAuditDomain) AdminListAudits(ctx context.Context, in sitein.AuditListInp) ([]entity.SiteAudit, int, error) {
	page, size := s.normalizeListPage(in.Page, in.Size)
	columns := dao.SiteAudit.Columns()
	m := dao.SiteAudit.Ctx(ctx)
	if in.Level != nil {
		m = m.Where(columns.Level, *in.Level)
	}
	if in.Action != "" {
		m = m.Where(columns.Action, in.Action)
	}
	if in.TargetType != "" {
		m = m.Where(columns.TargetType, in.TargetType)
	}
	if in.UserId > 0 {
		m = m.Where(columns.UserId, in.UserId)
	}
	if in.StartAt != nil {
		m = m.WhereGTE(columns.CreatedAt, in.StartAt)
	}
	if in.EndAt != nil {
		m = m.WhereLTE(columns.CreatedAt, in.EndAt)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []entity.SiteAudit
	err = m.Page(page, size).OrderDesc(columns.Id).Scan(&list)
	return list, total, err
}

func (s *sSiteAuditDomain) normalizeListPage(page int, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > siteAuditListMaxSize {
		size = siteAuditListMaxSize
	}
	return page, size
}

func (s *sSiteAuditDomain) Create(ctx context.Context, in sitein.AuditCreateInp) error {
	if in.Action == "" || in.TargetType == "" {
		return nil
	}
	if in.CreatedAt == nil {
		in.CreatedAt = gtime.Now()
	}
	_, err := dao.SiteAudit.Ctx(ctx).Data(do.SiteAudit{
		UserId:     in.UserId,
		Action:     in.Action,
		TargetType: in.TargetType,
		TargetId:   in.TargetId,
		Detail:     in.Detail,
		Ip:         in.Ip,
		Level:      in.Level,
		CreatedAt:  in.CreatedAt,
	}).Insert()
	return err
}
