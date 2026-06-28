package mod

import (
	"context"

	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"
)

type sModReportDomain struct{}

func init() {
	service.RegisterModReportDomain(NewModReportDomain())
}

func NewModReportDomain() *sModReportDomain {
	return &sModReportDomain{}
}

func (s *sModReportDomain) GetPendingCount(ctx context.Context, reporterId uint64, targetType string, targetId uint64) (int, error) {
	return dao.ModReport.Ctx(ctx).Where(dao.ModReport.Columns().ReporterId, reporterId).
		Where(dao.ModReport.Columns().TargetType, targetType).
		Where(dao.ModReport.Columns().TargetId, targetId).
		Where(dao.ModReport.Columns().Status, 0).
		Count()
}

func (s *sModReportDomain) Create(ctx context.Context, report entity.ModReport) error {
	_, err := dao.ModReport.Ctx(ctx).Data(report).Insert()
	return err
}

func (s *sModReportDomain) GetById(ctx context.Context, id uint64) (*entity.ModReport, error) {
	var report entity.ModReport
	err := dao.ModReport.Ctx(ctx).WherePri(id).Scan(&report)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (s *sModReportDomain) Update(ctx context.Context, id uint64, data interface{}) error {
	_, err := dao.ModReport.Ctx(ctx).WherePri(id).Data(data).Update()
	return err
}

func (s *sModReportDomain) QueryReports(ctx context.Context, status int, targetType string, page, size int) ([]entity.ModReport, int, error) {
	m := dao.ModReport.Ctx(ctx)
	if status >= 0 {
		m = m.Where(dao.ModReport.Columns().Status, status)
	}
	if targetType != "" {
		m = m.Where(dao.ModReport.Columns().TargetType, targetType)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var records []entity.ModReport
	err = m.Page(page, size).OrderDesc(dao.ModReport.Columns().CreatedAt).Scan(&records)
	return records, total, err
}

func (s *sModReportDomain) DeleteReportsByTarget(ctx context.Context, targetType string, targetId uint64) error {
	_, err := dao.ModReport.Ctx(ctx).Where(dao.ModReport.Columns().TargetType, targetType).
		Where(dao.ModReport.Columns().TargetId, targetId).Delete()
	return err
}
