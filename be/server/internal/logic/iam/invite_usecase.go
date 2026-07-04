package iam

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/iamin"
	"server/internal/model/out/iamout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

type sIamInviteUsecase struct{}

func init() {
	service.RegisterIamInviteUsecase(NewIamInviteUsecase())
}

func NewIamInviteUsecase() *sIamInviteUsecase {
	return &sIamInviteUsecase{}
}

func (s *sIamInviteUsecase) List(ctx context.Context, actor *model.Actor, in iamin.InviteListInp) (*iamout.InviteListOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	list, total, err := service.IamInviteDomain().QueryInvitesByInviter(ctx, actor.Id, in.Page, in.Size, in.Status)
	if err != nil {
		return nil, err
	}

	inviteeIds := make([]uint64, 0)
	for _, item := range list {
		if item.InviteeId > 0 {
			inviteeIds = append(inviteeIds, item.InviteeId)
		}
	}

	usernameMap := make(map[uint64]string)
	if len(inviteeIds) > 0 {
		var users []entity.IamUser
		users, err = service.IamUserDomain().GetUsersByIds(ctx, inviteeIds)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			usernameMap[u.Id] = u.Username
		}
	}

	outList := make([]iamout.InviteListItem, len(list))
	for i, item := range list {
		outList[i] = iamout.InviteListItem{
			Id:           item.Id,
			InviterId:    item.InviterId,
			InviteeEmail: item.InviteeEmail,
			InviteeId:    item.InviteeId,
			InviteeName:  usernameMap[item.InviteeId],
			Hash:         item.Hash,
			Status:       item.Status,
			IsTemporary:  item.IsTemporary,
			ExpireAt:     item.ExpireAt,
			UsedAt:       item.UsedAt,
			CreatedAt:    item.CreatedAt,
		}
	}

	return &iamout.InviteListOut{
		List:  outList,
		Total: total,
	}, nil
}

func (s *sIamInviteUsecase) Send(ctx context.Context, actor *model.Actor, in iamin.InviteSendInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	invite, err := service.IamInviteDomain().GetInvitesByInviterIdAndHash(ctx, actor.Id, in.Hash)
	if err != nil {
		return err
	}
	if invite == nil {
		return gerror.New(gi18n.T(ctx, "iam.invite.not_found"))
	}
	if invite.Status != consts.IamInviteStatusUnused {
		return gerror.New(gi18n.T(ctx, "iam.invite.invalid_status"))
	}
	if invite.ExpireAt != nil && !invite.ExpireAt.After(gtime.Now()) {
		return gerror.New(gi18n.T(ctx, "iam.invite.expired"))
	}

	return service.IamInviteDomain().UpdateInvite(ctx, invite.Id, do.IamInvite{
		Status:       consts.IamInviteStatusSent,
		InviteeEmail: in.Email,
	})
}

func (s *sIamInviteUsecase) Check(ctx context.Context, in iamin.InviteCheckInp) (*iamout.InviteCheckOut, error) {
	invite, err := service.IamInviteDomain().GetInviteByHash(ctx, in.Hash)
	if err != nil {
		return nil, err
	}
	if invite == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.invite.not_found"))
	}
	if invite.Status != consts.IamInviteStatusSent && invite.Status != consts.IamInviteStatusUnused {
		return nil, gerror.New(gi18n.T(ctx, "iam.invite.invalid_status"))
	}
	if invite.ExpireAt != nil && !invite.ExpireAt.After(gtime.Now()) {
		return nil, gerror.New(gi18n.T(ctx, "iam.invite.expired"))
	}

	user, err := service.IamUserDomain().GetUserById(ctx, invite.InviterId)
	if err != nil {
		return nil, err
	}

	var inviterUsername string
	if user != nil {
		inviterUsername = user.Username
	}

	return &iamout.InviteCheckOut{
		Hash:            invite.Hash,
		InviterId:       invite.InviterId,
		InviterUsername: inviterUsername,
	}, nil
}

func (s *sIamInviteUsecase) CleanupExpired(ctx context.Context) (int64, error) {
	return service.IamInviteDomain().ExpireInvites(ctx, gtime.Now())
}
