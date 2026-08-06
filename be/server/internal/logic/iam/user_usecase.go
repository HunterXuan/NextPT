package iam

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/iamin"
	"server/internal/model/out/iamout"
	"server/internal/service"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"golang.org/x/crypto/bcrypt"
)

type sIamUserUsecase struct{}

const temporaryTokenRateLimitScript = `
local current = redis.call("INCR", KEYS[1])
if current == 1 then
    redis.call("EXPIRE", KEYS[1], tonumber(ARGV[1]))
end
return current
`

const temporaryTokenStoreScript = `
redis.call("SET", KEYS[1], ARGV[1], "EX", tonumber(ARGV[3]))
redis.call("SET", KEYS[2], ARGV[2], "EX", tonumber(ARGV[3]))
return 1
`

const temporaryTokenConsumeScript = `
if redis.call("GET", KEYS[1]) ~= ARGV[1] then
    return 0
end
if redis.call("GET", KEYS[2]) ~= ARGV[2] then
    return 0
end
redis.call("DEL", KEYS[1], KEYS[2])
return 1
`

func init() {
	service.RegisterIamUserUsecase(NewIamUserUsecase())
}

func NewIamUserUsecase() *sIamUserUsecase {
	return &sIamUserUsecase{}
}

func (s *sIamUserUsecase) InvalidateUserCache(ctx context.Context, userId uint64) {
	// 1. Invalidate user ACLs cache
	aclsKey := service.SysCache().KeyIamUserAcls(ctx, userId)
	_, _ = gcache.Remove(ctx, aclsKey)
	_ = service.SysCache().PublishInvalidate(ctx, aclsKey)

	// 2. Invalidate actor details cache
	actorKey := service.SysCache().KeyIamActor(ctx, userId)
	_, _ = gcache.Remove(ctx, actorKey)
	_ = service.SysCache().PublishInvalidate(ctx, actorKey)

	// 3. Invalidate tracker passkey actor cache
	user, err := service.IamUserDomain().GetUserById(ctx, userId)
	if err == nil && user != nil && user.Passkey != "" {
		s.invalidatePasskeyActorCache(ctx, user.Passkey)
	}
}

func (s *sIamUserUsecase) LoadActor(ctx context.Context, userId uint64) (*model.Actor, error) {
	cacheKey := service.SysCache().KeyIamActor(ctx, userId)
	return s.loadActorWithCache(ctx, cacheKey, func(ctx context.Context) (*entity.IamUser, error) {
		return service.IamUserDomain().GetUserById(ctx, userId)
	})
}

func (s *sIamUserUsecase) LoadActorByPasskey(ctx context.Context, passkey string) (*model.Actor, error) {
	cacheKey := service.SysCache().KeyIamPasskeyActor(ctx, passkey)
	return s.loadActorWithCache(ctx, cacheKey, func(ctx context.Context) (*entity.IamUser, error) {
		user, err := service.IamUserDomain().GetUserByPasskey(ctx, passkey)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, gerror.New(gi18n.T(ctx, "iam.session.invalid_passkey"))
		}
		return user, nil
	})
}

func (s *sIamUserUsecase) loadActorWithCache(ctx context.Context, cacheKey string, loadUser func(context.Context) (*entity.IamUser, error)) (*model.Actor, error) {
	if actor, ok := s.getCachedActor(ctx, cacheKey); ok {
		roleVersion, err := s.getRoleActorVersion(ctx, actor.RoleId)
		if err == nil && actor.RoleVersion == roleVersion {
			return actor, nil
		}
		if err != nil {
			g.Log().Warningf(ctx, "iam actor cache: failed to read role actor version: %v", err)
		}
	}

	user, err := loadUser(ctx)
	if err != nil {
		return nil, err
	}

	actor, err := s.buildActor(ctx, user)
	if err != nil {
		return nil, err
	}

	roleVersion, err := s.getRoleActorVersion(ctx, actor.RoleId)
	if err != nil {
		g.Log().Warningf(ctx, "iam actor cache: failed to read role actor version: %v", err)
		return actor, nil
	}
	actor.RoleVersion = roleVersion
	_ = gcache.Set(ctx, cacheKey, actor, 24*time.Hour)
	return actor, nil
}

func (s *sIamUserUsecase) getCachedActor(ctx context.Context, cacheKey string) (*model.Actor, bool) {
	v, err := gcache.Get(ctx, cacheKey)
	if err != nil || v == nil || v.IsNil() {
		return nil, false
	}

	var actor model.Actor
	if err := v.Scan(&actor); err != nil || actor.Id == 0 {
		return nil, false
	}
	return &actor, true
}

