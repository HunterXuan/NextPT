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
