package admin

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/modin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/modout"
	"server/internal/service"
)

type sAdminModCheaterUsecase struct{}

func init() {
	service.RegisterAdminModCheaterUsecase(NewAdminModCheaterUsecase())
}

func NewAdminModCheaterUsecase() *sAdminModCheaterUsecase {
	return &sAdminModCheaterUsecase{}
}

func (s *sAdminModCheaterUsecase) List(ctx context.Context, actor *model.Actor, in adminin.ModCheaterListInp) (*modout.ListCheaterLogsOut, error) {
	out, err := service.ModCheaterUsecase().List(ctx, actor, modin.ListCheaterLogsInp{
		IsDealt: &in.Status,
		Page:    in.Page,
		Size:    in.Size,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *sAdminModCheaterUsecase) Resolve(ctx context.Context, actor *model.Actor, in adminin.ModCheaterResolveInp) error {
	if err := service.ModCheaterUsecase().Resolve(ctx, actor, modin.ResolveCheaterLogInp{
		Id:      in.Id,
		Comment: in.Comment,
	}); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeModCheaterLog,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"operation": consts.SiteAuditOperationResolve,
		},
	})
	return nil
}
