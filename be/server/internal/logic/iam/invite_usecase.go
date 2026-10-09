package iam

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/iamin"
	"server/internal/model/out/iamout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/glog"
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
			Hash:         s.inviteListHash(item),
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

	invite, err := service.IamInviteDomain().GetInviteByInviterIdAndId(ctx, actor.Id, in.Id)
	if err != nil {
		return err
	}
	if invite == nil {
		return gerror.New(gi18n.T(ctx, "iam.invite.not_found"))
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		lockedInvite, err := service.IamInviteDomain().GetInviteByHashForUpdate(ctx, invite.Hash)
		if err != nil {
			return err
		}
		return s.markInviteSent(ctx, actor, in, lockedInvite)
	})
	if err != nil {
		return err
	}
	s.sendInvitationMail(ctx, strings.TrimSpace(in.Email), invite)
	return nil
}

func (s *sIamInviteUsecase) markInviteSent(ctx context.Context, actor *model.Actor, in iamin.InviteSendInp, invite *entity.IamInvite) error {
	if invite == nil || invite.InviterId != actor.Id || invite.Id != in.Id {
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
		InviteeEmail: strings.TrimSpace(in.Email),
	})
}

func (s *sIamInviteUsecase) sendInvitationMail(ctx context.Context, recipient string, invite *entity.IamInvite) {
	mail := s.invitationMail(ctx, recipient, invite)
	textBody, htmlBody := mail.Bodies()
	if err := service.SysMailgun().SendHtmlMail(ctx, mail.Subject, textBody, htmlBody, mail.Recipient); err != nil {
		glog.Warningf(ctx, "send invitation mail failed: inviteId=%d error=%v", invite.Id, err)
	}
}

func (s *sIamInviteUsecase) invitationMail(ctx context.Context, recipient string, invite *entity.IamInvite) model.IamAccountActionMail {
	mailCtx := gi18n.WithLanguage(ctx, g.Cfg().MustGet(ctx, "i18n.default", "zh-CN").String())
	siteName := strings.TrimSpace(g.Cfg().MustGet(ctx, "site.name", "NextPT").String())
	siteURL := strings.TrimRight(strings.TrimSpace(g.Cfg().MustGet(ctx, "site.url", "http://localhost:3000").String()), "/")
	if siteName == "" {
		siteName = "NextPT"
	}
	if siteURL == "" {
		siteURL = "http://localhost:3000"
	}
	expiry := gi18n.T(mailCtx, "iam.invite.mail_permanent")
	if invite.ExpireAt != nil {
		expiry = fmt.Sprintf(gi18n.T(mailCtx, "iam.invite.mail_expiry"), invite.ExpireAt.Format("Y-m-d H:i:s T"))
	}
	return model.IamAccountActionMail{
		Kind: "invitation", Recipient: recipient,
		Subject:   fmt.Sprintf(gi18n.T(mailCtx, "iam.invite.mail_subject"), siteName),
		Greeting:  gi18n.T(mailCtx, "iam.invite.mail_greeting"),
		Intro:     fmt.Sprintf(gi18n.T(mailCtx, "iam.invite.mail_intro"), siteName, recipient),
		Action:    gi18n.T(mailCtx, "iam.invite.mail_action"),
		ActionURL: siteURL + "/register?invite=" + url.QueryEscape(invite.Hash),
		Expiry:    expiry, Note: gi18n.T(mailCtx, "iam.invite.mail_ignore"),
	}
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

func (s *sIamInviteUsecase) inviteListHash(invite entity.IamInvite) string {
	if invite.Status == consts.IamInviteStatusUnused {
		return s.maskInviteHash(invite.Hash)
	}
	return invite.Hash
}

func (s *sIamInviteUsecase) maskInviteHash(hash string) string {
	if len(hash) <= 12 {
		return hash
	}
	return hash[:6] + "..." + hash[len(hash)-6:]
}
