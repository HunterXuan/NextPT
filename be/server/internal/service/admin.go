// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"
	"server/internal/model/out/modout"
)

type (
	IAdminCatalogCategoryUsecase interface {
		Create(ctx context.Context, actor *model.Actor, in adminin.CatalogCategoryCreateInp) error
		Update(ctx context.Context, actor *model.Actor, in adminin.CatalogCategoryUpdateInp) error
		Delete(ctx context.Context, actor *model.Actor, in adminin.CatalogCategoryDeleteInp) error
		List(ctx context.Context, actor *model.Actor, in adminin.CatalogCategoryListInp) (*adminout.CatalogCategoryListOut, error)
	}
	IAdminCatalogTorrentUsecase interface {
		Delete(ctx context.Context, actor *model.Actor, in adminin.CatalogTorrentDeleteInp) error
	}
	IAdminForumCategoryUsecase interface {
		Create(ctx context.Context, actor *model.Actor, in adminin.ForumCategoryCreateInp) error
		Update(ctx context.Context, actor *model.Actor, in adminin.ForumCategoryUpdateInp) error
		Delete(ctx context.Context, actor *model.Actor, in adminin.ForumCategoryDeleteInp) error
		List(ctx context.Context, actor *model.Actor, in adminin.ForumCategoryListInp) (*adminout.ForumCategoryListOut, error)
	}
	IAdminForumNodeUsecase interface {
		Create(ctx context.Context, actor *model.Actor, in adminin.ForumNodeCreateInp) error
		Update(ctx context.Context, actor *model.Actor, in adminin.ForumNodeUpdateInp) error
		Delete(ctx context.Context, actor *model.Actor, in adminin.ForumNodeDeleteInp) error
		List(ctx context.Context, actor *model.Actor, in adminin.ForumNodeListInp) (*adminout.ForumNodeListOut, error)
	}
	IAdminForumTopicUsecase interface {
		Lock(ctx context.Context, actor *model.Actor, in adminin.ForumTopicLockInp) error
		Unlock(ctx context.Context, actor *model.Actor, in adminin.ForumTopicUnlockInp) error
		Pin(ctx context.Context, actor *model.Actor, in adminin.ForumTopicPinInp) error
		Unpin(ctx context.Context, actor *model.Actor, in adminin.ForumTopicUnpinInp) error
		Move(ctx context.Context, actor *model.Actor, in adminin.ForumTopicMoveInp) error
		Delete(ctx context.Context, actor *model.Actor, in adminin.ForumTopicDeleteInp) error
	}
	IAdminIamInviteUsecase interface {
		Grant(ctx context.Context, actor *model.Actor, in adminin.IamInviteGrantInp) error
	}
	IAdminIamPermissionUsecase interface {
		GrantUserAcl(ctx context.Context, actor *model.Actor, in adminin.IamUserPermissionGrantInp) error
		RevokeUserAcl(ctx context.Context, actor *model.Actor, in adminin.IamUserPermissionRevokeInp) error
		List(ctx context.Context, actor *model.Actor, in adminin.IamPermissionListInp) (*adminout.IamPermissionListOut, error)
	}
	IAdminIamRoleUsecase interface {
		List(ctx context.Context, actor *model.Actor) (*adminout.IamRoleListOut, error)
		Create(ctx context.Context, actor *model.Actor, in adminin.IamRoleCreateInp) (uint, error)
		Update(ctx context.Context, actor *model.Actor, in adminin.IamRoleUpdateInp) error
		Delete(ctx context.Context, actor *model.Actor, in adminin.IamRoleDeleteInp) error
	}
	IAdminIamSessionUsecase interface {
		DeleteByUser(ctx context.Context, actor *model.Actor, in adminin.IamSessionDeleteInp) error
	}
	IAdminIamUserUsecase interface {
		List(ctx context.Context, actor *model.Actor, in adminin.IamUserListInp) (*adminout.IamUserListOut, error)
		Update(ctx context.Context, actor *model.Actor, in adminin.IamUserUpdateInp) error
		StatDetail(ctx context.Context, actor *model.Actor, in adminin.IamUserStatDetailInp) (*adminout.IamUserStatDetailOut, error)
		StatUpdate(ctx context.Context, actor *model.Actor, in adminin.IamUserStatUpdateInp) error
		Ban(ctx context.Context, actor *model.Actor, in adminin.IamUserBanInp) error
	}
	IAdminModCheaterUsecase interface {
		List(ctx context.Context, actor *model.Actor, in adminin.ModCheaterListInp) (*modout.ListCheaterLogsOut, error)
		Resolve(ctx context.Context, actor *model.Actor, in adminin.ModCheaterResolveInp) error
	}
	IAdminModReportUsecase interface {
		List(ctx context.Context, actor *model.Actor, in adminin.ModReportListInp) (*modout.ListReportsOut, error)
		Resolve(ctx context.Context, actor *model.Actor, in adminin.ModReportResolveInp) error
	}
	IAdminModUserUsecase interface {
		List(ctx context.Context, actor *model.Actor, in adminin.ModUserListInp) (*modout.ListUserOut, error)
		Apply(ctx context.Context, actor *model.Actor, in adminin.ModUserApplyInp) error
		Remove(ctx context.Context, actor *model.Actor, in adminin.ModUserRemoveInp) error
	}
	IAdminSiteAuditUsecase interface {
		List(ctx context.Context, actor *model.Actor, in adminin.SiteAuditListInp) (*adminout.SiteAuditListOut, error)
	}
	IAdminSiteConfigUsecase interface {
		List(ctx context.Context, actor *model.Actor, in adminin.SiteConfigListInp) (*adminout.SiteConfigListOut, error)
		Update(ctx context.Context, actor *model.Actor, in adminin.SiteConfigUpdateInp) error
	}
	IAdminSysCronUsecase interface {
		List(ctx context.Context, actor *model.Actor) (*adminout.SysCronListOut, error)
		LogList(ctx context.Context, actor *model.Actor, jobName string, page int, size int) (*adminout.SysCronLogListOut, error)
	}
)

