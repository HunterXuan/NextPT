// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"server/internal/model"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/iamin"
	"server/internal/model/out/iamout"

	"github.com/goflyfox/gtoken/v2/gtoken"
	"github.com/gogf/gf/v2/os/gtime"
)

type (
	IIamInviteDomain interface {
		CreateInvite(ctx context.Context, inviterId uint64, isTemporary bool, expireAt *gtime.Time) (uint64, error)
		GetInviteByHashForUpdate(ctx context.Context, hash string) (*entity.IamInvite, error)
		GetInviteByHash(ctx context.Context, hash string) (*entity.IamInvite, error)
		GetInviteByInviterIdAndId(ctx context.Context, inviterId uint64, id uint64) (*entity.IamInvite, error)
		QueryInvitesByInviter(ctx context.Context, inviterId uint64, page int, size int, status *uint) ([]entity.IamInvite, int, error)
		AdminQuerySiteInvites(ctx context.Context, page int, size int, status *uint) ([]entity.IamInvite, int, error)
		AdminGetSiteInviteById(ctx context.Context, id uint64) (*entity.IamInvite, error)
		AdminUpdateSiteInviteStatus(ctx context.Context, id uint64, status int) error
		ExpireInvites(ctx context.Context, now *gtime.Time) (int64, error)
		UpdateInviteStatus(ctx context.Context, id uint64, status uint) error
		UpdateInvite(ctx context.Context, id uint64, data do.IamInvite) error
	}
	IIamInviteUsecase interface {
		List(ctx context.Context, actor *model.Actor, in iamin.InviteListInp) (*iamout.InviteListOut, error)
		Send(ctx context.Context, actor *model.Actor, in iamin.InviteSendInp) error
		Check(ctx context.Context, in iamin.InviteCheckInp) (*iamout.InviteCheckOut, error)
		CleanupExpired(ctx context.Context) (int64, error)
	}
	IIamLoginLogDomain interface {
		Create(ctx context.Context, data do.IamLoginLog) error
		List(ctx context.Context, options model.IamLoginLogListOptions) ([]entity.IamLoginLog, int, error)
	}
	IIamPermissionDomain interface {
		CheckPermissionWithList(ctx context.Context, rolePerms []string, userAcls []string, permKey string) (bool, error)
		GrantUserPermission(ctx context.Context, userId uint64, permKey string, isDeny bool) error
		GrantUserPermissions(ctx context.Context, userId uint64, permKeys []string, isDeny bool) error
		RevokeUserPermission(ctx context.Context, userId uint64, permKey string, isDeny bool) error
		RevokeUserPermissionsByIds(ctx context.Context, userId uint64, ids []uint64) error
		GrantUserPermissionsBySource(ctx context.Context, userId uint64, permKeys []string, sourceType int, sourceId uint64, expireAt *gtime.Time) error
		DeactivateUserPermissionsBySource(ctx context.Context, userId uint64, sourceType int, sourceId uint64) error
		GetAllPermissions(ctx context.Context) []string
		GetUserPermissions(ctx context.Context, userId uint64) ([]entity.IamUserPermission, error)
		ListUserPermissions(ctx context.Context, userId uint64, options model.IamUserPermissionListOptions) ([]entity.IamUserPermission, int, error)
	}
	IIamRoleDomain interface {
		GetRoleById(ctx context.Context, roleId uint) (*entity.IamRole, error)
		GetRolesByIds(ctx context.Context, roleIds []uint) ([]entity.IamRole, error)
		ListRoles(ctx context.Context) ([]entity.IamRole, error)
		AdminListRoles(ctx context.Context) ([]entity.IamRole, error)
		AdminCreateRole(ctx context.Context, level int, nameI18N []byte, rules []byte, permissions []byte, isStaff bool) (uint, error)
		AdminUpdateRole(ctx context.Context, id uint, level *int, nameI18N []byte, rules []byte, permissions []byte, isStaff *bool) error
		AdminDeleteRole(ctx context.Context, id uint) (int, error)
	}
	IIamRoleUsecase interface {
		List(ctx context.Context, actor *model.Actor) (*iamout.RoleListOut, error)
		SyncRanks(ctx context.Context) (*model.IamRankSyncResult, error)
	}
	IIamSessionDomain interface {
		GetGFToken() gtoken.Token
		GetGFMiddleware() gtoken.Middleware
		GenerateToken(ctx context.Context, userKey string, data any) (string, error)
		RemoveToken(ctx context.Context, userKey string) error
	}
	IIamSessionUsecase interface {
		Create(ctx context.Context, in iamin.SessionCreateInp) (*iamout.SessionCreateOut, error)
		VerifyTwoStep(ctx context.Context, in iamin.SessionTwoStepVerifyInp) (*iamout.SessionCreateOut, error)
		Delete(ctx context.Context, actor *model.Actor) error
		VerifyPasskey(ctx context.Context, passkey string) (*model.Actor, error)
	}
	IIamTwoStepDomain interface {
		ReplaceRecoveryCodeHashes(ctx context.Context, userId uint64, hashes []string) error
		HasUnusedRecoveryCode(ctx context.Context, userId uint64, codeHash string) (bool, error)
		ConsumeRecoveryCode(ctx context.Context, userId uint64, codeHash string) (bool, error)
		DeleteRecoveryCodes(ctx context.Context, userId uint64) error
	}
	IIamTwoStepUsecase interface {
		Setup(ctx context.Context, actor *model.Actor, in iamin.UserTwoStepSetupInp) (*iamout.UserTwoStepSetupOut, error)
		Confirm(ctx context.Context, actor *model.Actor, in iamin.UserTwoStepConfirmInp) (*iamout.UserTwoStepRecoveryCodesOut, error)
		CreateRecoveryCodes(ctx context.Context, actor *model.Actor, in iamin.UserTwoStepRecoveryCodesCreateInp) (*iamout.UserTwoStepRecoveryCodesOut, error)
		Disable(ctx context.Context, actor *model.Actor, in iamin.UserTwoStepDisableInp) error
		CreateLoginChallenge(ctx context.Context, userId uint64) (string, error)
		VerifyLogin(ctx context.Context, in iamin.SessionTwoStepVerifyInp) (uint64, error)
	}
	IIamUserDomain interface {
		GetUserByLogin(ctx context.Context, login string) (*entity.IamUser, error)
		GetUserByEmail(ctx context.Context, email string) (*entity.IamUser, error)
		GetUserById(ctx context.Context, id uint64) (*entity.IamUser, error)
		GetUserPasswordHash(ctx context.Context, userId uint64) (string, error)
		GetUserProfile(ctx context.Context, userId uint64) (*entity.IamUserProfile, error)
		GetUserStat(ctx context.Context, userId uint64) (*entity.IamUserStat, error)
		CheckUsernameExists(ctx context.Context, username string) (bool, error)
		CheckEmailExists(ctx context.Context, email string) (bool, error)
		InsertUser(ctx context.Context, data do.IamUser) (uint64, error)
		InsertUserProfile(ctx context.Context, data do.IamUserProfile) error
		InsertUserStat(ctx context.Context, data do.IamUserStat) error
		UpdateUserProfile(ctx context.Context, userId uint64, avatar string, info string, signature string) error
		ConfirmUserEmail(ctx context.Context, userId uint64) (bool, error)
		UpdatePasswordHash(ctx context.Context, userId uint64, passwordHash string) error
		UpdatePasskey(ctx context.Context, userId uint64, passkey string) error
		UpdateTwoStep(ctx context.Context, userId uint64, twoStepType int, encryptedSecret string) error
		UpdateLoginTrace(ctx context.Context, userId uint64, loginAt *gtime.Time, ip string) error
		// GetUserByPasskey 通过 Passkey 获取用户（无缓存，纯领域逻辑）
		GetUserByPasskey(ctx context.Context, passkey string) (*entity.IamUser, error)
		AdminListUsers(ctx context.Context, search string, order string, page int, size int) ([]*entity.IamUser, int, error)
		AdminUpdateUser(ctx context.Context, id uint64, status *int, role *uint, passkey *string) error
		AdminGetUserStat(ctx context.Context, id uint64) (*entity.IamUserStat, error)
		AdminUpdateUserStat(ctx context.Context, id uint64, uploadedDiff *int64, downloadedDiff *int64) (int64, error)
		AddUserUploaded(ctx context.Context, userId uint64, amount uint64) error
		ReduceUserDownloaded(ctx context.Context, userId uint64, amount uint64) error
		ExtendUserVip(ctx context.Context, userId uint64, durationDays int, remark string) error
		GetUsersByIds(ctx context.Context, ids []uint64) ([]entity.IamUser, error)
		GetUserIdsByRoles(ctx context.Context, roleIds []uint) ([]uint64, error)
		GetUserProfilesByUserIds(ctx context.Context, userIds []uint64) ([]entity.IamUserProfile, error)
		QueryRankCandidates(ctx context.Context, page int, size int) ([]model.IamRankCandidate, error)
		UpdateUserRole(ctx context.Context, userId uint64, roleId uint) error
	}
	IIamUserUsecase interface {
		InvalidateUserCache(ctx context.Context, userId uint64)
		LoadActor(ctx context.Context, userId uint64) (*model.Actor, error)
		LoadActorByPasskey(ctx context.Context, passkey string) (*model.Actor, error)
		EnsureCanAuthenticate(ctx context.Context, user *entity.IamUser) error
		CheckPermission(ctx context.Context, actor *model.Actor, permKey string) (bool, error)
		Permissions(ctx context.Context, actor *model.Actor) (*iamout.UserPermissionListOut, error)
		LoginLogs(ctx context.Context, actor *model.Actor, in iamin.UserLoginLogListInp) (*iamout.UserLoginLogListOut, error)
		Create(ctx context.Context, in iamin.UserCreateInp) (uint64, error)
		Me(ctx context.Context, actor *model.Actor) (*iamout.UserMeOut, error)
		Get(ctx context.Context, actor *model.Actor, in iamin.UserGetInp) (*iamout.UserGetOut, error)
		UpdateProfile(ctx context.Context, actor *model.Actor, in iamin.UserProfileUpdateInp) error
		ChangePassword(ctx context.Context, actor *model.Actor, in iamin.UserPasswordChangeInp) error
		CreateEmailVerificationRequest(ctx context.Context, in iamin.EmailVerificationRequestCreateInp) error
		CreateEmailVerification(ctx context.Context, in iamin.EmailVerificationCreateInp) error
		CreatePasswordResetRequest(ctx context.Context, in iamin.PasswordResetRequestCreateInp) error
		CreatePasswordReset(ctx context.Context, in iamin.PasswordResetCreateInp) error
		ResetPasskey(ctx context.Context, actor *model.Actor) (string, error)
	}
)

