package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/modin"
	"server/internal/service"
)

type sAdminModUserUsecase struct{}

func init() {
	service.RegisterAdminModUserUsecase(NewAdminModUserUsecase())
}

func NewAdminModUserUsecase() *sAdminModUserUsecase {
	return &sAdminModUserUsecase{}
}

func (s *sAdminModUserUsecase) Apply(ctx context.Context, actor *model.Actor, in adminin.ModUserApplyInp) error {
	return service.ModUserUsecase().Apply(ctx, actor, modin.ApplyModInp{
		UserId:          in.Id,
		ModType:         int(in.Type),
		Reason:          in.Reason,
		DurationSeconds: in.Duration,
	})
}

func (s *sAdminModUserUsecase) Remove(ctx context.Context, actor *model.Actor, in adminin.ModUserRemoveInp) error {
	return service.ModUserUsecase().Remove(ctx, actor, modin.RemoveModInp{
		Id:     in.ModId,
		UserId: in.Id,
	})
}
