package mod

import (
	"context"
	"fmt"
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

type sModReportUsecase struct{}

func init() {
	service.RegisterModReportUsecase(NewModReportUsecase())
}

func NewModReportUsecase() *sModReportUsecase {
	return &sModReportUsecase{}
}

func (s *sModReportUsecase) Create(ctx context.Context, actor *model.Actor, in modin.CreateReportInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	count, err := service.ModReportDomain().GetPendingCount(ctx, actor.Id, in.TargetType, in.TargetId)
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New(gi18n.T(ctx, "mod.report.duplicate"))
	}

	return service.ModReportDomain().Create(ctx, entity.ModReport{
		ReporterId: actor.Id,
		TargetType: in.TargetType,
		TargetId:   in.TargetId,
		Reason:     in.Reason,
		Status:     consts.ModReportStatusPending,
	})
}

func (s *sModReportUsecase) List(ctx context.Context, actor *model.Actor, in modin.ListReportsInp) (*modout.ListReportsOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	records, total, err := service.ModReportDomain().QueryReports(ctx, in.Status, in.TargetType, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	targetMap := s.loadReportTargetSummaryMap(ctx, records)
	userIds := make([]uint64, 0, len(records)*2)
	for _, r := range records {
		userIds = append(userIds, r.ReporterId)
		if r.DealtBy > 0 {
			userIds = append(userIds, r.DealtBy)
		}
		if target := targetMap[r.Id]; target.Author.Id > 0 {
			userIds = append(userIds, target.Author.Id)
		}
	}
	userMap := s.loadUserSummaryMap(ctx, userIds)

	var list []modout.ReportItem
	for _, r := range records {
		target := targetMap[r.Id]
		if target.Author.Id > 0 {
			target.Author = s.userSummary(userMap, target.Author.Id)
		}

		list = append(list, modout.ReportItem{
			Id:           r.Id,
			ReporterId:   r.ReporterId,
			Reporter:     s.userSummary(userMap, r.ReporterId),
			TargetType:   r.TargetType,
			TargetId:     r.TargetId,
			Target:       target,
			Reason:       r.Reason,
			Status:       r.Status,
			DealtBy:      r.DealtBy,
			DealtUser:    s.userSummary(userMap, r.DealtBy),
			DealtComment: r.DealtComment,
			DealtAt:      r.DealtAt,
			CreatedAt:    r.CreatedAt,
		})
	}

	return &modout.ListReportsOut{
		Page:  in.Page,
		Size:  in.Size,
		Total: total,
		List:  list,
	}, nil
}

func (s *sModReportUsecase) Resolve(ctx context.Context, actor *model.Actor, in modin.ResolveReportInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	report, err := service.ModReportDomain().GetById(ctx, in.Id)
	if err != nil || report == nil {
		return gerror.New(gi18n.T(ctx, "moderation.report.not_found"))
	}

	if report.Status != consts.ModReportStatusPending {
		return gerror.New(gi18n.T(ctx, "moderation.report.already_resolved"))
	}

	err = service.ModReportDomain().Update(ctx, in.Id, do.ModReport{
		Status:       in.Status,
		DealtBy:      actor.Id,
		DealtComment: in.Comment,
		DealtAt:      gtime.Now(),
	})
	if err != nil {
		return err
	}

	titleKey := "site.message.report.accepted.title"
	contentKey := "site.message.report.accepted.content"
	if in.Status == consts.ModReportStatusRejected {
		titleKey = "site.message.report.rejected.title"
		contentKey = "site.message.report.rejected.content"
	}
	content := strings.TrimSpace(in.Comment)
	target := s.loadReportTargetSummaryMap(ctx, []entity.ModReport{*report})[report.Id]
	targetType, targetId := target.SiteMessageTarget()
	notifyInp := sitein.MessageNotifyInp{
		SenderId:   actor.Id,
		ReceiverId: report.ReporterId,
		TitleKey:   titleKey,
		Content:    content,
		TargetType: targetType,
		TargetId:   targetId,
	}
	if content == "" {
		notifyInp.ContentKey = contentKey
	}
	service.SiteMessageUsecase().Notify(ctx, notifyInp)
	return nil
}

func (s *sModReportUsecase) loadReportTargetSummaryMap(ctx context.Context, records []entity.ModReport) map[uint64]modout.ReportTargetSummary {
	targets := make(map[uint64]modout.ReportTargetSummary, len(records))
	grouped := s.groupReportTargets(records)

	s.loadCatalogTorrentTargets(ctx, grouped[consts.ModReportTargetTypeCatalogTorrent], targets)
	s.loadCatalogCommentTargets(ctx, grouped[consts.ModReportTargetTypeCatalogComment], targets)
	s.loadCatalogSubtitleTargets(ctx, grouped[consts.ModReportTargetTypeCatalogSubtitle], targets)
	s.loadForumTopicTargets(ctx, grouped[consts.ModReportTargetTypeForumTopic], targets)
	s.loadForumReplyTargets(ctx, grouped[consts.ModReportTargetTypeForumReply], targets)

	for _, record := range records {
		if _, ok := targets[record.Id]; !ok {
			targets[record.Id] = s.deletedReportTargetSummary(record)
		}
	}
	return targets
}

func (s *sModReportUsecase) groupReportTargets(records []entity.ModReport) map[string]map[uint64][]uint64 {
	grouped := make(map[string]map[uint64][]uint64)
	for _, record := range records {
		if record.TargetType == "" || record.TargetId == 0 {
			continue
		}
		if grouped[record.TargetType] == nil {
			grouped[record.TargetType] = make(map[uint64][]uint64)
		}
		grouped[record.TargetType][record.TargetId] = append(grouped[record.TargetType][record.TargetId], record.Id)
	}
	return grouped
}

func (s *sModReportUsecase) loadCatalogTorrentTargets(ctx context.Context, reportIdsByTarget map[uint64][]uint64, targets map[uint64]modout.ReportTargetSummary) {
	if len(reportIdsByTarget) == 0 {
		return
	}
	torrents, err := service.CatalogTorrentDomain().GetTorrentsByIds(ctx, s.reportTargetIds(reportIdsByTarget))
	if err != nil {
		s.assignUnavailableTargetSummaries(reportIdsByTarget, targets, consts.ModReportTargetTypeCatalogTorrent)
		return
	}
	for _, torrent := range torrents {
		if torrent == nil {
			continue
		}
		s.assignReportTargetSummary(reportIdsByTarget[torrent.Id], targets, modout.ReportTargetSummary{
			Type:   consts.ModReportTargetTypeCatalogTorrent,
			Id:     torrent.Id,
			Title:  torrent.Name,
			Author: model.IamUserSummary{Id: torrent.OwnerId},
		})
	}
}

func (s *sModReportUsecase) loadCatalogCommentTargets(ctx context.Context, reportIdsByTarget map[uint64][]uint64, targets map[uint64]modout.ReportTargetSummary) {
	if len(reportIdsByTarget) == 0 {
		return
	}
	comments, err := service.CatalogCommentDomain().GetCommentsByIds(ctx, s.reportTargetIds(reportIdsByTarget))
	if err != nil {
		s.assignUnavailableTargetSummaries(reportIdsByTarget, targets, consts.ModReportTargetTypeCatalogComment)
		return
	}
	for _, comment := range comments {
		parentType := ""
		if comment.TargetType == consts.CatalogCommentTargetTypeCatalogTorrent && comment.TargetId > 0 {
			parentType = consts.ModReportTargetTypeCatalogTorrent
		}
		s.assignReportTargetSummary(reportIdsByTarget[comment.Id], targets, modout.ReportTargetSummary{
			Type:       consts.ModReportTargetTypeCatalogComment,
			Id:         comment.Id,
			Title:      s.contentSnippet(comment.Content),
			ParentType: parentType,
			ParentId:   comment.TargetId,
			Author:     model.IamUserSummary{Id: comment.UserId},
		})
	}
}

func (s *sModReportUsecase) loadCatalogSubtitleTargets(ctx context.Context, reportIdsByTarget map[uint64][]uint64, targets map[uint64]modout.ReportTargetSummary) {
	if len(reportIdsByTarget) == 0 {
		return
	}
	subtitles, err := service.CatalogSubtitleDomain().GetSubtitlesByIds(ctx, s.reportTargetIds(reportIdsByTarget))
	if err != nil {
		s.assignUnavailableTargetSummaries(reportIdsByTarget, targets, consts.ModReportTargetTypeCatalogSubtitle)
		return
	}
	for _, subtitle := range subtitles {
		title := strings.TrimSpace(subtitle.Title)
		if title == "" {
			title = strings.TrimSpace(subtitle.FileName)
		}
		s.assignReportTargetSummary(reportIdsByTarget[subtitle.Id], targets, modout.ReportTargetSummary{
			Type:       consts.ModReportTargetTypeCatalogSubtitle,
			Id:         subtitle.Id,
			Title:      title,
			ParentType: consts.ModReportTargetTypeCatalogTorrent,
			ParentId:   subtitle.TorrentId,
			Author:     model.IamUserSummary{Id: subtitle.UserId},
		})
	}
}

func (s *sModReportUsecase) loadForumTopicTargets(ctx context.Context, reportIdsByTarget map[uint64][]uint64, targets map[uint64]modout.ReportTargetSummary) {
	if len(reportIdsByTarget) == 0 {
		return
	}
	topics, err := service.ForumTopicDomain().GetTopicsByIds(ctx, s.reportTargetIds(reportIdsByTarget))
	if err != nil {
		s.assignUnavailableTargetSummaries(reportIdsByTarget, targets, consts.ModReportTargetTypeForumTopic)
		return
	}
	for _, topic := range topics {
		s.assignReportTargetSummary(reportIdsByTarget[topic.Id], targets, modout.ReportTargetSummary{
			Type:   consts.ModReportTargetTypeForumTopic,
			Id:     topic.Id,
			Title:  topic.Subject,
			Author: model.IamUserSummary{Id: topic.UserId},
		})
	}
}

func (s *sModReportUsecase) loadForumReplyTargets(ctx context.Context, reportIdsByTarget map[uint64][]uint64, targets map[uint64]modout.ReportTargetSummary) {
	if len(reportIdsByTarget) == 0 {
		return
	}
	replies, err := service.ForumReplyDomain().GetRepliesByIds(ctx, s.reportTargetIds(reportIdsByTarget))
	if err != nil {
		s.assignUnavailableTargetSummaries(reportIdsByTarget, targets, consts.ModReportTargetTypeForumReply)
		return
	}
	for _, reply := range replies {
		s.assignReportTargetSummary(reportIdsByTarget[reply.Id], targets, modout.ReportTargetSummary{
			Type:       consts.ModReportTargetTypeForumReply,
			Id:         reply.Id,
			Title:      s.contentSnippet(reply.Content),
			ParentType: consts.ModReportTargetTypeForumTopic,
			ParentId:   reply.TopicId,
			Author:     model.IamUserSummary{Id: reply.UserId},
		})
	}
}

func (s *sModReportUsecase) assignUnavailableTargetSummaries(reportIdsByTarget map[uint64][]uint64, targets map[uint64]modout.ReportTargetSummary, targetType string) {
	for targetId, reportIds := range reportIdsByTarget {
		s.assignReportTargetSummary(reportIds, targets, modout.ReportTargetSummary{
			Type:   targetType,
			Id:     targetId,
			Status: consts.ModReportTargetStatusUnavailable,
		})
	}
}

func (s *sModReportUsecase) assignReportTargetSummary(reportIds []uint64, targets map[uint64]modout.ReportTargetSummary, summary modout.ReportTargetSummary) {
	if summary.Status == "" {
		summary.Status = consts.ModReportTargetStatusNormal
	}
	if strings.TrimSpace(summary.Title) == "" {
		summary.Title = fmt.Sprintf("#%d", summary.Id)
	}
	for _, reportId := range reportIds {
		targets[reportId] = summary
	}
}

func (s *sModReportUsecase) reportTargetIds(reportIdsByTarget map[uint64][]uint64) []uint64 {
	ids := make([]uint64, 0, len(reportIdsByTarget))
	for id := range reportIdsByTarget {
		if id == 0 {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

func (s *sModReportUsecase) deletedReportTargetSummary(record entity.ModReport) modout.ReportTargetSummary {
	return modout.ReportTargetSummary{
		Type:   record.TargetType,
		Id:     record.TargetId,
		Status: consts.ModReportTargetStatusDeleted,
	}
}

func (s *sModReportUsecase) loadUserSummaryMap(ctx context.Context, userIds []uint64) map[uint64]model.IamUserSummary {
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

func (s *sModReportUsecase) uniqueUserIds(userIds []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(userIds))
	uniqueIds := make([]uint64, 0, len(userIds))
	for _, id := range userIds {
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

func (s *sModReportUsecase) userSummary(userMap map[uint64]model.IamUserSummary, userId uint64) model.IamUserSummary {
	if userId == 0 {
		return model.IamUserSummary{}
	}
	if summary, ok := userMap[userId]; ok {
		return summary
	}
	return model.IamUserSummary{Id: userId}
}

func (s *sModReportUsecase) contentSnippet(content string) string {
	normalized := strings.Join(strings.Fields(content), " ")
	runes := []rune(normalized)
	if len(runes) <= 80 {
		return normalized
	}
	return string(runes[:80]) + "..."
}