var (
	localIamInviteDomain     IIamInviteDomain
	localIamInviteUsecase    IIamInviteUsecase
	localIamLoginLogDomain   IIamLoginLogDomain
	localIamPermissionDomain IIamPermissionDomain
	localIamRoleDomain       IIamRoleDomain
	localIamRoleUsecase      IIamRoleUsecase
	localIamSessionDomain    IIamSessionDomain
	localIamSessionUsecase   IIamSessionUsecase
	localIamTwoStepDomain    IIamTwoStepDomain
	localIamTwoStepUsecase   IIamTwoStepUsecase
	localIamUserDomain       IIamUserDomain
	localIamUserUsecase      IIamUserUsecase
)

func IamInviteDomain() IIamInviteDomain {
	if localIamInviteDomain == nil {
		panic("implement not found for interface IIamInviteDomain, forgot register?")
	}
	return localIamInviteDomain
}

func RegisterIamInviteDomain(i IIamInviteDomain) {
	localIamInviteDomain = i
}

func IamInviteUsecase() IIamInviteUsecase {
	if localIamInviteUsecase == nil {
		panic("implement not found for interface IIamInviteUsecase, forgot register?")
	}
	return localIamInviteUsecase
}

func RegisterIamInviteUsecase(i IIamInviteUsecase) {
	localIamInviteUsecase = i
}

