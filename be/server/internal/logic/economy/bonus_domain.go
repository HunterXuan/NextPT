package economy

import (
	"context"
	"math"
	"time"

	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/economyin"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sEconomyBonusDomain struct{}

func init() {
	service.RegisterEconomyBonusDomain(NewEconomyBonusDomain())
}

func NewEconomyBonusDomain() *sEconomyBonusDomain {
	return &sEconomyBonusDomain{}
}

func (s *sEconomyBonusDomain) GetUserBonus(ctx context.Context, userId uint64) (float64, error) {
	var userStat entity.IamUserStat
	err := dao.IamUserStat.Ctx(ctx).Where(dao.IamUserStat.Columns().UserId, userId).Scan(&userStat)
	if err != nil || userStat.Id == 0 {
		return 0, gerror.New(gi18n.T(ctx, "economy.bonus.user_stat_error"))
	}
	return userStat.Bonus, nil
}

func (s *sEconomyBonusDomain) DebitBonusIfEnough(ctx context.Context, userId uint64, amount float64) error {
	res, err := dao.IamUserStat.Ctx(ctx).
		Where(dao.IamUserStat.Columns().UserId, userId).
		Where("bonus >= ?", amount).
		Decrement(dao.IamUserStat.Columns().Bonus, amount)
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "economy.bonus.deduct_failed"))
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return gerror.New(gi18n.T(ctx, "economy.bonus.insufficient_balance"))
	}
	return nil
}

func (s *sEconomyBonusDomain) CreditBonus(ctx context.Context, userId uint64, amount float64) error {
	_, err := dao.IamUserStat.Ctx(ctx).
		Where(dao.IamUserStat.Columns().UserId, userId).
		Increment(dao.IamUserStat.Columns().Bonus, amount)
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "economy.bonus.add_failed"))
	}
	return nil
}

func (s *sEconomyBonusDomain) InsertBonusLogs(ctx context.Context, logs []entity.EconomyBonusLog) error {
	doLogs := make([]do.EconomyBonusLog, len(logs))
	for i, log := range logs {
		var period interface{} = log.Period
		if log.Period == "" {
			period = nil
		}
		doLogs[i] = do.EconomyBonusLog{
			UserId:       log.UserId,
			Amount:       log.Amount,
			BalanceAfter: log.BalanceAfter,
			Action:       log.Action,
			TargetType:   log.TargetType,
			TargetId:     log.TargetId,
			Period:       period,
			Remark:       log.Remark,
			CreatedAt:    log.CreatedAt,
		}
	}
	_, err := dao.EconomyBonusLog.Ctx(ctx).Data(doLogs).Insert()
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "economy.bonus.write_log_failed"))
	}
	return nil
}

func (s *sEconomyBonusDomain) InsertBonusLog(ctx context.Context, log entity.EconomyBonusLog) error {
	var period interface{} = log.Period
	if log.Period == "" {
		period = nil
	}
	_, err := dao.EconomyBonusLog.Ctx(ctx).Data(do.EconomyBonusLog{
		UserId:       log.UserId,
		Amount:       log.Amount,
		BalanceAfter: log.BalanceAfter,
		Action:       log.Action,
		TargetType:   log.TargetType,
		TargetId:     log.TargetId,
		Period:       period,
		Remark:       log.Remark,
		CreatedAt:    log.CreatedAt,
	}).Insert()
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "economy.bonus.write_log_failed"))
	}
	return nil
}

func (s *sEconomyBonusDomain) CalculateBonusForPeers(peers []economyin.BonusPeerSnapshot, torrentMap map[uint64]economyin.BonusTorrentSnapshot, config economyin.BonusFormulaConfig, now time.Time) float64 {
	var scoreA float64
	var validCount int

	for _, p := range peers {
		t, ok := torrentMap[p.TorrentId]
		if !ok {
			continue
		}

		sizeGB := float64(t.Size) / (1024 * 1024 * 1024)
		if sizeGB >= 1.0 && validCount < 50 {
			validCount++
		}

		var Ti float64
		if t.CreatedAt != nil {
			Ti = now.Sub(t.CreatedAt.Time).Hours() / (24 * 7)
		}

		Ni := float64(t.Seeders)
		if Ni < 1 {
			Ni = 1
		}

		part1 := 1.0 - math.Pow(10, -Ti/config.T0)
		part3 := 1.0 + math.Sqrt(2)*math.Pow(10, -(Ni-1)/(config.N0-1))
		A := part1 * sizeGB * part3
		scoreA += A
	}

	B := config.B0 * (2.0 / math.Pi) * math.Atan(scoreA/config.L)
	return config.BasePoints*float64(validCount) + B
}

func (s *sEconomyBonusDomain) QueryBonusLogs(ctx context.Context, userId uint64, action string, page, size int) ([]entity.EconomyBonusLog, int, error) {
	m := dao.EconomyBonusLog.Ctx(ctx).Where(dao.EconomyBonusLog.Columns().UserId, userId)
	if action != "" {
		m = m.WhereLike(dao.EconomyBonusLog.Columns().Action, "%"+action+"%")
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var logs []entity.EconomyBonusLog
	err = m.Page(page, size).OrderDesc(dao.EconomyBonusLog.Columns().CreatedAt).Scan(&logs)
	return logs, total, err
}
