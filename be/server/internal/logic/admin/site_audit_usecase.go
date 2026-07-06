package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/adminout"
	"server/internal/service"
)

type sAdminSiteAuditUsecase struct{}

func init() {
	service.RegisterAdminSiteAuditUsecase(NewAdminSiteAuditUsecase())
}

func NewAdminSiteAuditUsecase() *sAdminSiteAuditUsecase {
	return &sAdminSiteAuditUsecase{}
}

func (s *sAdminSiteAuditUsecase) List(ctx context.Context, actor *model.Actor, in adminin.SiteAuditListInp) (*adminout.SiteAuditListOut, error) {
	logs, total, err := service.SiteAuditDomain().AdminListAudits(ctx, sitein.AuditListInp{
		Page:       in.Page,
		Size:       in.Size,
		Level:      in.Level,
		Action:     in.Action,
		TargetType: in.TargetType,
		UserId:     in.UserId,
		StartAt:    in.StartAt,
		EndAt:      in.EndAt,
	})
	if err != nil {
		return nil, err
	}

	userMap, err := s.loadActorSummaryMap(ctx, logs)
	if err != nil {
		return nil, err
	}

	list := make([]*adminout.SiteAuditItem, 0, len(logs))
	for _, l := range logs {
		list = append(list, &adminout.SiteAuditItem{
			Id:         l.Id,
			UserId:     l.UserId,
			Actor:      s.actorSummary(userMap, l.UserId),
			Action:     l.Action,
			TargetType: l.TargetType,
			TargetId:   l.TargetId,
			Level:      l.Level,
			Ip:         l.Ip,
			Detail:     l.Detail,
			CreatedAt:  l.CreatedAt,
		})
	}

	return &adminout.SiteAuditListOut{
		List:  list,
		Total: total,
		Page:  in.Page,
		Size:  in.Size,
	}, nil
}

func (s *sAdminSiteAuditUsecase) loadActorSummaryMap(ctx context.Context, logs []entity.SiteAudit) (map[uint64]model.IamUserSummary, error) {
	userIds := make([]uint64, 0, len(logs))
	seen := make(map[uint64]struct{}, len(logs))
	for _, log := range logs {
		if log.UserId == 0 {
			continue
		}
		if _, ok := seen[log.UserId]; ok {
			continue
		}
		seen[log.UserId] = struct{}{}
		userIds = append(userIds, log.UserId)
	}

	userMap := make(map[uint64]model.IamUserSummary, len(userIds))
	for _, userId := range userIds {
		userMap[userId] = model.IamUserSummary{Id: userId}
	}
	if len(userIds) == 0 {
		return userMap, nil
	}

	users, err := service.IamUserDomain().GetUsersByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		summary := userMap[user.Id]
		summary.Id = user.Id
		summary.Username = user.Username
		userMap[user.Id] = summary
	}

	profiles, err := service.IamUserDomain().GetUserProfilesByUserIds(ctx, userIds)
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		summary := userMap[profile.UserId]
		summary.Id = profile.UserId
		summary.Avatar = profile.Avatar
		userMap[profile.UserId] = summary
	}
	return userMap, nil
}

func (s *sAdminSiteAuditUsecase) actorSummary(userMap map[uint64]model.IamUserSummary, userId uint64) model.IamUserSummary {
	if userId == 0 {
		return model.IamUserSummary{}
	}
	if summary, ok := userMap[userId]; ok {
		return summary
	}
	return model.IamUserSummary{Id: userId}
}
