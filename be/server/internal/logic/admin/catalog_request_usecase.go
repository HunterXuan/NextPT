package admin

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/in/catalogin"
	"server/internal/model/in/sitein"
	"server/internal/service"
)

type sAdminCatalogRequestUsecase struct{}

func init() {
	service.RegisterAdminCatalogRequestUsecase(NewAdminCatalogRequestUsecase())
}

func NewAdminCatalogRequestUsecase() *sAdminCatalogRequestUsecase {
	return &sAdminCatalogRequestUsecase{}
}

func (s *sAdminCatalogRequestUsecase) Complete(ctx context.Context, actor *model.Actor, in catalogin.RequestCompleteInp) error {
	if err := service.CatalogRequestUsecase().CompleteByAdmin(ctx, actor, in); err != nil {
		return err
	}
	s.recordUpdateAudit(ctx, actor, in.Id, consts.SiteAuditOperationComplete)
	return nil
}

func (s *sAdminCatalogRequestUsecase) Cancel(ctx context.Context, actor *model.Actor, in catalogin.RequestCancelInp) error {
	if err := service.CatalogRequestUsecase().CancelByAdmin(ctx, actor, in); err != nil {
		return err
	}
	s.recordUpdateAudit(ctx, actor, in.Id, consts.SiteAuditOperationCancel)
	return nil
}

func (s *sAdminCatalogRequestUsecase) recordUpdateAudit(ctx context.Context, actor *model.Actor, requestId uint64, operation string) {
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeCatalogRequest,
		TargetId:   requestId,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"operation": operation,
		},
	})
}
