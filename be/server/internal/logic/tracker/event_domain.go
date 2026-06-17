package tracker

import (
	"context"

	"server/internal/dao"
	"server/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

type sTrackerEventDomain struct{}

func init() {
	service.RegisterTrackerEventDomain(NewTrackerEventDomain())
}

func NewTrackerEventDomain() *sTrackerEventDomain {
	return &sTrackerEventDomain{}
}

// DeleteExpiredIdempotency 删除过期的幂等记录
func (s *sTrackerEventDomain) DeleteExpiredIdempotency(ctx context.Context, beforeTime *gtime.Time) (int64, error) {
	result, err := dao.TrackerEventIdempotency.Ctx(ctx).Where("created_at < ?", beforeTime).Delete()
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
}
