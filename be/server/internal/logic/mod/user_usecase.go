package mod

import (
	"context"
	"fmt"
	"time"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/modin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/modout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

type sModUserUsecase struct{}

const expiredUserModCleanupLimit = 500

func init() {
	service.RegisterModUserUsecase(NewModUserUsecase())
}

func NewModUserUsecase() *sModUserUsecase {
	return &sModUserUsecase{}
}

func (s *sModUserUsecase) Apply(ctx context.Context, actor *model.Actor, in modin.ApplyModInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	if !s.isValidModType(in.ModType) {
		return gerror.New("invalid moderation type")
	}
	if in.DurationDays < 0 || in.DurationSeconds < 0 {
		return gerror.New("invalid moderation duration")
	}

	var expireAt *gtime.Time
	if in.DurationSeconds > 0 {
		expireAt = gtime.Now().Add(time.Duration(in.DurationSeconds) * time.Second)
	} else if in.DurationDays > 0 {
		expireAt = gtime.Now().Add(time.Duration(in.DurationDays) * 24 * time.Hour)
	}

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		modId, err := service.ModUserDomain().Create(ctx, entity.ModUserLog{
			UserId:   in.UserId,
			ModType:  in.ModType,
			Reason:   in.Reason,
			ExpireAt: expireAt,
			ModBy:    actor.Id,
			IsActive: true,
		})
		if err != nil {
			return err
		}

		return service.IamPermissionDomain().GrantUserPermissionsBySource(
			ctx,
			in.UserId,
			s.getModDenyPermissions(in.ModType),
			consts.IamUserPermissionSourceUserMod,
			modId,
			expireAt,
		)
	})
	if err != nil {
		return err
	}

	service.IamUserUsecase().InvalidateUserCache(ctx, in.UserId)
	notifyInp := sitein.MessageNotifyInp{
		ReceiverId:  in.UserId,
		TitleKey:    s.restrictionTitleKey("applied", in.ModType),
		ContentKey:  "site.message.restriction.applied.permanent.content",
		ContentArgs: []any{in.Reason},
	}
	if expireAt != nil {
		notifyInp.ContentKey = "site.message.restriction.applied.temporary.content"
		notifyInp.ContentArgs = []any{expireAt.Format("Y-m-d H:i"), in.Reason}
	}
	if in.Reason == consts.IamUserRankAutoBanReason {
		notifyInp.TitleKey = "site.message.rank.auto_banned.title"
		notifyInp.ContentKey = "site.message.rank.auto_banned.content"
		notifyInp.ContentArgs = nil
	}
	service.SiteMessageUsecase().Notify(ctx, notifyInp)
	return nil
}

func (s *sModUserUsecase) isValidModType(modType int) bool {
	switch modType {
	case consts.ModUserTypeWarned,
		consts.ModUserTypeBanned,
		consts.ModUserTypeLeechWarned,
		consts.ModUserTypeUploadBanned,
		consts.ModUserTypeDownloadBanned,
		consts.ModUserTypeForumBanned:
		return true
	default:
		return false
	}
}

func (s *sModUserUsecase) restrictionTitleKey(action string, modType int) string {
	typeKey := ""
	switch modType {
	case consts.ModUserTypeWarned:
		typeKey = "warned"
	case consts.ModUserTypeBanned:
		typeKey = "banned"
	case consts.ModUserTypeLeechWarned:
		typeKey = "leech_warned"
	case consts.ModUserTypeUploadBanned:
		typeKey = "upload_banned"
	case consts.ModUserTypeDownloadBanned:
		typeKey = "download_banned"
	case consts.ModUserTypeForumBanned:
		typeKey = "forum_banned"
	}
	if typeKey == "" {
		return fmt.Sprintf("site.message.restriction.%s.title", action)
	}
	return fmt.Sprintf("site.message.restriction.%s.%s.title", action, typeKey)
}

