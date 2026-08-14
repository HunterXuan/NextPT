// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/modin"
	"server/internal/model/out/modout"

	"github.com/gogf/gf/v2/os/gtime"
)

type (
	IModCheaterDomain interface {
		Create(ctx context.Context, log entity.ModCheaterLog) error
		Resolve(ctx context.Context, id uint64, dealtBy uint64, dealtComment string, dealtAt *gtime.Time) error
		QueryCheaterLogs(ctx context.Context, isDealt *int, page int, size int) ([]entity.ModCheaterLog, int, error)
		DeleteCheaterLogsByTorrentId(ctx context.Context, torrentId uint64) error
	}
	IModCheaterUsecase interface {
		Record(ctx context.Context, in modin.RecordCheaterLogInp) error
		List(ctx context.Context, actor *model.Actor, in modin.ListCheaterLogsInp) (*modout.ListCheaterLogsOut, error)
		Resolve(ctx context.Context, actor *model.Actor, in modin.ResolveCheaterLogInp) error
	}
	IModReportDomain interface {
		GetPendingCount(ctx context.Context, reporterId uint64, targetType string, targetId uint64) (int, error)
		Create(ctx context.Context, report entity.ModReport) error
		GetById(ctx context.Context, id uint64) (*entity.ModReport, error)
		Update(ctx context.Context, id uint64, data interface{}) error
		QueryReports(ctx context.Context, status int, targetType string, page int, size int) ([]entity.ModReport, int, error)
		DeleteReportsByTarget(ctx context.Context, targetType string, targetId uint64) error
		DeleteReportsByTargets(ctx context.Context, targetType string, targetIds []uint64) error
	}
	IModReportUsecase interface {
		Create(ctx context.Context, actor *model.Actor, in modin.CreateReportInp) error
		List(ctx context.Context, actor *model.Actor, in modin.ListReportsInp) (*modout.ListReportsOut, error)
		Resolve(ctx context.Context, actor *model.Actor, in modin.ResolveReportInp) error
	}
	IModStaffMessageDomain interface {
		Create(ctx context.Context, message entity.ModStaffMessage) (uint64, error)
		GetById(ctx context.Context, id uint64) (*entity.ModStaffMessage, error)
		ListBySender(ctx context.Context, senderId uint64, in modin.StaffMessageListInp) ([]entity.ModStaffMessage, int, error)
		AdminList(ctx context.Context, in modin.AdminStaffMessageListInp) ([]entity.ModStaffMessage, int, error)
		Update(ctx context.Context, id uint64, data interface{}) error
	}
	IModStaffMessageUsecase interface {
		Create(ctx context.Context, actor *model.Actor, in modin.StaffMessageCreateInp) (*modout.StaffMessageCreateOut, error)
		List(ctx context.Context, actor *model.Actor, in modin.StaffMessageListInp) (*modout.StaffMessageListOut, error)
		AdminList(ctx context.Context, actor *model.Actor, in modin.AdminStaffMessageListInp) (*modout.StaffMessageListOut, error)
		AdminProcess(ctx context.Context, actor *model.Actor, in modin.StaffMessageProcessInp) error
	}
	IModUserDomain interface {
		Create(ctx context.Context, mod entity.ModUserLog) (uint64, error)
		GetById(ctx context.Context, id uint64) (*entity.ModUserLog, error)
		Update(ctx context.Context, id uint64, data interface{}) error
		QueryExpiredActiveMods(ctx context.Context, now *gtime.Time, limit int) ([]entity.ModUserLog, error)
		HasActiveMod(ctx context.Context, userId uint64, modTypes []int) (bool, error)
		QueryUserLogs(ctx context.Context, userId uint64, page int, size int) ([]entity.ModUserLog, int, error)
	}
	IModUserUsecase interface {
		Apply(ctx context.Context, actor *model.Actor, in modin.ApplyModInp) error
		Remove(ctx context.Context, actor *model.Actor, in modin.RemoveModInp) error
		CleanupExpired(ctx context.Context) (int, error)
		List(ctx context.Context, actor *model.Actor, in modin.ListUserInp) (*modout.ListUserOut, error)
	}
)

var (
	localModCheaterDomain       IModCheaterDomain
	localModCheaterUsecase      IModCheaterUsecase
	localModReportDomain        IModReportDomain
	localModReportUsecase       IModReportUsecase
	localModStaffMessageDomain  IModStaffMessageDomain
	localModStaffMessageUsecase IModStaffMessageUsecase
	localModUserDomain          IModUserDomain
	localModUserUsecase         IModUserUsecase
)

func ModCheaterDomain() IModCheaterDomain {
	if localModCheaterDomain == nil {
		panic("implement not found for interface IModCheaterDomain, forgot register?")
	}
	return localModCheaterDomain
}

func RegisterModCheaterDomain(i IModCheaterDomain) {
	localModCheaterDomain = i
}

func ModCheaterUsecase() IModCheaterUsecase {
	if localModCheaterUsecase == nil {
		panic("implement not found for interface IModCheaterUsecase, forgot register?")
	}
	return localModCheaterUsecase
}

func RegisterModCheaterUsecase(i IModCheaterUsecase) {
	localModCheaterUsecase = i
}

func ModReportDomain() IModReportDomain {
	if localModReportDomain == nil {
		panic("implement not found for interface IModReportDomain, forgot register?")
	}
	return localModReportDomain
}

func RegisterModReportDomain(i IModReportDomain) {
	localModReportDomain = i
}

func ModReportUsecase() IModReportUsecase {
	if localModReportUsecase == nil {
		panic("implement not found for interface IModReportUsecase, forgot register?")
	}
	return localModReportUsecase
}

func RegisterModReportUsecase(i IModReportUsecase) {
	localModReportUsecase = i
}

func ModStaffMessageDomain() IModStaffMessageDomain {
	if localModStaffMessageDomain == nil {
		panic("implement not found for interface IModStaffMessageDomain, forgot register?")
	}
	return localModStaffMessageDomain
}

func RegisterModStaffMessageDomain(i IModStaffMessageDomain) {
	localModStaffMessageDomain = i
}

func ModStaffMessageUsecase() IModStaffMessageUsecase {
	if localModStaffMessageUsecase == nil {
		panic("implement not found for interface IModStaffMessageUsecase, forgot register?")
	}
	return localModStaffMessageUsecase
}

func RegisterModStaffMessageUsecase(i IModStaffMessageUsecase) {
	localModStaffMessageUsecase = i
}

func ModUserDomain() IModUserDomain {
	if localModUserDomain == nil {
		panic("implement not found for interface IModUserDomain, forgot register?")
	}
	return localModUserDomain
}

func RegisterModUserDomain(i IModUserDomain) {
	localModUserDomain = i
}

func ModUserUsecase() IModUserUsecase {
	if localModUserUsecase == nil {
		panic("implement not found for interface IModUserUsecase, forgot register?")
	}
	return localModUserUsecase
}

func RegisterModUserUsecase(i IModUserUsecase) {
	localModUserUsecase = i
}
