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

type sAdminModUserUsecase struct{}

func init() {
	service.RegisterAdminModUserUsecase(NewAdminModUserUsecase())
}

func NewAdminModUserUsecase() *sAdminModUserUsecase {
	return &sAdminModUserUsecase{}
}

func (s *sAdminModUserUsecase) List(ctx context.Context, actor *model.Actor, in adminin.ModUserListInp) (*modout.ListUserOut, error) {
	return service.ModUserUsecase().List(ctx, actor, modin.ListUserInp{
		UserId: in.Id,
	})
}

func (s *sAdminModUserUsecase) Apply(ctx context.Context, actor *model.Actor, in adminin.ModUserApplyInp) error {
	if err := service.ModUserUsecase().Apply(ctx, actor, modin.ApplyModInp{
		UserId:          in.Id,
		ModType:         int(in.Type),
		Reason:          in.Reason,
		DurationSeconds: in.Duration,
	}); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionApply,
		TargetType: consts.SiteAuditTargetTypeIamUser,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"modType":         in.Type,
			"durationSeconds": in.Duration,
			"reason":          in.Reason,
		},
	})
	return nil
}

func (s *sAdminModUserUsecase) Remove(ctx context.Context, actor *model.Actor, in adminin.ModUserRemoveInp) error {
	if err := service.ModUserUsecase().Remove(ctx, actor, modin.RemoveModInp{
		Id:     in.ModId,
		UserId: in.Id,
	}); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionRemove,
		TargetType: consts.SiteAuditTargetTypeModUserLog,
		TargetId:   in.ModId,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"userId": in.Id,
		},
	})
	return nil
}
