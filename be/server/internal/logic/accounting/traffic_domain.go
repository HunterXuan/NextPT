package accounting

import (
	"context"
	"fmt"

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

type trafficStatColumns struct {
	userId     string
	period     string
	uploaded   string
	downloaded string
	seedTime   string
	leechTime  string
	bonus      string
	createdAt  string
}

type trafficStatTarget struct {
	model   *gdb.Model
	period  string
	columns trafficStatColumns
}

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

	// 1. Update user_stat
	statUpdate := g.Map{
		dao.IamUserStat.Columns().Uploaded:   gdb.Raw("uploaded + " + gconv.String(diffUp)),
		dao.IamUserStat.Columns().Downloaded: gdb.Raw("downloaded + " + gconv.String(diffDn)),
	}
	if isSeeder {
		statUpdate[dao.IamUserStat.Columns().SeedTime] = gdb.Raw("seed_time + " + gconv.String(timeDiff))
	} else {
		statUpdate[dao.IamUserStat.Columns().LeechTime] = gdb.Raw("leech_time + " + gconv.String(timeDiff))
	}
	res, err := dao.IamUserStat.Ctx(ctx).Where("user_id", userId).Update(statUpdate)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return gerror.Newf("user stat record missing for user_id: %d", userId)
	}

	// 2. Update daily and monthly stat.
	if eventTime == nil {
		eventTime = gtime.Now()
	}
	dateStr := eventTime.Format("Y-m-d")
	monthStr := eventTime.Format("Y-m")

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

	if err = s.upsertTrafficStat(s.dailyTrafficStatTarget(ctx, dateStr), change); err != nil {
		return err
	}

	return s.upsertTrafficStat(s.monthlyTrafficStatTarget(ctx, monthStr), change)
}

func (s *sAccountingTrafficDomain) dailyTrafficStatTarget(ctx context.Context, period string) trafficStatTarget {
	columns := dao.IamUserDailyStat.Columns()
	return trafficStatTarget{
		model:  dao.IamUserDailyStat.Ctx(ctx),
		period: period,
		columns: trafficStatColumns{
			userId:     columns.UserId,
			period:     columns.Date,
			uploaded:   columns.Uploaded,
			downloaded: columns.Downloaded,
			seedTime:   columns.SeedTime,
			leechTime:  columns.LeechTime,
			bonus:      columns.Bonus,
			createdAt:  columns.CreatedAt,
		},
	}
}

func (s *sAccountingTrafficDomain) monthlyTrafficStatTarget(ctx context.Context, period string) trafficStatTarget {
	columns := dao.IamUserMonthlyStat.Columns()
	return trafficStatTarget{
		model:  dao.IamUserMonthlyStat.Ctx(ctx),
		period: period,
		columns: trafficStatColumns{
			userId:     columns.UserId,
			period:     columns.YearMonth,
			uploaded:   columns.Uploaded,
			downloaded: columns.Downloaded,
			seedTime:   columns.SeedTime,
			leechTime:  columns.LeechTime,
			bonus:      columns.Bonus,
			createdAt:  columns.CreatedAt,
		},
	}
}

func (s *sAccountingTrafficDomain) upsertTrafficStat(target trafficStatTarget, change trafficStatChange) error {
	columns := target.columns
	_, err := target.model.
		Data(g.Map{
			columns.userId:     change.userId,
			columns.period:     target.period,
			columns.uploaded:   change.diffUp,
			columns.downloaded: change.diffDn,
			columns.seedTime:   change.seedTimeDiff,
			columns.leechTime:  change.leechTimeDiff,
			columns.bonus:      0,
			columns.createdAt:  change.eventTime,
		}).
		OnDuplicate(g.Map{
			columns.uploaded:   s.incrementByInsertedValue(target.model, columns.uploaded),
			columns.downloaded: s.incrementByInsertedValue(target.model, columns.downloaded),
			columns.seedTime:   s.incrementByInsertedValue(target.model, columns.seedTime),
			columns.leechTime:  s.incrementByInsertedValue(target.model, columns.leechTime),
		}).
		Insert()
	return err
}

func (s *sAccountingTrafficDomain) incrementByInsertedValue(m *gdb.Model, column string) gdb.Raw {
	column = m.QuoteWord(column)
	return gdb.Raw(fmt.Sprintf("%s + VALUES(%s)", column, column))
}

func (s *sAccountingTrafficDomain) QueryDailyStats(ctx context.Context, userId uint64, startDate, endDate *gtime.Time) ([]entity.IamUserDailyStat, error) {
	columns := dao.IamUserDailyStat.Columns()
	q := dao.IamUserDailyStat.Ctx(ctx).Where(columns.UserId, userId)
	if startDate != nil {
		q = q.WhereGTE(columns.Date, startDate.Format("Y-m-d"))
	}
	if endDate != nil {
		q = q.WhereLTE(columns.Date, endDate.Format("Y-m-d"))
	}
	var list []entity.IamUserDailyStat
	err := q.OrderDesc(columns.Date).Limit(100).Scan(&list)
	return list, err
}

func (s *sAccountingTrafficDomain) QueryMonthlyStats(ctx context.Context, userId uint64, startDate, endDate *gtime.Time) ([]entity.IamUserMonthlyStat, error) {
	columns := dao.IamUserMonthlyStat.Columns()
	q := dao.IamUserMonthlyStat.Ctx(ctx).Where(columns.UserId, userId)
	if startDate != nil {
		q = q.WhereGTE(columns.YearMonth, startDate.Format("Y-m"))
	}
	if endDate != nil {
		q = q.WhereLTE(columns.YearMonth, endDate.Format("Y-m"))
	}
	var list []entity.IamUserMonthlyStat
	err := q.OrderDesc(columns.YearMonth).Limit(100).Scan(&list)
	return list, err
}
