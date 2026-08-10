package accounting

import (
	"context"
	"fmt"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/model/in/accountingin"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sAccountingTrafficDomain struct{}

type trafficStatChange struct {
	userId            uint64
	uploadedDiff      int64
	downloadedDiff    int64
	rawUploadedDiff   int64
	rawDownloadedDiff int64
	seedTimeDiff      int
	leechTimeDiff     int
	eventTime         *gtime.Time
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

func (s *sAccountingTrafficDomain) RecordTraffic(ctx context.Context, in accountingin.RecordTrafficInp) error {
	if in.UploadedDiff <= 0 && in.DownloadedDiff <= 0 && in.RawUploadedDiff <= 0 && in.RawDownloadedDiff <= 0 && in.TimeDiff <= 0 {
		return nil
	}

	statColumns := dao.IamUserStat.Columns()
	statModel := dao.IamUserStat.Ctx(ctx)
	statUpdate := g.Map{
		statColumns.Uploaded:      s.incrementByValue(statModel, statColumns.Uploaded, in.UploadedDiff),
		statColumns.Downloaded:    s.incrementByValue(statModel, statColumns.Downloaded, in.DownloadedDiff),
		statColumns.RawUploaded:   s.incrementByValue(statModel, statColumns.RawUploaded, in.RawUploadedDiff),
		statColumns.RawDownloaded: s.incrementByValue(statModel, statColumns.RawDownloaded, in.RawDownloadedDiff),
	}
	if in.IsSeeder {
		statUpdate[statColumns.SeedTime] = s.incrementByValue(statModel, statColumns.SeedTime, int64(in.TimeDiff))
	} else {
		statUpdate[statColumns.LeechTime] = s.incrementByValue(statModel, statColumns.LeechTime, int64(in.TimeDiff))
	}
	res, err := statModel.Where(statColumns.UserId, in.UserId).Update(statUpdate)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return gerror.Newf("user stat record missing for user_id: %d", in.UserId)
	}

	// 2. Update daily and monthly period stat.
	if in.EventTime == nil {
		in.EventTime = gtime.Now()
	}

	var seedTimeDiff, leechTimeDiff int
	if in.IsSeeder {
		seedTimeDiff = in.TimeDiff
	} else {
		leechTimeDiff = in.TimeDiff
	}

	change := trafficStatChange{
		userId:            in.UserId,
		uploadedDiff:      in.UploadedDiff,
		downloadedDiff:    in.DownloadedDiff,
		rawUploadedDiff:   in.RawUploadedDiff,
		rawDownloadedDiff: in.RawDownloadedDiff,
		seedTimeDiff:      seedTimeDiff,
		leechTimeDiff:     leechTimeDiff,
		eventTime:         in.EventTime,
	}

	if err = s.upsertPeriodStat(ctx, consts.AccountingStatPeriodDaily, in.EventTime.Format("Y-m-d"), change); err != nil {
		return err
	}

	return s.upsertPeriodStat(ctx, consts.AccountingStatPeriodMonthly, in.EventTime.Format("Y-m"), change)
}

func (s *sAccountingTrafficDomain) upsertPeriodStat(ctx context.Context, periodType int, periodKey string, change trafficStatChange) error {
	columns := dao.IamUserPeriodStat.Columns()
	m := dao.IamUserPeriodStat.Ctx(ctx)
	_, err := m.
		Data(g.Map{
			columns.UserId:        change.userId,
			columns.PeriodType:    periodType,
			columns.PeriodKey:     periodKey,
			columns.Uploaded:      change.uploadedDiff,
			columns.Downloaded:    change.downloadedDiff,
			columns.RawUploaded:   change.rawUploadedDiff,
			columns.RawDownloaded: change.rawDownloadedDiff,
			columns.SeedTime:      change.seedTimeDiff,
			columns.LeechTime:     change.leechTimeDiff,
			columns.Bonus:         0,
			columns.CreatedAt:     change.eventTime,
		}).
		OnDuplicate(g.Map{
			columns.Uploaded:      s.incrementByInsertedValue(m, columns.Uploaded),
			columns.Downloaded:    s.incrementByInsertedValue(m, columns.Downloaded),
			columns.RawUploaded:   s.incrementByInsertedValue(m, columns.RawUploaded),
			columns.RawDownloaded: s.incrementByInsertedValue(m, columns.RawDownloaded),
			columns.SeedTime:      s.incrementByInsertedValue(m, columns.SeedTime),
			columns.LeechTime:     s.incrementByInsertedValue(m, columns.LeechTime),
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

func (s *sAccountingTrafficDomain) CleanupPeriodStats(ctx context.Context, now *gtime.Time) (int, error) {
	if now == nil {
		now = gtime.Now()
	}
	dailyBefore := gtime.NewFromTime(now.Time.AddDate(0, 0, -consts.AccountingDailyStatRetentionDays))
	monthlyBefore := gtime.NewFromTime(now.Time.AddDate(0, -consts.AccountingMonthlyStatRetentionMonths, 0))
	columns := dao.IamUserPeriodStat.Columns()
	deleted := 0

	rows, err := s.deletePeriodStats(dao.IamUserPeriodStat.Ctx(ctx).
		Where(columns.PeriodType, consts.AccountingStatPeriodDaily).
		WhereLT(columns.PeriodKey, dailyBefore.Format("Y-m-d")))
	if err != nil {
		return deleted, err
	}
	deleted += rows

	rows, err = s.deletePeriodStats(dao.IamUserPeriodStat.Ctx(ctx).
		Where(columns.PeriodType, consts.AccountingStatPeriodMonthly).
		WhereLT(columns.PeriodKey, monthlyBefore.Format("Y-m")))
	if err != nil {
		return deleted, err
	}
	return deleted + rows, nil
}

func (s *sAccountingTrafficDomain) deletePeriodStats(m *gdb.Model) (int, error) {
	result, err := m.Delete()
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return int(rows), err
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