func (s *sIamUserUsecase) buildActor(ctx context.Context, user *entity.IamUser) (*model.Actor, error) {
	if user == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.user.not_found"))
	}
	if err := s.EnsureCanAuthenticate(ctx, user); err != nil {
		return nil, err
	}

	role, err := service.IamRoleDomain().GetRoleById(ctx, user.Role)
	if err != nil || role == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.session.role_missing"))
	}

	return &model.Actor{
		Id:        user.Id,
		Email:     user.Email,
		RoleId:    user.Role,
		RoleLevel: role.Level,
		IsStaff:   role.IsStaff,
	}, nil
}

func (s *sIamUserUsecase) EnsureCanAuthenticate(ctx context.Context, user *entity.IamUser) error {
	if user == nil {
		return gerror.New(gi18n.T(ctx, "iam.user.not_found"))
	}
	if user.Status == consts.IamUserStatusPending {
		return gerror.New(gi18n.T(ctx, "iam.session.email_unverified"))
	}
	if user.Status != consts.IamUserStatusConfirmed {
		return gerror.New(gi18n.T(ctx, "iam.session.account_banned"))
	}
	banned, err := service.ModUserDomain().HasActiveMod(ctx, user.Id, []int{consts.ModUserTypeBanned})
	if err != nil {
		return err
	}
	if banned {
		return gerror.New(gi18n.T(ctx, "iam.session.account_banned"))
	}
	return nil
}

func (s *sIamUserUsecase) getRoleActorVersion(ctx context.Context, roleId uint) (int64, error) {
	v, err := g.Redis().Do(ctx, "GET", service.SysCache().KeyIamRoleActorVersion(ctx, roleId))
	if err != nil {
		return 0, err
	}
	if v == nil || v.IsNil() {
		return 0, nil
	}
	return v.Int64(), nil
}

func (s *sIamUserUsecase) CheckPermission(ctx context.Context, actor *model.Actor, permKey string) (bool, error) {
	rolePerms, userAcls, err := s.loadPermissionLists(ctx, actor)
	if err != nil {
		return false, err
	}

	return service.IamPermissionDomain().CheckPermissionWithList(ctx, rolePerms, userAcls, permKey)
}

func (s *sIamUserUsecase) Permissions(ctx context.Context, actor *model.Actor) (*iamout.UserPermissionListOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	rolePerms, userAcls, err := s.loadPermissionLists(ctx, actor)
	if err != nil {
		return nil, err
	}

	permissions := make([]string, 0)
	for _, permission := range service.IamPermissionDomain().GetAllPermissions(ctx) {
		permission = strings.TrimSpace(permission)
		if permission == "" || permission == consts.IamPermissionAll || !strings.HasSuffix(permission, ":*") {
			continue
		}
		ok, err := service.IamPermissionDomain().CheckPermissionWithList(ctx, rolePerms, userAcls, permission)
		if err != nil {
			return nil, err
		}
		if ok {
			permissions = append(permissions, permission)
		}
	}

	return &iamout.UserPermissionListOut{Permissions: permissions}, nil
}

func (s *sIamUserUsecase) LoginLogs(ctx context.Context, actor *model.Actor, in iamin.UserLoginLogListInp) (*iamout.UserLoginLogListOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	logs, total, err := service.IamLoginLogDomain().List(ctx, model.IamLoginLogListOptions{
		UserId: actor.Id,
		Result: in.Result,
		Page:   in.Page,
		Size:   in.Size,
	})
	if err != nil {
		return nil, err
	}
	items := make([]model.IamLoginLogItem, 0, len(logs))
	for _, log := range logs {
		items = append(items, model.NewIamLoginLogItem(log))
	}
	return &iamout.UserLoginLogListOut{
		List:  items,
		Total: total,
		Page:  in.Page,
		Size:  in.Size,
	}, nil
}

