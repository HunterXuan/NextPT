// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package economy

import (
	"context"

	"server/api/economy/v1"
)

type IEconomyV1 interface {
	BonusLogs(ctx context.Context, req *v1.BonusLogsReq) (res *v1.BonusLogsRes, err error)
	HourlyBonus(ctx context.Context, req *v1.HourlyBonusReq) (res *v1.HourlyBonusRes, err error)
}
