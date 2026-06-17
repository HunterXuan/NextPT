package economy

import (
	"context"

	v1 "server/api/economy/v1"
	"server/internal/library/contexts"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

func (c *ControllerV1) HourlyBonus(ctx context.Context, req *v1.HourlyBonusReq) (res *v1.HourlyBonusRes, err error) {
	actor := contexts.GetActor(ctx)
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	out, err := service.EconomyBonusUsecase().GetMyHourlyBonus(ctx, actor, req.HourlyBonusInp)
	if err != nil {
		return nil, err
	}

	return &v1.HourlyBonusRes{HourlyBonusOut: *out}, nil
}