func (s *sIamUserUsecase) loadPermissionLists(ctx context.Context, actor *model.Actor) ([]string, []string, error) {
	if actor == nil {
		return nil, nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	// 1. Get role perms from cache or DB
	rolePermsKey := service.SysCache().KeyIamRolePerms(ctx, actor.RoleId)
	rolePermsVar, err := gcache.GetOrSetFunc(ctx, rolePermsKey, func(ctx context.Context) (value interface{}, err error) {
		role, err := service.IamRoleDomain().GetRoleById(ctx, actor.RoleId)
		if err != nil || role == nil {
			return []string{}, err
		}
		var perms []string
		if role.Permissions != nil {
			_ = role.Permissions.Scan(&perms)
		}
		return perms, nil
	}, 24*time.Hour)
	if err != nil {
		return nil, nil, err
	}
	var rolePerms []string
	_ = rolePermsVar.Scan(&rolePerms)

	// 2. Get user ACLs from cache or DB
	userAclsKey := service.SysCache().KeyIamUserAcls(ctx, actor.Id)
	userAclsVar, err := gcache.GetOrSetFunc(ctx, userAclsKey, func(ctx context.Context) (value interface{}, err error) {
		acls, err := service.IamPermissionDomain().GetUserPermissions(ctx, actor.Id)
		if err != nil {
			return []string{}, err
		}
		var list []string
		for _, acl := range acls {
			list = append(list, acl.PermKey)
		}
		return list, nil
	}, time.Minute)
	if err != nil {
		return nil, nil, err
	}
	var userAcls []string
	_ = userAclsVar.Scan(&userAcls)

	return rolePerms, userAcls, nil
}

func (s *sIamUserUsecase) Create(ctx context.Context, in iamin.UserCreateInp) (uint64, error) {
	in.Username = strings.TrimSpace(in.Username)
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.InviteHash = strings.TrimSpace(in.InviteHash)

	registerEnabled := s.getIamConfigCache(ctx, consts.SiteConfigIamRegisterEnabled).Bool()
	if !registerEnabled && in.InviteHash == "" {
		return 0, gerror.New(gi18n.T(ctx, "iam.invite.register_disabled"))
	}

	exists, err := service.IamUserDomain().CheckUsernameExists(ctx, in.Username)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, gerror.New(gi18n.T(ctx, "iam.user.username_exists"))
	}

	exists, err = service.IamUserDomain().CheckEmailExists(ctx, in.Email)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, gerror.New(gi18n.T(ctx, "iam.user.email_exists"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	passkey := gmd5.MustEncryptString(grand.S(32))

	var newUserId uint64

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var invitedById uint64
		if in.InviteHash != "" {
			invite, err := service.IamInviteDomain().GetInviteByHashForUpdate(ctx, in.InviteHash)
			if err != nil {
				return err
			}
			if invite == nil {
				return gerror.New(gi18n.T(ctx, "iam.invite.not_found"))
			}
			if invite.Status != consts.IamInviteStatusSent && invite.Status != consts.IamInviteStatusUnused {
				return gerror.New(gi18n.T(ctx, "iam.invite.invalid_status"))
			}
			if invite.ExpireAt != nil && !invite.ExpireAt.After(gtime.Now()) {
				return gerror.New(gi18n.T(ctx, "iam.invite.expired"))
			}
			if invite.Status == consts.IamInviteStatusSent && invite.InviteeEmail != "" && invite.InviteeEmail != in.Email {
				return gerror.New(gi18n.T(ctx, "iam.invite.email_mismatch"))
			}
			invitedById = invite.InviterId
		}

		userId, err := service.IamUserDomain().InsertUser(ctx, do.IamUser{
			Username:     in.Username,
			Email:        in.Email,
			PasswordHash: string(hash),
			Passkey:      passkey,
			Status:       consts.IamUserStatusPending,
			Role:         s.getIamConfigCache(ctx, consts.SiteConfigIamDefaultRegisterRole).Uint(),
			InvitedBy:    invitedById,
		})
		if err != nil {
			return err
		}
		newUserId = userId

		if err := service.IamUserDomain().InsertUserProfile(ctx, do.IamUserProfile{
			UserId: userId,
			Avatar: "",
			Info:   "",
		}); err != nil {
			return err
		}

		if err := service.IamUserDomain().InsertUserStat(ctx, do.IamUserStat{
			UserId:     userId,
			Uploaded:   0,
			Downloaded: 0,
			Bonus:      0,
		}); err != nil {
			return err
		}

		if in.InviteHash != "" {
			invite, _ := service.IamInviteDomain().GetInviteByHashForUpdate(ctx, in.InviteHash)
			if invite != nil {
				_ = service.IamInviteDomain().UpdateInvite(ctx, invite.Id, do.IamInvite{
					Status:    consts.IamInviteStatusUsed,
					InviteeId: userId,
					UsedAt:    gtime.Now(),
				})
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	if err := s.CreateEmailVerificationRequest(ctx, iamin.EmailVerificationRequestCreateInp{Email: in.Email}); err != nil {
		glog.Warningf(ctx, "send registration email verification failed: error=%v", err)
	}
	return newUserId, nil
}

func (s *sIamUserUsecase) Me(ctx context.Context, actor *model.Actor) (*iamout.UserMeOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	user, err := service.IamUserDomain().GetUserById(ctx, actor.Id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.user.not_found"))
	}

	profile, err := service.IamUserDomain().GetUserProfile(ctx, actor.Id)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		profile = &entity.IamUserProfile{}
	}

	stat, err := service.IamUserDomain().GetUserStat(ctx, actor.Id)
	if err != nil {
		return nil, err
	}
	if stat == nil {
		stat = &entity.IamUserStat{}
	}

	role, err := service.IamRoleDomain().GetRoleById(ctx, user.Role)
	if err != nil || role == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.session.role_missing"))
	}

	return &iamout.UserMeOut{
		User: iamout.UserMeAccountOut{
			Id:             user.Id,
			Username:       user.Username,
			Email:          user.Email,
			Passkey:        user.Passkey,
			Status:         user.Status,
			TwoStepEnabled: user.TwoStepType == consts.IamTwoStepTypeTOTP,
			VipUntil:       user.VipUntil,
			CreatedAt:      user.CreatedAt,
		},
		Role: iamout.UserMeRoleOut{
			Id:      user.Role,
			Name:    s.localizeRoleName(ctx, role),
			Level:   actor.RoleLevel,
			IsStaff: actor.IsStaff,
		},
		Profile: iamout.UserMeProfileOut{
			Avatar:    profile.Avatar,
			Info:      profile.Info,
			Signature: profile.Signature,
		},
		Stat: iamout.UserMeStatOut{
			Uploaded:      stat.Uploaded,
			Downloaded:    stat.Downloaded,
			RawUploaded:   stat.RawUploaded,
			RawDownloaded: stat.RawDownloaded,
			Bonus:         stat.Bonus,
			ShareRatio:    s.calculateShareRatio(stat.Uploaded, stat.Downloaded),
		},
	}, nil
}

func (s *sIamUserUsecase) Get(ctx context.Context, actor *model.Actor, in iamin.UserGetInp) (*iamout.UserGetOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	language := gi18n.LanguageFromCtx(ctx)
	if language == "" {
		language = g.Cfg().MustGet(ctx, "i18n.default", "zh-CN").String()
	}
	cacheKey := service.SysCache().KeyIamUserPublic(ctx, in.Id, language)
	value, err := gcache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		return s.buildPublicUser(ctx, in.Id)
	}, 5*time.Minute)
	if err != nil {
		return nil, err
	}
	if out, ok := value.Val().(*iamout.UserGetOut); ok {
		return out, nil
	}

	var out iamout.UserGetOut
	if err := value.Scan(&out); err != nil {
		return s.buildPublicUser(ctx, in.Id)
	}
	return &out, nil
}

