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
		TransferBonus(ctx context.Context, fromUserId uint64, toUserId uint64, amount float64, targetType string, targetId uint64) error
		AddBonus(ctx context.Context, userId uint64, amount float64, action string, targetType string, targetId uint64, remark string, period string) error
		GetMyHourlyBonus(ctx context.Context, actor *model.Actor, in economyin.HourlyBonusInp) (*economyout.HourlyBonusOut, error)
		// DistributeBonusPoints 魔力值自动发放 (NexusPHP Formula)
		DistributeBonusPoints(ctx context.Context) error
		// CalculateHourlyBonus 计算指定用户当前每小时可获得魔力值
		CalculateHourlyBonus(ctx context.Context, userId uint64) (float64, error)
	}
	IEconomyRewardDomain interface {
		InsertRewardRecord(ctx context.Context, record entity.EconomyRewardRecord) error
		QueryRewardSummaries(ctx context.Context, targetType string, targetId uint64, page int, size int) ([]economyout.RewardSummary, int, error)
		DeleteRewardRecordsByTarget(ctx context.Context, targetType string, targetId uint64) error
		DeleteRewardRecordsByTargets(ctx context.Context, targetType string, targetIds []uint64) error
	}
	IEconomyShopDomain interface {
		LoadProducts(ctx context.Context) (model.EconomyShopProducts, error)
		InsertOrder(ctx context.Context, order entity.EconomyShopOrder) (uint64, error)
		CompleteOrder(ctx context.Context, id uint64, targetType string, targetId uint64) error
		QueryOrdersByUser(ctx context.Context, userId uint64, page int, size int) ([]entity.EconomyShopOrder, int, error)
	}
	IEconomyShopUsecase interface {
		ListProducts(ctx context.Context, actor *model.Actor, in economyin.ShopProductListInp) (*economyout.ShopProductListOut, error)
		CreateOrder(ctx context.Context, actor *model.Actor, in economyin.ShopOrderCreateInp) (*economyout.ShopOrderCreateOut, error)
		ListMyOrders(ctx context.Context, actor *model.Actor, in economyin.ShopOrderListInp) (*economyout.ShopOrderListOut, error)
	}
)

var (
	localEconomyBonusDomain  IEconomyBonusDomain
	localEconomyBonusUsecase IEconomyBonusUsecase
	localEconomyRewardDomain IEconomyRewardDomain
	localEconomyShopDomain   IEconomyShopDomain
	localEconomyShopUsecase  IEconomyShopUsecase
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

func EconomyRewardDomain() IEconomyRewardDomain {
	if localEconomyRewardDomain == nil {
		panic("implement not found for interface IEconomyRewardDomain, forgot register?")
	}
	return localEconomyRewardDomain
}

func RegisterEconomyRewardDomain(i IEconomyRewardDomain) {
	localEconomyRewardDomain = i
}

func EconomyShopDomain() IEconomyShopDomain {
	if localEconomyShopDomain == nil {
		panic("implement not found for interface IEconomyShopDomain, forgot register?")
	}
	return localEconomyShopDomain
}

func RegisterEconomyShopDomain(i IEconomyShopDomain) {
	localEconomyShopDomain = i
}

func EconomyShopUsecase() IEconomyShopUsecase {
	if localEconomyShopUsecase == nil {
		panic("implement not found for interface IEconomyShopUsecase, forgot register?")
	}
	return localEconomyShopUsecase
}

func RegisterEconomyShopUsecase(i IEconomyShopUsecase) {
	localEconomyShopUsecase = i
}
