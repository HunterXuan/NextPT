package iam

import (
	"context"
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
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"golang.org/x/crypto/bcrypt"
)

type sIamUserUsecase struct{}

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
		return false, err
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
		return false, err
	}
	var userAcls []string
	_ = userAclsVar.Scan(&userAcls)

	return service.IamPermissionDomain().CheckPermissionWithList(ctx, rolePerms, userAcls, permKey)
}

func (s *sIamUserUsecase) Create(ctx context.Context, in iamin.UserCreateInp) (uint64, error) {
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
			Status:       consts.IamUserStatusConfirmed,
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
		Id:         user.Id,
		Username:   user.Username,
		Email:      user.Email,
		Passkey:    user.Passkey,
		Role:       user.Role,
		RoleName:   s.localizeRoleName(ctx, role),
		RoleLevel:  actor.RoleLevel,
		IsStaff:    actor.IsStaff,
		Status:     user.Status,
		VipUntil:   user.VipUntil,
		Avatar:     profile.Avatar,
		Info:       profile.Info,
		Signature:  profile.Signature,
		Uploaded:   stat.Uploaded,
		Downloaded: stat.Downloaded,
		Bonus:      stat.Bonus,
		ShareRatio: s.calculateShareRatio(stat.Uploaded, stat.Downloaded),
		CreatedAt:  user.CreatedAt,
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
	return service.IamSessionDomain().RemoveToken(ctx, gconv.String(actor.Id))
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
