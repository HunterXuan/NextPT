package accounting

import (
	"context"
	"fmt"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sAccountingTrafficDomain struct{}

type trafficStatChange struct {
	userId        uint64
	diffUp        int64
	diffDn        int64
	seedTimeDiff  int
	leechTimeDiff int
	eventTime     *gtime.Time
}

func init() {
	service.RegisterAccountingTrafficDomain(NewAccountingTrafficDomain())
}

func NewAccountingTrafficDomain() *sAccountingTrafficDomain {
	return &sAccountingTrafficDomain{}
}

func (s *sAccountingTrafficDomain) GetUserStat(ctx context.Context, userId uint64) (*entity.IamUserStat, error) {
	var stat *entity.IamUserStat
	err := dao.IamUserStat.Ctx(ctx).Where("user_id", userId).Scan(&stat)
	return stat, err
}

func (s *sAccountingTrafficDomain) RecordTraffic(ctx context.Context, userId uint64, diffUp, diffDn int64, isSeeder bool, timeDiff int, eventTime *gtime.Time) error {
	if diffUp <= 0 && diffDn <= 0 && timeDiff <= 0 {
		return nil
	}

	statColumns := dao.IamUserStat.Columns()
	statModel := dao.IamUserStat.Ctx(ctx)
	statUpdate := g.Map{
		statColumns.Uploaded:   s.incrementByValue(statModel, statColumns.Uploaded, diffUp),
		statColumns.Downloaded: s.incrementByValue(statModel, statColumns.Downloaded, diffDn),
	}
	if isSeeder {
		statUpdate[statColumns.SeedTime] = s.incrementByValue(statModel, statColumns.SeedTime, int64(timeDiff))
	} else {
		statUpdate[statColumns.LeechTime] = s.incrementByValue(statModel, statColumns.LeechTime, int64(timeDiff))
	}
	res, err := statModel.Where(statColumns.UserId, userId).Update(statUpdate)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return gerror.Newf("user stat record missing for user_id: %d", userId)
	}

	// 2. Update daily and monthly period stat.
	if eventTime == nil {
		eventTime = gtime.Now()
	}

	var seedTimeDiff, leechTimeDiff int
	if isSeeder {
		seedTimeDiff = timeDiff
	} else {
		leechTimeDiff = timeDiff
	}

	change := trafficStatChange{
		userId:        userId,
		diffUp:        diffUp,
		diffDn:        diffDn,
		seedTimeDiff:  seedTimeDiff,
		leechTimeDiff: leechTimeDiff,
		eventTime:     eventTime,
	}

	if err = s.upsertPeriodStat(ctx, consts.AccountingStatPeriodDaily, eventTime.Format("Y-m-d"), change); err != nil {
		return err
	}

	return s.upsertPeriodStat(ctx, consts.AccountingStatPeriodMonthly, eventTime.Format("Y-m"), change)
}

func (s *sAccountingTrafficDomain) upsertPeriodStat(ctx context.Context, periodType int, periodKey string, change trafficStatChange) error {
	columns := dao.IamUserPeriodStat.Columns()
	m := dao.IamUserPeriodStat.Ctx(ctx)
	_, err := m.
		Data(g.Map{
			columns.UserId:     change.userId,
			columns.PeriodType: periodType,
			columns.PeriodKey:  periodKey,
			columns.Uploaded:   change.diffUp,
			columns.Downloaded: change.diffDn,
			columns.SeedTime:   change.seedTimeDiff,
			columns.LeechTime:  change.leechTimeDiff,
			columns.Bonus:      0,
			columns.CreatedAt:  change.eventTime,
		}).
		OnDuplicate(g.Map{
			columns.Uploaded:   s.incrementByInsertedValue(m, columns.Uploaded),
			columns.Downloaded: s.incrementByInsertedValue(m, columns.Downloaded),
			columns.SeedTime:   s.incrementByInsertedValue(m, columns.SeedTime),
			columns.LeechTime:  s.incrementByInsertedValue(m, columns.LeechTime),
		}).
		Save()
	return err
}

func (s *sAccountingTrafficDomain) incrementByInsertedValue(m *gdb.Model, column string) gdb.Raw {
	column = m.QuoteWord(column)
	return gdb.Raw(fmt.Sprintf("%s + VALUES(%s)", column, column))
}

func (s *sAccountingTrafficDomain) incrementByValue(m *gdb.Model, column string, delta int64) gdb.Raw {
	column = m.QuoteWord(column)
	return gdb.Raw(fmt.Sprintf("%s + %s", column, gconv.String(delta)))
}

func (s *sAccountingTrafficDomain) QueryPeriodStats(ctx context.Context, userId uint64, periodType int, startDate, endDate *gtime.Time) ([]entity.IamUserPeriodStat, error) {
	columns := dao.IamUserPeriodStat.Columns()
	q := dao.IamUserPeriodStat.Ctx(ctx).
		Where(columns.UserId, userId).
		Where(columns.PeriodType, periodType)
	if startDate != nil {
		startKey, err := s.formatPeriodKey(periodType, startDate)
		if err != nil {
			return nil, err
		}
		q = q.WhereGTE(columns.PeriodKey, startKey)
	}
	if endDate != nil {
		endKey, err := s.formatPeriodKey(periodType, endDate)
		if err != nil {
			return nil, err
		}
		q = q.WhereLTE(columns.PeriodKey, endKey)
	}
	var list []entity.IamUserPeriodStat
	err := q.OrderDesc(columns.PeriodKey).Limit(100).Scan(&list)
	return list, err
}

func (s *sAccountingTrafficDomain) formatPeriodKey(periodType int, t *gtime.Time) (string, error) {
	switch periodType {
	case consts.AccountingStatPeriodDaily:
		return t.Format("Y-m-d"), nil
	case consts.AccountingStatPeriodMonthly:
		return t.Format("Y-m"), nil
	default:
		return "", gerror.Newf("unsupported traffic stat period type: %d", periodType)
	}
}