var (
	localAdminCatalogCategoryUsecase IAdminCatalogCategoryUsecase
	localAdminCatalogTorrentUsecase  IAdminCatalogTorrentUsecase
	localAdminForumCategoryUsecase   IAdminForumCategoryUsecase
	localAdminForumNodeUsecase       IAdminForumNodeUsecase
	localAdminForumTopicUsecase      IAdminForumTopicUsecase
	localAdminIamInviteUsecase       IAdminIamInviteUsecase
	localAdminIamPermissionUsecase   IAdminIamPermissionUsecase
	localAdminIamRoleUsecase         IAdminIamRoleUsecase
	localAdminIamSessionUsecase      IAdminIamSessionUsecase
	localAdminIamUserUsecase         IAdminIamUserUsecase
	localAdminModCheaterUsecase      IAdminModCheaterUsecase
	localAdminModReportUsecase       IAdminModReportUsecase
	localAdminModUserUsecase         IAdminModUserUsecase
	localAdminSiteAuditUsecase       IAdminSiteAuditUsecase
	localAdminSiteConfigUsecase      IAdminSiteConfigUsecase
	localAdminSysCronUsecase         IAdminSysCronUsecase
)

func AdminCatalogCategoryUsecase() IAdminCatalogCategoryUsecase {
	if localAdminCatalogCategoryUsecase == nil {
		panic("implement not found for interface IAdminCatalogCategoryUsecase, forgot register?")
	}
	return localAdminCatalogCategoryUsecase
}

func RegisterAdminCatalogCategoryUsecase(i IAdminCatalogCategoryUsecase) {
	localAdminCatalogCategoryUsecase = i
}

func AdminCatalogTorrentUsecase() IAdminCatalogTorrentUsecase {
	if localAdminCatalogTorrentUsecase == nil {
		panic("implement not found for interface IAdminCatalogTorrentUsecase, forgot register?")
	}
	return localAdminCatalogTorrentUsecase
}

func RegisterAdminCatalogTorrentUsecase(i IAdminCatalogTorrentUsecase) {
	localAdminCatalogTorrentUsecase = i
}

func AdminForumCategoryUsecase() IAdminForumCategoryUsecase {
	if localAdminForumCategoryUsecase == nil {
		panic("implement not found for interface IAdminForumCategoryUsecase, forgot register?")
	}
	return localAdminForumCategoryUsecase
}

func RegisterAdminForumCategoryUsecase(i IAdminForumCategoryUsecase) {
	localAdminForumCategoryUsecase = i
}

func AdminForumNodeUsecase() IAdminForumNodeUsecase {
	if localAdminForumNodeUsecase == nil {
		panic("implement not found for interface IAdminForumNodeUsecase, forgot register?")
	}
	return localAdminForumNodeUsecase
}

func RegisterAdminForumNodeUsecase(i IAdminForumNodeUsecase) {
	localAdminForumNodeUsecase = i
}

