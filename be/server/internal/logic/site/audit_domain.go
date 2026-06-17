package site

import (
	"context"

	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"
)

type sSiteAuditDomain struct{}

func init() {
	service.RegisterSiteAuditDomain(NewSiteAuditDomain())
}

func NewSiteAuditDomain() *sSiteAuditDomain {
	return &sSiteAuditDomain{}
}

func (s *sSiteAuditDomain) AdminListAudits(ctx context.Context, page, size int) ([]entity.SiteAudit, int, error) {
	m := dao.SiteAudit.Ctx(ctx)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []entity.SiteAudit
	err = m.Page(page, size).OrderDesc(dao.SiteAudit.Columns().Id).Scan(&list)
	return list, total, err
}
