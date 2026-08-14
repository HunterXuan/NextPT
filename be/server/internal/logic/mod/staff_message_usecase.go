package mod

import (
	"context"
	"strings"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/model/in/modin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/modout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

type sModStaffMessageUsecase struct{}

func init() {
	service.RegisterModStaffMessageUsecase(NewModStaffMessageUsecase())
}

func NewModStaffMessageUsecase() *sModStaffMessageUsecase {
	return &sModStaffMessageUsecase{}
}

func (s *sModStaffMessageUsecase) Create(ctx context.Context, actor *model.Actor, in modin.StaffMessageCreateInp) (*modout.StaffMessageCreateOut, error) {
	if actor == nil || actor.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	subject := strings.TrimSpace(in.Subject)
	content := strings.TrimSpace(in.Content)
	if subject == "" || content == "" {
		return nil, gerror.New(gi18n.T(ctx, "mod.staff_message.content_required"))
	}
	id, err := service.ModStaffMessageDomain().Create(ctx, entity.ModStaffMessage{
		SenderId: actor.Id,
		Subject:  subject,
		Content:  content,
	})
	if err != nil {
		return nil, err
	}
	return &modout.StaffMessageCreateOut{Id: id}, nil
}

func (s *sModStaffMessageUsecase) List(ctx context.Context, actor *model.Actor, in modin.StaffMessageListInp) (*modout.StaffMessageListOut, error) {
	if actor == nil || actor.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	list, total, err := service.ModStaffMessageDomain().ListBySender(ctx, actor.Id, in)
	if err != nil {
		return nil, err
	}
	return s.buildList(ctx, list, total, in.Page, in.Size), nil
}

func (s *sModStaffMessageUsecase) AdminList(ctx context.Context, actor *model.Actor, in modin.AdminStaffMessageListInp) (*modout.StaffMessageListOut, error) {
	if actor == nil || actor.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	list, total, err := service.ModStaffMessageDomain().AdminList(ctx, in)
	if err != nil {
		return nil, err
	}
	return s.buildList(ctx, list, total, in.Page, in.Size), nil
}

func (s *sModStaffMessageUsecase) AdminProcess(ctx context.Context, actor *model.Actor, in modin.StaffMessageProcessInp) error {
	if actor == nil || actor.Id == 0 {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	message, err := service.ModStaffMessageDomain().GetById(ctx, in.Id)
	if err != nil || message == nil {
		return gerror.New(gi18n.T(ctx, "mod.staff_message.not_found"))
	}
	answer := strings.TrimSpace(in.Answer)
	if message.Status == consts.ModStaffMessageStatusProcessed {
		return gerror.New(gi18n.T(ctx, "mod.staff_message.already_processed"))
	}
	now := gtime.Now()
	if err := service.ModStaffMessageDomain().Update(ctx, in.Id, do.ModStaffMessage{
		Status:     consts.ModStaffMessageStatusProcessed,
		AnsweredBy: actor.Id,
		Answer:     answer,
		AnsweredAt: now,
		UpdatedAt:  now,
	}); err != nil {
		return err
	}

	if answer != "" {
		service.SiteMessageUsecase().Notify(ctx, sitein.MessageNotifyInp{
			ReceiverId:  message.SenderId,
			TitleKey:    "site.message.staff_message.processed_with_reply.title",
			ContentKey:  "site.message.staff_message.processed_with_reply.content",
			ContentArgs: []any{message.Subject, answer},
		})
	} else {
		service.SiteMessageUsecase().Notify(ctx, sitein.MessageNotifyInp{
			ReceiverId:  message.SenderId,
			TitleKey:    "site.message.staff_message.processed.title",
			ContentKey:  "site.message.staff_message.processed.content",
			ContentArgs: []any{message.Subject},
		})
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeModStaffMessage,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"status": consts.ModStaffMessageStatusProcessed,
		},
	})
	return nil
}

func (s *sModStaffMessageUsecase) buildList(ctx context.Context, list []entity.ModStaffMessage, total int, page int, size int) *modout.StaffMessageListOut {
	userIds := make([]uint64, 0, len(list)*2)
	for _, item := range list {
		if item.SenderId > 0 {
			userIds = append(userIds, item.SenderId)
		}
		if item.AnsweredBy > 0 {
			userIds = append(userIds, item.AnsweredBy)
		}
	}
	users := s.loadUserSummaryMap(ctx, userIds)
	items := make([]modout.StaffMessageItem, 0, len(list))
	for _, item := range list {
		items = append(items, modout.StaffMessageItem{
			Id:             item.Id,
			SenderId:       item.SenderId,
			Sender:         users[item.SenderId],
			Subject:        item.Subject,
			Content:        item.Content,
			Status:         item.Status,
			AnsweredBy:     item.AnsweredBy,
			AnsweredByUser: users[item.AnsweredBy],
			Answer:         item.Answer,
			AnsweredAt:     item.AnsweredAt,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}
	return &modout.StaffMessageListOut{List: items, Total: total, Page: page, Size: size}
}

func (s *sModStaffMessageUsecase) loadUserSummaryMap(ctx context.Context, userIds []uint64) map[uint64]model.IamUserSummary {
	ids := make(map[uint64]struct{}, len(userIds))
	for _, id := range userIds {
		if id > 0 {
			ids[id] = struct{}{}
		}
	}
	result := make(map[uint64]model.IamUserSummary, len(ids))
	if len(ids) == 0 {
		return result
	}
	orderedIds := make([]uint64, 0, len(ids))
	for id := range ids {
		orderedIds = append(orderedIds, id)
		result[id] = model.IamUserSummary{Id: id}
	}
	users, err := service.IamUserDomain().GetUsersByIds(ctx, orderedIds)
	if err == nil {
		for _, user := range users {
			summary := result[user.Id]
			summary.Username = user.Username
			result[user.Id] = summary
		}
	}
	profiles, err := service.IamUserDomain().GetUserProfilesByUserIds(ctx, orderedIds)
	if err == nil {
		for _, profile := range profiles {
			summary := result[profile.UserId]
			summary.Avatar = profile.Avatar
			result[profile.UserId] = summary
		}
	}
	return result
}