func IamLoginLogDomain() IIamLoginLogDomain {
	if localIamLoginLogDomain == nil {
		panic("implement not found for interface IIamLoginLogDomain, forgot register?")
	}
	return localIamLoginLogDomain
}

func RegisterIamLoginLogDomain(i IIamLoginLogDomain) {
	localIamLoginLogDomain = i
}

func IamPermissionDomain() IIamPermissionDomain {
	if localIamPermissionDomain == nil {
		panic("implement not found for interface IIamPermissionDomain, forgot register?")
	}
	return localIamPermissionDomain
}

func RegisterIamPermissionDomain(i IIamPermissionDomain) {
	localIamPermissionDomain = i
}

func IamRoleDomain() IIamRoleDomain {
	if localIamRoleDomain == nil {
		panic("implement not found for interface IIamRoleDomain, forgot register?")
	}
	return localIamRoleDomain
}

func RegisterIamRoleDomain(i IIamRoleDomain) {
	localIamRoleDomain = i
}

func IamRoleUsecase() IIamRoleUsecase {
	if localIamRoleUsecase == nil {
		panic("implement not found for interface IIamRoleUsecase, forgot register?")
	}
	return localIamRoleUsecase
}

func RegisterIamRoleUsecase(i IIamRoleUsecase) {
	localIamRoleUsecase = i
}

