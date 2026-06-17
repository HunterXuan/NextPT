// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/economyin"
	"server/internal/model/out/economyout"
	"time"
)

type (
	IEconomyBonusDomain interface {
		GetUserBonus(ctx context.Context, userId uint64) (float64, error)
		DebitBonusIfEnough(ctx context.Context, userId uint64, amount float64) error
		CreditBonus(ctx context.Context, userId uint64, amount float64) error
		InsertBonusLogs(ctx context.Context, logs []entity.EconomyBonusLog) error
		InsertBonusLog(ctx context.Context, log entity.EconomyBonusLog) error
		CalculateBonusForPeers(peers []economyin.BonusPeerSnapshot, torrentMap map[uint64]economyin.BonusTorrentSnapshot, config economyin.BonusFormulaConfig, now time.Time) float64
		QueryBonusLogs(ctx context.Context, userId uint64, action string, page int, size int) ([]entity.EconomyBonusLog, int, error)
	}
	IEconomyBonusUsecase interface {
		ListMyBonusLogs(ctx context.Context, actor *model.Actor, in economyin.BonusLogsInp) (*economyout.BonusLogsOut, error)
		TransferBonus(ctx context.Context, fromUserId uint64, toUserId uint64, amount float64, targetType string, targetId uint64, remarkFrom string, remarkTo string) error
		AddBonus(ctx context.Context, userId uint64, amount float64, action string, targetType string, targetId uint64, remark string, period string) error
		GetMyHourlyBonus(ctx context.Context, actor *model.Actor, in economyin.HourlyBonusInp) (*economyout.HourlyBonusOut, error)
		// DistributeBonusPoints 魔力值自动发放 (NexusPHP Formula)
		DistributeBonusPoints(ctx context.Context) error
		// CalculateHourlyBonus 计算指定用户当前每小时可获得魔力值
		CalculateHourlyBonus(ctx context.Context, userId uint64) (float64, error)
	}
)

var (
	localEconomyBonusDomain  IEconomyBonusDomain
	localEconomyBonusUsecase IEconomyBonusUsecase
)

func EconomyBonusDomain() IEconomyBonusDomain {
	if localEconomyBonusDomain == nil {
		panic("implement not found for interface IEconomyBonusDomain, forgot register?")
	}
	return localEconomyBonusDomain
}

func RegisterEconomyBonusDomain(i IEconomyBonusDomain) {
	localEconomyBonusDomain = i
}

func EconomyBonusUsecase() IEconomyBonusUsecase {
	if localEconomyBonusUsecase == nil {
		panic("implement not found for interface IEconomyBonusUsecase, forgot register?")
	}
	return localEconomyBonusUsecase
}

func RegisterEconomyBonusUsecase(i IEconomyBonusUsecase) {
	localEconomyBonusUsecase = i
}