func AdminForumTopicUsecase() IAdminForumTopicUsecase {
	if localAdminForumTopicUsecase == nil {
		panic("implement not found for interface IAdminForumTopicUsecase, forgot register?")
	}
	return localAdminForumTopicUsecase
}

func RegisterAdminForumTopicUsecase(i IAdminForumTopicUsecase) {
	localAdminForumTopicUsecase = i
}

func AdminIamInviteUsecase() IAdminIamInviteUsecase {
	if localAdminIamInviteUsecase == nil {
		panic("implement not found for interface IAdminIamInviteUsecase, forgot register?")
	}
	return localAdminIamInviteUsecase
}

func RegisterAdminIamInviteUsecase(i IAdminIamInviteUsecase) {
	localAdminIamInviteUsecase = i
}

func AdminIamPermissionUsecase() IAdminIamPermissionUsecase {
	if localAdminIamPermissionUsecase == nil {
		panic("implement not found for interface IAdminIamPermissionUsecase, forgot register?")
	}
	return localAdminIamPermissionUsecase
}

func RegisterAdminIamPermissionUsecase(i IAdminIamPermissionUsecase) {
	localAdminIamPermissionUsecase = i
}

func AdminIamRoleUsecase() IAdminIamRoleUsecase {
	if localAdminIamRoleUsecase == nil {
		panic("implement not found for interface IAdminIamRoleUsecase, forgot register?")
	}
	return localAdminIamRoleUsecase
}

func RegisterAdminIamRoleUsecase(i IAdminIamRoleUsecase) {
	localAdminIamRoleUsecase = i
}

func AdminIamSessionUsecase() IAdminIamSessionUsecase {
	if localAdminIamSessionUsecase == nil {
		panic("implement not found for interface IAdminIamSessionUsecase, forgot register?")
	}
	return localAdminIamSessionUsecase
}

func RegisterAdminIamSessionUsecase(i IAdminIamSessionUsecase) {
	localAdminIamSessionUsecase = i
}

func AdminIamUserUsecase() IAdminIamUserUsecase {
	if localAdminIamUserUsecase == nil {
		panic("implement not found for interface IAdminIamUserUsecase, forgot register?")
	}
	return localAdminIamUserUsecase
}

func RegisterAdminIamUserUsecase(i IAdminIamUserUsecase) {
	localAdminIamUserUsecase = i
}

func AdminModCheaterUsecase() IAdminModCheaterUsecase {
	if localAdminModCheaterUsecase == nil {
		panic("implement not found for interface IAdminModCheaterUsecase, forgot register?")
	}
	return localAdminModCheaterUsecase
}

func RegisterAdminModCheaterUsecase(i IAdminModCheaterUsecase) {
	localAdminModCheaterUsecase = i
}

func AdminModReportUsecase() IAdminModReportUsecase {
	if localAdminModReportUsecase == nil {
		panic("implement not found for interface IAdminModReportUsecase, forgot register?")
	}
	return localAdminModReportUsecase
}

func RegisterAdminModReportUsecase(i IAdminModReportUsecase) {
	localAdminModReportUsecase = i
}

func AdminModUserUsecase() IAdminModUserUsecase {
	if localAdminModUserUsecase == nil {
		panic("implement not found for interface IAdminModUserUsecase, forgot register?")
	}
	return localAdminModUserUsecase
}

func RegisterAdminModUserUsecase(i IAdminModUserUsecase) {
	localAdminModUserUsecase = i
}

func AdminSiteAuditUsecase() IAdminSiteAuditUsecase {
	if localAdminSiteAuditUsecase == nil {
		panic("implement not found for interface IAdminSiteAuditUsecase, forgot register?")
	}
	return localAdminSiteAuditUsecase
}

func RegisterAdminSiteAuditUsecase(i IAdminSiteAuditUsecase) {
	localAdminSiteAuditUsecase = i
}

func AdminSiteConfigUsecase() IAdminSiteConfigUsecase {
	if localAdminSiteConfigUsecase == nil {
		panic("implement not found for interface IAdminSiteConfigUsecase, forgot register?")
	}
	return localAdminSiteConfigUsecase
}

func RegisterAdminSiteConfigUsecase(i IAdminSiteConfigUsecase) {
	localAdminSiteConfigUsecase = i
}

func AdminSysCronUsecase() IAdminSysCronUsecase {
	if localAdminSysCronUsecase == nil {
		panic("implement not found for interface IAdminSysCronUsecase, forgot register?")
	}
	return localAdminSysCronUsecase
}

func RegisterAdminSysCronUsecase(i IAdminSysCronUsecase) {
	localAdminSysCronUsecase = i
}
