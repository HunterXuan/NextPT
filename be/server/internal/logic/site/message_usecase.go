package site

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/sitein"
	"server/internal/model/out/siteout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
)

type sSiteMessageUsecase struct{}

func init() {
	service.RegisterSiteMessageUsecase(NewSiteMessageUsecase())
}

func NewSiteMessageUsecase() *sSiteMessageUsecase {
	return &sSiteMessageUsecase{}
}

func (s *sSiteMessageUsecase) List(ctx context.Context, actor *model.Actor, in sitein.MessageListInp) (*siteout.MessageListOut, error) {
	list, total, err := service.SiteMessageDomain().ListByReceiver(ctx, s.actorId(actor), in)
	if err != nil {
		return nil, err
	}
	return &siteout.MessageListOut{
		List:  s.buildMessageItems(ctx, list),
		Total: total,
		Page:  in.Page,
		Size:  in.Size,
	}, nil
}

func (s *sSiteMessageUsecase) MarkRead(ctx context.Context, actor *model.Actor, in sitein.MessageReadInp) error {
	return service.SiteMessageDomain().MarkRead(ctx, s.actorId(actor), in.Id)
}

func (s *sSiteMessageUsecase) MarkAllRead(ctx context.Context, actor *model.Actor, in sitein.MessageReadAllInp) error {
	return service.SiteMessageDomain().MarkAllRead(ctx, s.actorId(actor))
}

func (s *sSiteMessageUsecase) Notify(ctx context.Context, in sitein.MessageCreateInp) error {
	if in.ReceiverId == 0 {
		return nil
	}
	_, err := service.SiteMessageDomain().Create(ctx, in)
	return err
}

func (s *sSiteMessageUsecase) AdminList(ctx context.Context, actor *model.Actor, in sitein.AdminMessageListInp) (*siteout.MessageListOut, error) {
	list, total, err := service.SiteMessageDomain().AdminList(ctx, in)
	if err != nil {
		return nil, err
	}
	return &siteout.MessageListOut{
		List:  s.buildMessageItems(ctx, list),
		Total: total,
		Page:  in.Page,
		Size:  in.Size,
	}, nil
}

func (s *sSiteMessageUsecase) AdminCreate(ctx context.Context, actor *model.Actor, in sitein.AdminMessageCreateInp) (*siteout.MessageCreateOut, error) {
	receiverIds := s.uniqueReceiverIds(in.ReceiverIds)
	if len(receiverIds) == 0 {
		return nil, gerror.New("receiver is required")
	}
	items := make([]sitein.MessageCreateInp, 0, len(receiverIds))
	for _, receiverId := range receiverIds {
		items = append(items, sitein.MessageCreateInp{
			SenderId:   s.actorId(actor),
			ReceiverId: receiverId,
			Title:      in.Title,
			Content:    in.Content,
			TargetType: in.TargetType,
			TargetId:   in.TargetId,
		})
	}
	if err := service.SiteMessageDomain().BatchCreate(ctx, items); err != nil {
		return nil, err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionCreate,
		TargetType: consts.SiteAuditTargetTypeSiteMessage,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"count":      len(items),
			"title":      in.Title,
			"targetType": in.TargetType,
			"targetId":   in.TargetId,
		},
	})
	return &siteout.MessageCreateOut{Count: len(items)}, nil
}

func (s *sSiteMessageUsecase) buildMessageItems(ctx context.Context, list []entity.SiteMessage) []*siteout.MessageItem {
	userMap := s.loadUserSummaryMap(ctx, s.messageUserIds(list))
	items := make([]*siteout.MessageItem, 0, len(list))
	for _, item := range list {
		items = append(items, &siteout.MessageItem{
			Id:         item.Id,
			SenderId:   item.SenderId,
			Sender:     s.userSummary(userMap, item.SenderId),
			ReceiverId: item.ReceiverId,
			Receiver:   s.userSummary(userMap, item.ReceiverId),
			Title:      item.Title,
			Content:    item.Content,
			TargetType: item.TargetType,
			TargetId:   item.TargetId,
			IsRead:     item.IsRead,
			ReadAt:     item.ReadAt,
			CreatedAt:  item.CreatedAt,
		})
	}
	return items
}

func (s *sSiteMessageUsecase) messageUserIds(list []entity.SiteMessage) []uint64 {
	ids := make([]uint64, 0, len(list)*2)
	for _, item := range list {
		if item.SenderId > 0 {
			ids = append(ids, item.SenderId)
		}
		if item.ReceiverId > 0 {
			ids = append(ids, item.ReceiverId)
		}
	}
	return ids
}

func (s *sSiteMessageUsecase) loadUserSummaryMap(ctx context.Context, userIds []uint64) map[uint64]model.IamUserSummary {
	userMap := make(map[uint64]model.IamUserSummary)
	uniqueIds := s.uniqueUserIds(userIds)
	for _, id := range uniqueIds {
		userMap[id] = model.IamUserSummary{Id: id}
	}
	if len(uniqueIds) == 0 {
		return userMap
	}

	users, err := service.IamUserDomain().GetUsersByIds(ctx, uniqueIds)
	if err != nil {
		return userMap
	}
	for _, user := range users {
		summary := userMap[user.Id]
		summary.Id = user.Id
		summary.Username = user.Username
		userMap[user.Id] = summary
	}

	profiles, err := service.IamUserDomain().GetUserProfilesByUserIds(ctx, uniqueIds)
	if err != nil {
		return userMap
	}
	for _, profile := range profiles {
		summary := userMap[profile.UserId]
		summary.Id = profile.UserId
		summary.Avatar = profile.Avatar
		userMap[profile.UserId] = summary
	}
	return userMap
}

func (s *sSiteMessageUsecase) userSummary(userMap map[uint64]model.IamUserSummary, userId uint64) model.IamUserSummary {
	if userId == 0 {
		return model.IamUserSummary{}
	}
	if summary, ok := userMap[userId]; ok {
		return summary
	}
	return model.IamUserSummary{Id: userId}
}

func (s *sSiteMessageUsecase) uniqueReceiverIds(ids []uint64) []uint64 {
	uniqueIds := s.uniqueUserIds(ids)
	out := make([]uint64, 0, len(uniqueIds))
	for _, id := range uniqueIds {
		if id > 0 {
			out = append(out, id)
		}
	}
	return out
}

func (s *sSiteMessageUsecase) uniqueUserIds(ids []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(ids))
	uniqueIds := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIds = append(uniqueIds, id)
	}
	return uniqueIds
}

func (s *sSiteMessageUsecase) actorId(actor *model.Actor) uint64 {
	if actor == nil {
		return 0
	}
	return actor.Id
}
