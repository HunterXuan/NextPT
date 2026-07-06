package admin

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/adminout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sAdminIamInviteUsecase struct{}

const (
	inviteHashAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	inviteHashLength   = 32
)

func NewAdminIamInviteUsecase() *sAdminIamInviteUsecase {
	return &sAdminIamInviteUsecase{}
}

func init() {
	service.RegisterAdminIamInviteUsecase(NewAdminIamInviteUsecase())
}

func (s *sAdminIamInviteUsecase) List(ctx context.Context, actor *model.Actor, in adminin.IamInviteListInp) (*adminout.IamInviteListOut, error) {
	list, total, err := service.IamInviteDomain().AdminQuerySiteInvites(ctx, in.Page, in.Size, in.Status)
	if err != nil {
		return nil, err
	}

	inviteeIds := make([]uint64, 0)
	for _, item := range list {
		if item.InviteeId > 0 {
			inviteeIds = append(inviteeIds, item.InviteeId)
		}
	}
	usernameMap, err := s.loadInviteeUsernameMap(ctx, inviteeIds)
	if err != nil {
		return nil, err
	}

	items := make([]adminout.IamInviteItem, 0, len(list))
	for _, item := range list {
		items = append(items, s.buildInviteItem(item, usernameMap))
	}
	return &adminout.IamInviteListOut{
		Invites: items,
		Total:   total,
	}, nil
}

func (s *sAdminIamInviteUsecase) Grant(ctx context.Context, actor *model.Actor, in adminin.IamInviteGrantInp) error {
	inviterIds, err := s.resolveInviteGrantInviters(ctx, in.TargetMode, in.RoleIds)
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "admin.invite.grant_failed"))
	}

	list := make([]do.IamInvite, 0, len(inviterIds)*in.Amount)
	isTemporary := in.IsTemp || in.ExpireAt != nil
	for _, inviterId := range inviterIds {
		for i := 0; i < in.Amount; i++ {
			list = append(list, do.IamInvite{
				InviterId:   inviterId,
				Hash:        s.generateInviteHash(),
				Status:      consts.IamInviteStatusUnused,
				IsTemporary: isTemporary,
				ExpireAt:    in.ExpireAt,
			})
		}
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return service.IamInviteDomain().AdminCreateInvites(ctx, list)
	})
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "admin.invite.grant_failed"))
	}

	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionCreate,
		TargetType: consts.SiteAuditTargetTypeIamInvite,
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"operation":    consts.SiteAuditOperationGrant,
			"amount":       in.Amount,
			"targetMode":   in.TargetMode,
			"roleIds":      in.RoleIds,
			"inviterCount": len(inviterIds),
			"isTemporary":  isTemporary,
			"expireAt":     in.ExpireAt,
		},
	})
	return nil
}

func (s *sAdminIamInviteUsecase) Recycle(ctx context.Context, actor *model.Actor, in adminin.IamInviteRecycleInp) error {
	invite, err := service.IamInviteDomain().AdminGetSiteInviteById(ctx, in.Id)
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "admin.invite.recycle_failed"))
	}
	if invite == nil {
		return gerror.New(gi18n.T(ctx, "admin.invite.not_found"))
	}
	if !s.canRecycleInvite(invite.Status) {
		return gerror.New(gi18n.T(ctx, "admin.invite.invalid_status"))
	}
	if err := service.IamInviteDomain().AdminUpdateSiteInviteStatus(ctx, in.Id, consts.IamInviteStatusRecycled); err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "admin.invite.recycle_failed"))
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeIamInvite,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"operation": consts.SiteAuditOperationRecycle,
		},
	})
	return nil
}

func (s *sAdminIamInviteUsecase) resolveInviteGrantInviters(ctx context.Context, targetMode string, roleIds []uint) ([]uint64, error) {
	switch targetMode {
	case consts.IamInviteGrantTargetSite:
		return []uint64{0}, nil
	case consts.IamInviteGrantTargetRoles:
		return s.resolveRoleInviteGrantInviters(ctx, roleIds)
	default:
		return nil, gerror.New(gi18n.T(ctx, "admin.invite.invalid_target"))
	}
}

func (s *sAdminIamInviteUsecase) resolveRoleInviteGrantInviters(ctx context.Context, roleIds []uint) ([]uint64, error) {
	roleIds = s.normalizeRoleIds(roleIds)
	if len(roleIds) == 0 {
		return nil, gerror.New(gi18n.T(ctx, "admin.invite.role_required"))
	}

	roles, err := service.IamRoleDomain().GetRolesByIds(ctx, roleIds)
	if err != nil {
		return nil, err
	}
	if len(roles) != len(roleIds) {
		return nil, gerror.New(gi18n.T(ctx, "admin.role.not_found"))
	}

	userIds, err := service.IamUserDomain().GetUserIdsByRoles(ctx, roleIds)
	if err != nil {
		return nil, err
	}
	if len(userIds) == 0 {
		return nil, gerror.New(gi18n.T(ctx, "admin.invite.no_target_users"))
	}
	return userIds, nil
}

func (s *sAdminIamInviteUsecase) normalizeRoleIds(roleIds []uint) []uint {
	seen := make(map[uint]struct{}, len(roleIds))
	normalized := make([]uint, 0, len(roleIds))
	for _, roleId := range roleIds {
		if roleId == 0 {
			continue
		}
		if _, ok := seen[roleId]; ok {
			continue
		}
		seen[roleId] = struct{}{}
		normalized = append(normalized, roleId)
	}
	return normalized
}

func (s *sAdminIamInviteUsecase) generateInviteHash() string {
	return grand.Str(inviteHashAlphabet, inviteHashLength)
}

func (s *sAdminIamInviteUsecase) loadInviteeUsernameMap(ctx context.Context, inviteeIds []uint64) (map[uint64]string, error) {
	usernameMap := make(map[uint64]string)
	if len(inviteeIds) == 0 {
		return usernameMap, nil
	}
	users, err := service.IamUserDomain().GetUsersByIds(ctx, inviteeIds)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		usernameMap[user.Id] = user.Username
	}
	return usernameMap, nil
}

func (s *sAdminIamInviteUsecase) buildInviteItem(invite entity.IamInvite, usernameMap map[uint64]string) adminout.IamInviteItem {
	return adminout.IamInviteItem{
		Id:           invite.Id,
		InviterId:    invite.InviterId,
		InviteeEmail: invite.InviteeEmail,
		InviteeId:    invite.InviteeId,
		InviteeName:  usernameMap[invite.InviteeId],
		Hash:         invite.Hash,
		Status:       invite.Status,
		IsTemporary:  invite.IsTemporary,
		ExpireAt:     invite.ExpireAt,
		UsedAt:       invite.UsedAt,
		CreatedAt:    invite.CreatedAt,
	}
}

func (s *sAdminIamInviteUsecase) canRecycleInvite(status int) bool {
	switch status {
	case consts.IamInviteStatusUnused,
		consts.IamInviteStatusSent,
		consts.IamInviteStatusExpired:
		return true
	default:
		return false
	}
}