func IamSessionDomain() IIamSessionDomain {
	if localIamSessionDomain == nil {
		panic("implement not found for interface IIamSessionDomain, forgot register?")
	}
	return localIamSessionDomain
}

func RegisterIamSessionDomain(i IIamSessionDomain) {
	localIamSessionDomain = i
}

func IamSessionUsecase() IIamSessionUsecase {
	if localIamSessionUsecase == nil {
		panic("implement not found for interface IIamSessionUsecase, forgot register?")
	}
	return localIamSessionUsecase
}

func RegisterIamSessionUsecase(i IIamSessionUsecase) {
	localIamSessionUsecase = i
}

func IamTwoStepDomain() IIamTwoStepDomain {
	if localIamTwoStepDomain == nil {
		panic("implement not found for interface IIamTwoStepDomain, forgot register?")
	}
	return localIamTwoStepDomain
}

func RegisterIamTwoStepDomain(i IIamTwoStepDomain) {
	localIamTwoStepDomain = i
}

func IamTwoStepUsecase() IIamTwoStepUsecase {
	if localIamTwoStepUsecase == nil {
		panic("implement not found for interface IIamTwoStepUsecase, forgot register?")
	}
	return localIamTwoStepUsecase
}

func RegisterIamTwoStepUsecase(i IIamTwoStepUsecase) {
	localIamTwoStepUsecase = i
}

func IamUserDomain() IIamUserDomain {
	if localIamUserDomain == nil {
		panic("implement not found for interface IIamUserDomain, forgot register?")
	}
	return localIamUserDomain
}

func RegisterIamUserDomain(i IIamUserDomain) {
	localIamUserDomain = i
}

func IamUserUsecase() IIamUserUsecase {
	if localIamUserUsecase == nil {
		panic("implement not found for interface IIamUserUsecase, forgot register?")
	}
	return localIamUserUsecase
}

func RegisterIamUserUsecase(i IIamUserUsecase) {
	localIamUserUsecase = i
}
