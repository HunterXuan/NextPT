package mod

import (
	"context"

	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/modin"
	"server/internal/model/out/modout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

type sModReportUsecase struct{}

func init() {
	service.RegisterModReportUsecase(NewModReportUsecase())
}

func NewModReportUsecase() *sModReportUsecase {
	return &sModReportUsecase{}
}

func (s *sModReportUsecase) Create(ctx context.Context, actor *model.Actor, in modin.CreateReportInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	count, err := service.ModReportDomain().GetPendingCount(ctx, actor.Id, in.TargetType, in.TargetId)
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New(gi18n.T(ctx, "mod.report.duplicate"))
	}

	return service.ModReportDomain().Create(ctx, entity.ModReport{
		ReporterId: actor.Id,
		TargetType: in.TargetType,
		TargetId:   in.TargetId,
		Reason:     in.Reason,
		Status:     0,
	})
}

func (s *sModReportUsecase) List(ctx context.Context, actor *model.Actor, in modin.ListReportsInp) (*modout.ListReportsOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	records, total, err := service.ModReportDomain().QueryReports(ctx, in.Status, in.TargetType, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	var list []modout.ReportItem
	for _, r := range records {
		list = append(list, modout.ReportItem{
			Id:           r.Id,
			ReporterId:   r.ReporterId,
			TargetType:   r.TargetType,
			TargetId:     r.TargetId,
			Reason:       r.Reason,
			Status:       r.Status,
			DealtBy:      r.DealtBy,
			DealtComment: r.DealtComment,
			DealtAt:      r.DealtAt,
			CreatedAt:    r.CreatedAt,
		})
	}

	return &modout.ListReportsOut{
		Page:  in.Page,
		Size:  in.Size,
		Total: total,
		List:  list,
	}, nil
}

func (s *sModReportUsecase) Resolve(ctx context.Context, actor *model.Actor, in modin.ResolveReportInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	report, err := service.ModReportDomain().GetById(ctx, in.Id)
	if err != nil || report == nil {
		return gerror.New(gi18n.T(ctx, "moderation.report.not_found"))
	}

	if report.Status != 0 {
		return gerror.New(gi18n.T(ctx, "moderation.report.already_resolved"))
	}

	return service.ModReportDomain().Update(ctx, in.Id, entity.ModReport{
		Status:       in.Status,
		DealtBy:      actor.Id,
		DealtComment: in.Comment,
		DealtAt:      gtime.Now(),
	})
}
