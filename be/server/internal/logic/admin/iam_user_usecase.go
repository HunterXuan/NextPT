package admin

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/modin"
	"server/internal/model/out/adminout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gcache"
)

type sAdminIamUserUsecase struct{}

func NewAdminIamUserUsecase() *sAdminIamUserUsecase {
	return &sAdminIamUserUsecase{}
}

func init() {
	service.RegisterAdminIamUserUsecase(NewAdminIamUserUsecase())
}

func (s *sAdminIamUserUsecase) List(ctx context.Context, actor *model.Actor, in adminin.IamUserListInp) (*adminout.IamUserListOut, error) {
	users, total, err := service.IamUserDomain().AdminListUsers(ctx, in.Search, in.Order, in.Page, in.Size)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "admin.user.fetch_failed"))
	}
	return &adminout.IamUserListOut{
		Users: users,
		Total: total,
	}, nil
}

func (s *sAdminIamUserUsecase) Update(ctx context.Context, actor *model.Actor, in adminin.IamUserUpdateInp) error {
	if in.Status != nil && !s.isValidUserStatus(*in.Status) {
		return gerror.New("invalid user status")
	}

	var oldPasskey string
	if in.Passkey != nil {
		oldUser, err := service.IamUserDomain().GetUserById(ctx, in.Id)
		if err != nil {
			return gerror.Wrap(err, gi18n.T(ctx, "admin.user.update_failed"))
		}
		if oldUser != nil {
			oldPasskey = oldUser.Passkey
		}
	}

	err := service.IamUserDomain().AdminUpdateUser(ctx, in.Id, in.Status, in.Role, in.Passkey)
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "admin.user.update_failed"))
	}

	// If role, status, or passkey was updated, invalidate cache so changes take effect immediately.
	if in.Role != nil || in.Status != nil || in.Passkey != nil {
		service.IamUserUsecase().InvalidateUserCache(ctx, in.Id)
	}
	if in.Passkey != nil && oldPasskey != "" && *in.Passkey != oldPasskey {
		s.invalidatePasskeyActorCache(ctx, oldPasskey)
	}

	return nil
}

func (s *sAdminIamUserUsecase) isValidUserStatus(status int) bool {
	switch status {
	case consts.IamUserStatusPending,
		consts.IamUserStatusConfirmed,
		consts.IamUserStatusDisabled:
		return true
	default:
		return false
	}
}

func (s *sAdminIamUserUsecase) invalidatePasskeyActorCache(ctx context.Context, passkey string) {
	if passkey == "" {
		return
	}
	key := service.SysCache().KeyIamPasskeyActor(ctx, passkey)
	_, _ = gcache.Remove(ctx, key)
	_ = service.SysCache().PublishInvalidate(ctx, key)
}

func (s *sAdminIamUserUsecase) StatDetail(ctx context.Context, actor *model.Actor, in adminin.IamUserStatDetailInp) (*adminout.IamUserStatDetailOut, error) {
	stat, err := service.IamUserDomain().AdminGetUserStat(ctx, in.Id)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "admin.user.fetch_stat_failed"))
	}
	if stat == nil || stat.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "admin.user.stat_not_found"))
	}
	return &adminout.IamUserStatDetailOut{IamUserStat: *stat}, nil
}

func (s *sAdminIamUserUsecase) StatUpdate(ctx context.Context, actor *model.Actor, in adminin.IamUserStatUpdateInp) error {
	rows, err := service.IamUserDomain().AdminUpdateUserStat(ctx, in.Id, in.UploadedDiff, in.DownloadedDiff, in.BonusDiff)
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "admin.user.update_stat_failed"))
	}
	if rows == 0 {
		return gerror.New(gi18n.T(ctx, "admin.user.stat_no_changes"))
	}
	return nil
}

func (s *sAdminIamUserUsecase) Ban(ctx context.Context, actor *model.Actor, in adminin.IamUserBanInp) error {
	err := service.ModUserUsecase().Apply(ctx, actor, modin.ApplyModInp{
		UserId:       in.Id,
		ModType:      consts.ModUserTypeBanned,
		Reason:       in.Reason,
		DurationDays: in.DurationDays,
	})
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "admin.user.ban_failed_mod"))
	}
	return nil
}
