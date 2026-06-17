package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/modin"
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
	return service.ModCheaterUsecase().Resolve(ctx, actor, modin.ResolveCheaterLogInp{
		Id:      in.Id,
		Comment: in.Comment,
	})
}
