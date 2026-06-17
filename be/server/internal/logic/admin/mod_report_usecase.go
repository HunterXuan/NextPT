package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/modin"
	"server/internal/model/out/modout"
	"server/internal/service"
)

type sAdminModReportUsecase struct{}

func init() {
	service.RegisterAdminModReportUsecase(NewAdminModReportUsecase())
}

func NewAdminModReportUsecase() *sAdminModReportUsecase {
	return &sAdminModReportUsecase{}
}

func (s *sAdminModReportUsecase) List(ctx context.Context, actor *model.Actor, in adminin.ModReportListInp) (*modout.ListReportsOut, error) {
	out, err := service.ModReportUsecase().List(ctx, actor, modin.ListReportsInp{
		Page:       in.Page,
		Size:       in.Size,
		Status:     in.Status,
		TargetType: in.TargetType,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *sAdminModReportUsecase) Resolve(ctx context.Context, actor *model.Actor, in adminin.ModReportResolveInp) error {
	return service.ModReportUsecase().Resolve(ctx, actor, modin.ResolveReportInp{
		Id:      in.Id,
		Status:  in.Status,
		Comment: in.Comment,
	})
}
