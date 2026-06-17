package economy

import (
	"context"

	"server/internal/library/contexts"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"

	v1 "server/api/economy/v1"
)

func (c *ControllerV1) BonusLogs(ctx context.Context, req *v1.BonusLogsReq) (res *v1.BonusLogsRes, err error) {
	actor := contexts.GetActor(ctx)
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	out, err := service.EconomyBonusUsecase().ListMyBonusLogs(ctx, actor, req.BonusLogsInp)
	if err != nil {
		return nil, err
	}

	return &v1.BonusLogsRes{BonusLogsOut: *out}, nil
}