func (s *sModUserUsecase) getModDenyPermissions(modType int) []string {
	switch modType {
	case consts.ModUserTypeBanned:
		return []string{s.denyPermission(consts.IamPermissionAll)}
	case consts.ModUserTypeUploadBanned:
		return []string{
			s.denyPermission(consts.IamPermissionCatalogTorrentCreate),
			s.denyPermission(consts.IamPermissionCatalogSubtitleCreate),
		}
	case consts.ModUserTypeDownloadBanned:
		return []string{
			s.denyPermission(consts.IamPermissionCatalogTorrentDownload),
			s.denyPermission(consts.IamPermissionCatalogSubtitleDownload),
		}
	case consts.ModUserTypeForumBanned:
		return []string{
			s.denyPermission(consts.IamPermissionForumTopicRead),
			s.denyPermission(consts.IamPermissionForumTopicCreate),
			s.denyPermission(consts.IamPermissionForumTopicUpdate),
			s.denyPermission(consts.IamPermissionForumReplyRead),
			s.denyPermission(consts.IamPermissionForumReplyCreate),
			s.denyPermission(consts.IamPermissionForumReplyUpdate),
		}
	default:
		return nil
	}
}

func (s *sModUserUsecase) denyPermission(perm string) string {
	return "-" + perm
}

func (s *sModUserUsecase) Remove(ctx context.Context, actor *model.Actor, in modin.RemoveModInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	mod, err := service.ModUserDomain().GetById(ctx, in.Id)
	if err != nil || mod == nil {
		return gerror.New("Mod record not found")
	}

	if mod.UserId != in.UserId {
		return gerror.New("User mismatch")
	}

	columns := dao.ModUserLog.Columns()
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := service.ModUserDomain().Update(ctx, in.Id, g.Map{
			columns.IsActive:  false,
			columns.UpdatedAt: gtime.Now(),
		}); err != nil {
			return err
		}

		return service.IamPermissionDomain().DeactivateUserPermissionsBySource(
			ctx,
			in.UserId,
			consts.IamUserPermissionSourceUserMod,
			in.Id,
		)
	})
	if err != nil {
		return err
	}
	service.IamUserUsecase().InvalidateUserCache(ctx, in.UserId)
	service.SiteMessageUsecase().Notify(ctx, sitein.MessageNotifyInp{
		ReceiverId: in.UserId,
		TitleKey:   s.restrictionTitleKey("removed", mod.ModType),
		ContentKey: "site.message.restriction.removed.content",
	})
	return nil
}

func (s *sModUserUsecase) CleanupExpired(ctx context.Context) (int, error) {
	now := gtime.Now()
	records, err := service.ModUserDomain().QueryExpiredActiveMods(ctx, now, expiredUserModCleanupLimit)
	if err != nil {
		return 0, err
	}
	if len(records) == 0 {
		return 0, nil
	}

	columns := dao.ModUserLog.Columns()
	userIds := make(map[uint64]struct{}, len(records))
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, record := range records {
			if record.Id == 0 || record.UserId == 0 {
				continue
			}

			if err := service.ModUserDomain().Update(ctx, record.Id, g.Map{
				columns.IsActive:  false,
				columns.UpdatedAt: now,
			}); err != nil {
				return err
			}

			if err := service.IamPermissionDomain().DeactivateUserPermissionsBySource(
				ctx,
				record.UserId,
				consts.IamUserPermissionSourceUserMod,
				record.Id,
			); err != nil {
				return err
			}
			userIds[record.UserId] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	for userId := range userIds {
		service.IamUserUsecase().InvalidateUserCache(ctx, userId)
	}
	for _, record := range records {
		if record.Id == 0 || record.UserId == 0 {
			continue
		}
		service.SiteMessageUsecase().Notify(ctx, sitein.MessageNotifyInp{
			ReceiverId: record.UserId,
			TitleKey:   s.restrictionTitleKey("expired", record.ModType),
			ContentKey: "site.message.restriction.expired.content",
		})
	}
	return len(records), nil
}

func (s *sModUserUsecase) List(ctx context.Context, actor *model.Actor, in modin.ListUserInp) (*modout.ListUserOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	var records []entity.ModUserLog
	records, _, err := service.ModUserDomain().QueryUserLogs(ctx, in.UserId, 1, 1000)
	if err != nil {
		return nil, err
	}

	var list []modout.UserModItem
	for _, r := range records {
		list = append(list, modout.UserModItem{
			Id:         r.Id,
			UserId:     r.UserId,
			ModType:    r.ModType,
			Reason:     r.Reason,
			ExpireAt:   r.ExpireAt,
			ModBy:      r.ModBy,
			ModComment: r.ModComment,
			IsActive:   r.IsActive,
			CreatedAt:  r.CreatedAt,
		})
	}

	return &modout.ListUserOut{List: list}, nil
}
