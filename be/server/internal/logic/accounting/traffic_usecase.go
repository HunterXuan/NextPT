package accounting

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/accountingin"
	"server/internal/model/out/accountingout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/util/gconv"
)

type sAccountingTrafficUsecase struct{}

func init() {
	service.RegisterAccountingTrafficUsecase(NewAccountingTrafficUsecase())
}

func NewAccountingTrafficUsecase() *sAccountingTrafficUsecase {
	return &sAccountingTrafficUsecase{}
}

func (s *sAccountingTrafficUsecase) GetMyTraffic(ctx context.Context, actor *model.Actor) (*accountingout.TrafficGetMeOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	stat, err := service.AccountingTrafficDomain().GetUserStat(ctx, actor.Id)
	if err != nil {
		return nil, err
	}

	var ratio float64
	if stat != nil {
		if stat.Downloaded == 0 && stat.Uploaded > 0 {
			ratio = 9999.9 // Infinity representation
		} else if stat.Downloaded > 0 {
			ratio = float64(stat.Uploaded) / float64(stat.Downloaded)
		}

		return &accountingout.TrafficGetMeOut{
			Uploaded:   stat.Uploaded,
			Downloaded: stat.Downloaded,
			ShareRatio: ratio,
			SeedTime:   stat.SeedTime,
			LeechTime:  stat.LeechTime,
		}, nil
	}

	return &accountingout.TrafficGetMeOut{}, nil
}

func (s *sAccountingTrafficUsecase) ListMyTrafficHistory(ctx context.Context, actor *model.Actor, in accountingin.TrafficHistoryListInp) (*accountingout.TrafficHistoryListOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	var outList []accountingout.TrafficHistoryItem

	if in.Period == "monthly" {
		list, err := service.AccountingTrafficDomain().QueryMonthlyStats(ctx, actor.Id, in.StartDate, in.EndDate)
		if err != nil {
			return nil, err
		}
		for _, item := range list {
			outList = append(outList, accountingout.TrafficHistoryItem{
				Date:       item.YearMonth,
				Uploaded:   item.Uploaded,
				Downloaded: item.Downloaded,
				SeedTime:   item.SeedTime,
				LeechTime:  item.LeechTime,
				Bonus:      gconv.String(item.Bonus),
			})
		}
	} else {
		list, err := service.AccountingTrafficDomain().QueryDailyStats(ctx, actor.Id, in.StartDate, in.EndDate)
		if err != nil {
			return nil, err
		}
		for _, item := range list {
			outList = append(outList, accountingout.TrafficHistoryItem{
				Date:       item.Date.Format("Y-m-d"),
				Uploaded:   item.Uploaded,
				Downloaded: item.Downloaded,
				SeedTime:   item.SeedTime,
				LeechTime:  item.LeechTime,
				Bonus:      gconv.String(item.Bonus),
			})
		}
	}

	if outList == nil {
		outList = []accountingout.TrafficHistoryItem{}
	}

	return &accountingout.TrafficHistoryListOut{List: outList}, nil
}