func (s *sIamUserUsecase) buildPublicUser(ctx context.Context, userId uint64) (*iamout.UserGetOut, error) {
	user, err := service.IamUserDomain().GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.user.not_found"))
	}

	profile, err := service.IamUserDomain().GetUserProfile(ctx, userId)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		profile = &entity.IamUserProfile{}
	}

	role, err := service.IamRoleDomain().GetRoleById(ctx, user.Role)
	if err != nil || role == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.session.role_missing"))
	}

	stat, err := service.AccountingTrafficDomain().GetUserStat(ctx, userId)
	if err != nil {
		return nil, err
	}
	if stat == nil {
		stat = &entity.IamUserStat{}
	}

	return &iamout.UserGetOut{
		User: iamout.UserGetAccountOut{
			Id:        user.Id,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		Role: iamout.UserGetRoleOut{
			Id:      user.Role,
			Name:    s.localizeRoleName(ctx, role),
			IsStaff: role.IsStaff,
		},
		Profile: iamout.UserGetProfileOut{
			Avatar:    profile.Avatar,
			Info:      profile.Info,
			Signature: profile.Signature,
		},
		Stat: iamout.UserGetStatOut{
			Uploaded:   stat.Uploaded,
			Downloaded: stat.Downloaded,
			ShareRatio: s.calculateShareRatio(stat.Uploaded, stat.Downloaded),
			SeedTime:   stat.SeedTime,
		},
	}, nil
}

