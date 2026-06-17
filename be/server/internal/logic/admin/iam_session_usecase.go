package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/util/gconv"
)

type sAdminIamSessionUsecase struct{}

func NewAdminIamSessionUsecase() *sAdminIamSessionUsecase {
	return &sAdminIamSessionUsecase{}
}

func init() {
	service.RegisterAdminIamSessionUsecase(NewAdminIamSessionUsecase())
}

func (s *sAdminIamSessionUsecase) DeleteByUser(ctx context.Context, actor *model.Actor, in adminin.IamSessionDeleteInp) error {
	err := service.IamSessionDomain().RemoveToken(ctx, gconv.String(in.UserId))
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "admin.session.delete_failed"))
	}
	return nil
}
