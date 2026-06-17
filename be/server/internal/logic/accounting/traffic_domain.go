package accounting

import (
	"context"

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

	// 2. Update daily and monthly stat (Raw SQL for ON DUPLICATE KEY UPDATE)
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

	// Update daily stat
	_, err = g.DB().Exec(ctx, `
		INSERT INTO iam_user_daily_stat (user_id, date, uploaded, downloaded, seed_time, leech_time, bonus, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, 0, ?)
		ON DUPLICATE KEY UPDATE 
			uploaded = uploaded + VALUES(uploaded),
			downloaded = downloaded + VALUES(downloaded),
			seed_time = seed_time + VALUES(seed_time),
			leech_time = leech_time + VALUES(leech_time)
	`, userId, dateStr, diffUp, diffDn, seedTimeDiff, leechTimeDiff, eventTime)
	if err != nil {
		return err
	}

	// Update monthly stat
	_, err = g.DB().Exec(ctx, `
		INSERT INTO iam_user_monthly_stat (user_id, year_month, uploaded, downloaded, seed_time, leech_time, bonus, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, 0, ?)
		ON DUPLICATE KEY UPDATE 
			uploaded = uploaded + VALUES(uploaded),
			downloaded = downloaded + VALUES(downloaded),
			seed_time = seed_time + VALUES(seed_time),
			leech_time = leech_time + VALUES(leech_time)
	`, userId, monthStr, diffUp, diffDn, seedTimeDiff, leechTimeDiff, eventTime)

	return err
}

func (s *sAccountingTrafficDomain) QueryDailyStats(ctx context.Context, userId uint64, startDate, endDate *gtime.Time) ([]entity.IamUserDailyStat, error) {
	q := dao.IamUserDailyStat.Ctx(ctx).Where(dao.IamUserDailyStat.Columns().UserId, userId)
	if startDate != nil {
		q = q.Where("date >= ?", startDate.Format("Y-m-d"))
	}
	if endDate != nil {
		q = q.Where("date <= ?", endDate.Format("Y-m-d"))
	}
	var list []entity.IamUserDailyStat
	err := q.OrderDesc(dao.IamUserDailyStat.Columns().Date).Limit(100).Scan(&list)
	return list, err
}

func (s *sAccountingTrafficDomain) QueryMonthlyStats(ctx context.Context, userId uint64, startDate, endDate *gtime.Time) ([]entity.IamUserMonthlyStat, error) {
	q := dao.IamUserMonthlyStat.Ctx(ctx).Where(dao.IamUserMonthlyStat.Columns().UserId, userId)
	if startDate != nil {
		q = q.Where("year_month >= ?", startDate.Format("Y-m"))
	}
	if endDate != nil {
		q = q.Where("year_month <= ?", endDate.Format("Y-m"))
	}
	var list []entity.IamUserMonthlyStat
	err := q.OrderDesc(dao.IamUserMonthlyStat.Columns().YearMonth).Limit(100).Scan(&list)
	return list, err
}
