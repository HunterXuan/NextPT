package admin

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/adminout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
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
	items, err := s.buildUserItems(ctx, users)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "admin.user.fetch_failed"))
	}
	return &adminout.IamUserListOut{
		Users: items,
		Total: total,
	}, nil
}

func (s *sAdminIamUserUsecase) buildUserItems(ctx context.Context, users []*entity.IamUser) ([]adminout.IamUserItem, error) {
	items := make([]adminout.IamUserItem, 0, len(users))
	if len(users) == 0 {
		return items, nil
	}

	userIds := make([]uint64, 0, len(users))
	for _, user := range users {
		if user != nil && user.Id > 0 {
			userIds = append(userIds, user.Id)
		}
	}
	profileMap, err := s.loadUserProfileMap(ctx, userIds)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user == nil {
			continue
		}
		item := adminout.IamUserItem{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			Passkey:   user.Passkey,
			Status:    user.Status,
			Role:      user.Role,
			VipUntil:  user.VipUntil,
			VipRemark: user.VipRemark,
			InvitedBy: user.InvitedBy,
			LastLogin: user.LastLogin,
			LastIp:    user.LastIp,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		if profile, ok := profileMap[user.Id]; ok {
			item.Avatar = profile.Avatar
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *sAdminIamUserUsecase) loadUserProfileMap(ctx context.Context, userIds []uint64) (map[uint64]entity.IamUserProfile, error) {
	profileMap := make(map[uint64]entity.IamUserProfile)
	if len(userIds) == 0 {
		return profileMap, nil
	}

	profiles, err := service.IamUserDomain().GetUserProfilesByUserIds(ctx, userIds)
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		profileMap[profile.UserId] = profile
	}
	return profileMap, nil
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

	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeIamUser,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"statusChanged":  in.Status != nil,
			"roleChanged":    in.Role != nil,
			"passkeyChanged": in.Passkey != nil,
		},
	})
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
	hasStatDiff := s.hasStatDiff(in.UploadedDiff, in.DownloadedDiff)
	hasBonusDiff := in.BonusDiff != nil && *in.BonusDiff != 0
	if !hasStatDiff && !hasBonusDiff {
		return gerror.New(gi18n.T(ctx, "admin.user.stat_no_changes"))
	}

	changed := false
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if hasStatDiff {
			rows, err := service.IamUserDomain().AdminUpdateUserStat(ctx, in.Id, in.UploadedDiff, in.DownloadedDiff)
			if err != nil {
				return err
			}
			if rows > 0 {
				changed = true
			}
		}

		if hasBonusDiff {
			if err := service.EconomyBonusUsecase().AddBonus(ctx, in.Id, *in.BonusDiff, consts.EconomyBonusActionAdminAdjustment, "", 0, "", ""); err != nil {
				return err
			}
			changed = true
		}

		return nil
	})
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "admin.user.update_stat_failed"))
	}
	if !changed {
		return gerror.New(gi18n.T(ctx, "admin.user.stat_no_changes"))
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdateStat,
		TargetType: consts.SiteAuditTargetTypeIamUser,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"uploadedDiff":   in.UploadedDiff,
			"downloadedDiff": in.DownloadedDiff,
			"bonusDiff":      in.BonusDiff,
		},
	})
	return nil
}

func (s *sAdminIamUserUsecase) hasStatDiff(uploadedDiff, downloadedDiff *int64) bool {
	return (uploadedDiff != nil && *uploadedDiff != 0) || (downloadedDiff != nil && *downloadedDiff != 0)
}