func (s *sIamUserUsecase) localizeRoleName(ctx context.Context, role *entity.IamRole) string {
	if role == nil || role.NameI18N == nil {
		return ""
	}

	var names map[string]string
	if err := role.NameI18N.Scan(&names); err != nil || len(names) == 0 {
		return ""
	}

	lang := gi18n.LanguageFromCtx(ctx)
	if lang != "" && names[lang] != "" {
		return names[lang]
	}

	defaultLang := g.Cfg().MustGet(ctx, "i18n.default", "zh-CN").String()
	if defaultLang != "" && names[defaultLang] != "" {
		return names[defaultLang]
	}

	for _, name := range names {
		if name != "" {
			return name
		}
	}
	return ""
}

func (s *sIamUserUsecase) UpdateProfile(ctx context.Context, actor *model.Actor, in iamin.UserProfileUpdateInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	return service.IamUserDomain().UpdateUserProfile(ctx, actor.Id, in.Avatar, in.Info, in.Signature)
}

func (s *sIamUserUsecase) ChangePassword(ctx context.Context, actor *model.Actor, in iamin.UserPasswordChangeInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	passwordHash, err := service.IamUserDomain().GetUserPasswordHash(ctx, actor.Id)
	if err != nil {
		return err
	}
	if passwordHash == "" {
		return gerror.New(gi18n.T(ctx, "iam.user.not_found"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(in.OldPassword)); err != nil {
		return gerror.New(gi18n.T(ctx, "iam.user.old_password_invalid"))
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(in.NewPassword)); err == nil {
		return gerror.New(gi18n.T(ctx, "iam.user.new_password_same"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := service.IamUserDomain().UpdatePasswordHash(ctx, actor.Id, string(hash)); err != nil {
		return err
	}
	return service.IamSessionDomain().RemoveByUser(ctx, actor.Id)
}

func (s *sIamUserUsecase) CreateEmailVerificationRequest(ctx context.Context, in iamin.EmailVerificationRequestCreateInp) error {
	email := strings.ToLower(strings.TrimSpace(in.Email))

	ipAllowed, err := s.temporaryTokenRateAllowed(
		ctx,
		service.SysCache().KeyIamEmailVerificationRateIp(ctx, s.temporaryTokenDigest(s.requestIp(ctx))),
		consts.IamEmailVerificationRateWindow,
		consts.IamEmailVerificationRateLimitByIp,
	)
	if err != nil {
		return err
	}
	emailAllowed, err := s.temporaryTokenRateAllowed(
		ctx,
		service.SysCache().KeyIamEmailVerificationRateEmail(ctx, s.temporaryTokenDigest(email)),
		consts.IamEmailVerificationRateWindow,
		consts.IamEmailVerificationRateLimitEmail,
	)
	if err != nil {
		return err
	}
	if !ipAllowed || !emailAllowed {
		return nil
	}

	user, err := service.IamUserDomain().GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil || user.Status != consts.IamUserStatusPending {
		return nil
	}

	token, err := s.newTemporaryToken()
	if err != nil {
		return err
	}
	tokenHash := s.temporaryTokenDigest(token)
	if err := s.storeTemporaryToken(
		ctx,
		service.SysCache().KeyIamEmailVerificationToken(ctx, tokenHash),
		service.SysCache().KeyIamEmailVerificationUser(ctx, user.Id),
		user.Id,
		tokenHash,
		consts.IamEmailVerificationTokenTTL,
	); err != nil {
		return err
	}

	s.sendEmailVerificationMail(user.Username, user.Email, token)
	return nil
}

func (s *sIamUserUsecase) CreateEmailVerification(ctx context.Context, in iamin.EmailVerificationCreateInp) error {
	tokenHash := s.temporaryTokenDigest(strings.TrimSpace(in.Token))
	tokenKey := service.SysCache().KeyIamEmailVerificationToken(ctx, tokenHash)
	userId, err := s.temporaryTokenUserId(ctx, tokenKey)
	if err != nil {
		return err
	}
	if userId == 0 {
		return s.emailVerificationTokenError(ctx)
	}

	user, err := service.IamUserDomain().GetUserById(ctx, userId)
	if err != nil {
		return err
	}
	if user == nil || user.Status != consts.IamUserStatusPending {
		return s.emailVerificationTokenError(ctx)
	}

	consumed, err := s.consumeTemporaryToken(
		ctx,
		tokenKey,
		service.SysCache().KeyIamEmailVerificationUser(ctx, userId),
		userId,
		tokenHash,
	)
	if err != nil {
		return err
	}
	if !consumed {
		return s.emailVerificationTokenError(ctx)
	}

	confirmed, err := service.IamUserDomain().ConfirmUserEmail(ctx, userId)
	if err != nil {
		return err
	}
	if !confirmed {
		return s.emailVerificationTokenError(ctx)
	}
	s.InvalidateUserCache(ctx, userId)
	return nil
}

func (s *sIamUserUsecase) CreatePasswordResetRequest(ctx context.Context, in iamin.PasswordResetRequestCreateInp) error {
	email := strings.ToLower(strings.TrimSpace(in.Email))

	ipAllowed, err := s.temporaryTokenRateAllowed(
		ctx,
		service.SysCache().KeyIamPasswordResetRateIp(ctx, s.temporaryTokenDigest(s.requestIp(ctx))),
		consts.IamPasswordResetRateWindow,
		consts.IamPasswordResetRateLimitByIp,
	)
	if err != nil {
		return err
	}
	emailAllowed, err := s.temporaryTokenRateAllowed(
		ctx,
		service.SysCache().KeyIamPasswordResetRateEmail(ctx, s.temporaryTokenDigest(email)),
		consts.IamPasswordResetRateWindow,
		consts.IamPasswordResetRateLimitEmail,
	)
	if err != nil {
		return err
	}
	if !ipAllowed || !emailAllowed {
		return nil
	}

	user, err := service.IamUserDomain().GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	token, err := s.newTemporaryToken()
	if err != nil {
		return err
	}
	tokenHash := s.temporaryTokenDigest(token)
	if err := s.storeTemporaryToken(
		ctx,
		service.SysCache().KeyIamPasswordResetToken(ctx, tokenHash),
		service.SysCache().KeyIamPasswordResetUser(ctx, user.Id),
		user.Id,
		tokenHash,
		consts.IamPasswordResetTokenTTL,
	); err != nil {
		return err
	}

	s.sendPasswordResetMail(user.Username, user.Email, token)
	return nil
}

func (s *sIamUserUsecase) CreatePasswordReset(ctx context.Context, in iamin.PasswordResetCreateInp) error {
	tokenHash := s.temporaryTokenDigest(strings.TrimSpace(in.Token))
	tokenKey := service.SysCache().KeyIamPasswordResetToken(ctx, tokenHash)
	userId, err := s.temporaryTokenUserId(ctx, tokenKey)
	if err != nil {
		return err
	}
	if userId == 0 {
		return s.passwordResetTokenError(ctx)
	}

	passwordHash, err := service.IamUserDomain().GetUserPasswordHash(ctx, userId)
	if err != nil {
		return err
	}
	if passwordHash == "" {
		return s.passwordResetTokenError(ctx)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(in.NewPassword)); err == nil {
		return gerror.New(gi18n.T(ctx, "iam.user.new_password_same"))
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	consumed, err := s.consumeTemporaryToken(
		ctx,
		tokenKey,
		service.SysCache().KeyIamPasswordResetUser(ctx, userId),
		userId,
		tokenHash,
	)
	if err != nil {
		return err
	}
	if !consumed {
		return s.passwordResetTokenError(ctx)
	}

	if err := service.IamUserDomain().UpdatePasswordHash(ctx, userId, string(newPasswordHash)); err != nil {
		return err
	}
	if err := service.IamSessionDomain().RemoveByUser(ctx, userId); err != nil {
		return err
	}
	s.InvalidateUserCache(ctx, userId)
	return nil
}

func (s *sIamUserUsecase) temporaryTokenRateAllowed(ctx context.Context, key string, window time.Duration, limit int) (bool, error) {
	value, err := g.Redis().Do(
		ctx,
		"EVAL",
		temporaryTokenRateLimitScript,
		1,
		key,
		int(window.Seconds()),
	)
	if err != nil {
		return false, err
	}
	return value.Int() <= limit, nil
}

func (s *sIamUserUsecase) storeTemporaryToken(ctx context.Context, tokenKey string, userKey string, userId uint64, tokenHash string, ttl time.Duration) error {
	_, err := g.Redis().Do(
		ctx,
		"EVAL",
		temporaryTokenStoreScript,
		2,
		tokenKey,
		userKey,
		gconv.String(userId),
		tokenHash,
		int(ttl.Seconds()),
	)
	return err
}

func (s *sIamUserUsecase) consumeTemporaryToken(ctx context.Context, tokenKey string, userKey string, userId uint64, tokenHash string) (bool, error) {
	value, err := g.Redis().Do(
		ctx,
		"EVAL",
		temporaryTokenConsumeScript,
		2,
		tokenKey,
		userKey,
		gconv.String(userId),
		tokenHash,
	)
	if err != nil {
		return false, err
	}
	return value.Int() == 1, nil
}

func (s *sIamUserUsecase) temporaryTokenUserId(ctx context.Context, tokenKey string) (uint64, error) {
	value, err := g.Redis().Do(ctx, "GET", tokenKey)
	if err != nil {
		return 0, err
	}
	if value == nil || value.IsNil() {
		return 0, nil
	}
	return value.Uint64(), nil
}

func (s *sIamUserUsecase) newTemporaryToken() (string, error) {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func (s *sIamUserUsecase) temporaryTokenDigest(value string) string {
	digest := sha256.Sum256([]byte(value))
	return hex.EncodeToString(digest[:])
}

func (s *sIamUserUsecase) requestIp(ctx context.Context) string {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil || strings.TrimSpace(r.GetClientIp()) == "" {
		return "unknown"
	}
	return strings.TrimSpace(r.GetClientIp())
}

func (s *sIamUserUsecase) emailVerificationTokenError(ctx context.Context) error {
	return gerror.New(gi18n.T(ctx, "iam.user.email_verification_token_invalid"))
}

func (s *sIamUserUsecase) passwordResetTokenError(ctx context.Context) error {
	return gerror.New(gi18n.T(ctx, "iam.user.password_reset_token_invalid"))
}

func (s *sIamUserUsecase) sendEmailVerificationMail(username, recipient, token string) {
	mailCtx := gctx.New()
	defaultLang := g.Cfg().MustGet(mailCtx, "i18n.default", "zh-CN").String()
	i18nCtx := gi18n.WithLanguage(mailCtx, defaultLang)
	siteName := strings.TrimSpace(g.Cfg().MustGet(mailCtx, "site.name", "NextPT").String())
	siteUrl := strings.TrimRight(strings.TrimSpace(g.Cfg().MustGet(mailCtx, "site.url", "http://localhost:3000").String()), "/")
	if siteName == "" {
		siteName = "NextPT"
	}
	if siteUrl == "" {
		siteUrl = "http://localhost:3000"
	}

	verificationUrl := siteUrl + "/verify-email?token=" + url.QueryEscape(token)
	subject := fmt.Sprintf(gi18n.T(i18nCtx, "iam.user.email_verification_mail_subject"), siteName)
	greeting := fmt.Sprintf(gi18n.T(i18nCtx, "iam.user.email_verification_mail_greeting"), username)
	intro := fmt.Sprintf(gi18n.T(i18nCtx, "iam.user.email_verification_mail_intro"), siteName)
	action := gi18n.T(i18nCtx, "iam.user.email_verification_mail_action")
	expiry := fmt.Sprintf(gi18n.T(i18nCtx, "iam.user.email_verification_mail_expiry"), int(consts.IamEmailVerificationTokenTTL.Hours()))
	note := gi18n.T(i18nCtx, "iam.user.email_verification_mail_ignore")
	s.sendAccountActionMail(mailCtx, model.AccountActionMail{
		Kind:      "email verification",
		Recipient: recipient,
		Subject:   subject,
		Greeting:  greeting,
		Intro:     intro,
		ActionURL: verificationUrl,
		Action:    action,
		Expiry:    expiry,
		Note:      note,
	})
}

func (s *sIamUserUsecase) sendPasswordResetMail(username, recipient, token string) {
	mailCtx := gctx.New()
	defaultLang := g.Cfg().MustGet(mailCtx, "i18n.default", "zh-CN").String()
	i18nCtx := gi18n.WithLanguage(mailCtx, defaultLang)
	siteName := strings.TrimSpace(g.Cfg().MustGet(mailCtx, "site.name", "NextPT").String())
	siteUrl := strings.TrimRight(strings.TrimSpace(g.Cfg().MustGet(mailCtx, "site.url", "http://localhost:3000").String()), "/")
	if siteName == "" {
		siteName = "NextPT"
	}
	if siteUrl == "" {
		siteUrl = "http://localhost:3000"
	}

	resetUrl := siteUrl + "/reset-password?token=" + url.QueryEscape(token)
	subject := fmt.Sprintf(gi18n.T(i18nCtx, "iam.user.password_reset_mail_subject"), siteName)
	greeting := fmt.Sprintf(gi18n.T(i18nCtx, "iam.user.password_reset_mail_greeting"), username)
	intro := fmt.Sprintf(gi18n.T(i18nCtx, "iam.user.password_reset_mail_intro"), siteName)
	action := gi18n.T(i18nCtx, "iam.user.password_reset_mail_action")
	expiry := fmt.Sprintf(gi18n.T(i18nCtx, "iam.user.password_reset_mail_expiry"), int(consts.IamPasswordResetTokenTTL.Minutes()))
	note := gi18n.T(i18nCtx, "iam.user.password_reset_mail_ignore")
	s.sendAccountActionMail(mailCtx, model.AccountActionMail{
		Kind:      "password reset",
		Recipient: recipient,
		Subject:   subject,
		Greeting:  greeting,
		Intro:     intro,
		ActionURL: resetUrl,
		Action:    action,
		Expiry:    expiry,
		Note:      note,
	})
}

func (s *sIamUserUsecase) sendAccountActionMail(ctx context.Context, mail model.AccountActionMail) {
	textBody := strings.Join([]string{mail.Greeting, mail.Intro, mail.Action, mail.ActionURL, mail.Expiry, mail.Note}, "\n\n")
	htmlBody := fmt.Sprintf(
		`<!doctype html><html><body style="margin:0;background:#f8fafc;font-family:Arial,sans-serif;color:#0f172a"><div style="max-width:560px;margin:0 auto;padding:32px 20px"><div style="background:#fff;border:1px solid #e2e8f0;border-radius:8px;padding:28px"><h1 style="margin:0 0 20px;font-size:20px">%s</h1><p style="margin:0 0 12px;line-height:1.7">%s</p><p style="margin:0 0 20px;line-height:1.7">%s</p><p style="margin:0 0 20px"><a href="%s" style="display:inline-block;border-radius:6px;background:#0284c7;padding:11px 18px;color:#fff;text-decoration:none;font-weight:600">%s</a></p><p style="margin:0 0 8px;color:#475569;line-height:1.7">%s</p><p style="margin:0;color:#64748b;line-height:1.7">%s</p></div></div></body></html>`,
		html.EscapeString(mail.Subject),
		html.EscapeString(mail.Greeting),
		html.EscapeString(mail.Intro),
		html.EscapeString(mail.ActionURL),
		html.EscapeString(mail.Action),
		html.EscapeString(mail.Expiry),
		html.EscapeString(mail.Note),
	)

	go func() {
		if err := service.SysMailgun().SendHtmlMail(ctx, mail.Subject, textBody, htmlBody, mail.Recipient); err != nil {
			glog.Warningf(ctx, "send %s mail failed: error=%v", mail.Kind, err)
		}
	}()
}

func (s *sIamUserUsecase) ResetPasskey(ctx context.Context, actor *model.Actor) (string, error) {
	if actor == nil {
		return "", gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	user, err := service.IamUserDomain().GetUserById(ctx, actor.Id)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", gerror.New(gi18n.T(ctx, "iam.user.not_found"))
	}

	oldPasskey := user.Passkey
	newPasskey := gmd5.MustEncryptString(grand.S(32))
	if err := service.IamUserDomain().UpdatePasskey(ctx, actor.Id, newPasskey); err != nil {
		return "", err
	}

	s.invalidatePasskeyActorCache(ctx, oldPasskey)
	s.InvalidateUserCache(ctx, actor.Id)
	return newPasskey, nil
}

func (s *sIamUserUsecase) invalidatePasskeyActorCache(ctx context.Context, passkey string) {
	if passkey == "" {
		return
	}
	passkeyKey := service.SysCache().KeyIamPasskeyActor(ctx, passkey)
	_, _ = gcache.Remove(ctx, passkeyKey)
	_ = service.SysCache().PublishInvalidate(ctx, passkeyKey)
}

func (s *sIamUserUsecase) calculateShareRatio(uploaded, downloaded uint64) float64 {
	if downloaded == 0 {
		return 0
	}
	return float64(uploaded) / float64(downloaded)
}

func (s *sIamUserUsecase) getIamConfigCache(ctx context.Context, key string) *gvar.Var {
	cacheKey := service.SysCache().KeySiteConfigFullPath(ctx, key)
	val, err := gcache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		return service.SiteConfigDomain().GetByPath(ctx, key).Val(), nil
	}, 5*time.Minute)
	if err != nil || val.IsNil() {
		return service.SiteConfigDomain().GetByPath(ctx, key)
	}
	return gvar.New(val.Val())
}
