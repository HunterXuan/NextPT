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
		GetInviteByHashForUpdate(ctx context.Context, hash string) (*entity.IamInvite, error)
		GetInviteByHash(ctx context.Context, hash string) (*entity.IamInvite, error)
		GetInvitesByInviterIdAndHash(ctx context.Context, inviterId uint64, hash string) (*entity.IamInvite, error)
		QueryInvitesByInviter(ctx context.Context, inviterId uint64, page int, size int, status *uint) ([]entity.IamInvite, int, error)
		ExpireInvites(ctx context.Context, now *gtime.Time) (int64, error)
		UpdateInviteStatus(ctx context.Context, id uint64, status uint) error
		UpdateInvite(ctx context.Context, id uint64, data do.IamInvite) error
		AdminCreateInvites(ctx context.Context, invites []do.IamInvite) error
	}
	IIamInviteUsecase interface {
		List(ctx context.Context, actor *model.Actor, in iamin.InviteListInp) (*iamout.InviteListOut, error)
		Send(ctx context.Context, actor *model.Actor, in iamin.InviteSendInp) error
		Check(ctx context.Context, in iamin.InviteCheckInp) (*iamout.InviteCheckOut, error)
		CleanupExpired(ctx context.Context) (int64, error)
	}
	IIamPermissionDomain interface {
		CheckPermissionWithList(ctx context.Context, rolePerms []string, userAcls []string, permKey string) (bool, error)
		GrantUserPermission(ctx context.Context, userId uint64, permKey string, isDeny bool) error
		RevokeUserPermission(ctx context.Context, userId uint64, permKey string, isDeny bool) error
		GrantUserPermissionsBySource(ctx context.Context, userId uint64, permKeys []string, sourceType int, sourceId uint64, expireAt *gtime.Time) error
		DeactivateUserPermissionsBySource(ctx context.Context, userId uint64, sourceType int, sourceId uint64) error
		GetAllPermissions(ctx context.Context) []string
		GetUserPermissions(ctx context.Context, userId uint64) ([]entity.IamUserPermission, error)
	}
	IIamRoleDomain interface {
		GetRoleById(ctx context.Context, roleId uint) (*entity.IamRole, error)
		AdminListRoles(ctx context.Context) ([]entity.IamRole, error)
		AdminCreateRole(ctx context.Context, level int, nameI18N []byte, rules []byte, permissions []byte, isStaff bool) (uint, error)
		AdminUpdateRole(ctx context.Context, id uint, level *int, nameI18N []byte, rules []byte, permissions []byte, isStaff *bool) error
		AdminDeleteRole(ctx context.Context, id uint) (int, error)
	}
	IIamSessionDomain interface {
		GetGFToken() gtoken.Token
		GetGFMiddleware() gtoken.Middleware
		GenerateToken(ctx context.Context, userKey string, data any) (string, error)
		RemoveToken(ctx context.Context, userKey string) error
	}
	IIamSessionUsecase interface {
		Create(ctx context.Context, in iamin.SessionCreateInp) (*iamout.SessionCreateOut, error)
		Delete(ctx context.Context, actor *model.Actor) error
		VerifyPasskey(ctx context.Context, passkey string) (*model.Actor, error)
	}
	IIamUserDomain interface {
		GetUserByLogin(ctx context.Context, login string) (*entity.IamUser, error)
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
		UpdatePasswordHash(ctx context.Context, userId uint64, passwordHash string) error
		UpdatePasskey(ctx context.Context, userId uint64, passkey string) error
		// GetUserByPasskey 通过 Passkey 获取用户（无缓存，纯领域逻辑）
		GetUserByPasskey(ctx context.Context, passkey string) (*entity.IamUser, error)
		AdminListUsers(ctx context.Context, search string, order string, page int, size int) ([]*entity.IamUser, int, error)
		AdminUpdateUser(ctx context.Context, id uint64, status *int, role *uint, passkey *string) error
		AdminGetUserStat(ctx context.Context, id uint64) (*entity.IamUserStat, error)
		AdminUpdateUserStat(ctx context.Context, id uint64, uploadedDiff *int64, downloadedDiff *int64) (int64, error)
		GetUsersByIds(ctx context.Context, ids []uint64) ([]entity.IamUser, error)
		GetUserProfilesByUserIds(ctx context.Context, userIds []uint64) ([]entity.IamUserProfile, error)
	}
	IIamUserUsecase interface {
		InvalidateUserCache(ctx context.Context, userId uint64)
		LoadActor(ctx context.Context, userId uint64) (*model.Actor, error)
		LoadActorByPasskey(ctx context.Context, passkey string) (*model.Actor, error)
		EnsureCanAuthenticate(ctx context.Context, user *entity.IamUser) error
		CheckPermission(ctx context.Context, actor *model.Actor, permKey string) (bool, error)
		Create(ctx context.Context, in iamin.UserCreateInp) (uint64, error)
		Me(ctx context.Context, actor *model.Actor) (*iamout.UserMeOut, error)
		UpdateProfile(ctx context.Context, actor *model.Actor, in iamin.UserProfileUpdateInp) error
		ChangePassword(ctx context.Context, actor *model.Actor, in iamin.UserPasswordChangeInp) error
		ResetPasskey(ctx context.Context, actor *model.Actor) (string, error)
	}
)

var (
	localIamInviteDomain     IIamInviteDomain
	localIamInviteUsecase    IIamInviteUsecase
	localIamPermissionDomain IIamPermissionDomain
	localIamRoleDomain       IIamRoleDomain
	localIamSessionDomain    IIamSessionDomain
	localIamSessionUsecase   IIamSessionUsecase
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
