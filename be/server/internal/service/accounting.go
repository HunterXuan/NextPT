// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/accountingin"
	"server/internal/model/out/accountingout"

	"github.com/gogf/gf/v2/os/gtime"
)

type (
	IAccountingPeerUsecase interface {
		ListMyPeers(ctx context.Context, actor *model.Actor, in accountingin.PeerListInp) (*accountingout.PeerListOut, error)
	}
	IAccountingSnatchDomain interface {
		RecordSnatch(ctx context.Context, in accountingin.RecordSnatchInp) (bool, error)
		ListSnatches(ctx context.Context, userId uint64, page int, size int, isFinished *bool) ([]*entity.TrackerSnatch, int, error)
		GetSnatch(ctx context.Context, userId uint64, torrentId uint64) (*entity.TrackerSnatch, error)
		DeleteSnatchesByTorrentId(ctx context.Context, torrentId uint64) error
	}
	IAccountingSnatchUsecase interface {
		ListMySnatches(ctx context.Context, actor *model.Actor, in accountingin.SnatchListInp) (*accountingout.SnatchListOut, error)
		GetMySnatch(ctx context.Context, actor *model.Actor, in accountingin.SnatchGetInp) (*accountingout.SnatchGetOut, error)
	}
	IAccountingTrafficDomain interface {
		GetUserStat(ctx context.Context, userId uint64) (*entity.IamUserStat, error)
		RecordTraffic(ctx context.Context, in accountingin.RecordTrafficInp) error
		QueryPeriodStats(ctx context.Context, userId uint64, periodType int, startDate *gtime.Time, endDate *gtime.Time) ([]entity.IamUserPeriodStat, error)
	}
	IAccountingTrafficUsecase interface {
		GetMyTraffic(ctx context.Context, actor *model.Actor) (*accountingout.TrafficGetMeOut, error)
		ListMyTrafficHistory(ctx context.Context, actor *model.Actor, in accountingin.TrafficHistoryListInp) (*accountingout.TrafficHistoryListOut, error)
	}
)

var (
	localAccountingPeerUsecase    IAccountingPeerUsecase
	localAccountingSnatchDomain   IAccountingSnatchDomain
	localAccountingSnatchUsecase  IAccountingSnatchUsecase
	localAccountingTrafficDomain  IAccountingTrafficDomain
	localAccountingTrafficUsecase IAccountingTrafficUsecase
)

func AccountingPeerUsecase() IAccountingPeerUsecase {
	if localAccountingPeerUsecase == nil {
		panic("implement not found for interface IAccountingPeerUsecase, forgot register?")
	}
	return localAccountingPeerUsecase
}

func RegisterAccountingPeerUsecase(i IAccountingPeerUsecase) {
	localAccountingPeerUsecase = i
}

func AccountingSnatchDomain() IAccountingSnatchDomain {
	if localAccountingSnatchDomain == nil {
		panic("implement not found for interface IAccountingSnatchDomain, forgot register?")
	}
	return localAccountingSnatchDomain
}

func RegisterAccountingSnatchDomain(i IAccountingSnatchDomain) {
	localAccountingSnatchDomain = i
}

func AccountingSnatchUsecase() IAccountingSnatchUsecase {
	if localAccountingSnatchUsecase == nil {
		panic("implement not found for interface IAccountingSnatchUsecase, forgot register?")
	}
	return localAccountingSnatchUsecase
}

func RegisterAccountingSnatchUsecase(i IAccountingSnatchUsecase) {
	localAccountingSnatchUsecase = i
}

func AccountingTrafficDomain() IAccountingTrafficDomain {
	if localAccountingTrafficDomain == nil {
		panic("implement not found for interface IAccountingTrafficDomain, forgot register?")
	}
	return localAccountingTrafficDomain
}

func RegisterAccountingTrafficDomain(i IAccountingTrafficDomain) {
	localAccountingTrafficDomain = i
}

func AccountingTrafficUsecase() IAccountingTrafficUsecase {
	if localAccountingTrafficUsecase == nil {
		panic("implement not found for interface IAccountingTrafficUsecase, forgot register?")
	}
	return localAccountingTrafficUsecase
}

func RegisterAccountingTrafficUsecase(i IAccountingTrafficUsecase) {
	localAccountingTrafficUsecase = i
}
