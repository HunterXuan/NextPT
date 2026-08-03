package catalog

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/catalogin"
	"server/internal/model/in/modin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/catalogout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sCatalogCommentUsecase struct{}

type catalogCommentTarget struct {
	Type              string
	Id                uint64
	OwnerId           uint64
	Title             string
	MessageTargetType string
}

func init() {
	service.RegisterCatalogCommentUsecase(NewCatalogCommentUsecase())
}

func NewCatalogCommentUsecase() *sCatalogCommentUsecase {
	return &sCatalogCommentUsecase{}
}

func (s *sCatalogCommentUsecase) Create(ctx context.Context, actor *model.Actor, in catalogin.CommentCreateInp) (*catalogout.CommentCreateOut, error) {
	return s.createForTarget(ctx, actor, consts.CatalogCommentTargetTypeCatalogTorrent, in.Id, in.Content)
}

func (s *sCatalogCommentUsecase) CreateForRequest(ctx context.Context, actor *model.Actor, in catalogin.RequestCommentCreateInp) (*catalogout.CommentCreateOut, error) {
	return s.createForTarget(ctx, actor, consts.CatalogCommentTargetTypeCatalogRequest, in.Id, in.Content)
}

func (s *sCatalogCommentUsecase) createForTarget(ctx context.Context, actor *model.Actor, targetType string, targetId uint64, content string) (*catalogout.CommentCreateOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	target, err := s.loadTarget(ctx, actor, targetType, targetId)
	if err != nil {
		return nil, err
	}
	id, err := service.CatalogCommentDomain().CreateComment(ctx, target.Type, target.Id, actor.Id, content)
	if err != nil {
		return nil, err
	}
	titleKey := "site.message.catalog_comment.title"
	if target.Type == consts.CatalogCommentTargetTypeCatalogRequest {
		titleKey = "site.message.catalog_request_comment.title"
	}
	service.SiteMessageUsecase().Notify(ctx, sitein.MessageNotifyInp{
		ActorId:    actor.Id,
		ReceiverId: target.OwnerId,
		TitleKey:   titleKey,
		TitleArgs:  []any{target.Title},
		Content:    content,
		TargetType: target.MessageTargetType,
		TargetId:   target.Id,
	})
	return &catalogout.CommentCreateOut{Id: id}, nil
}

func (s *sCatalogCommentUsecase) List(ctx context.Context, actor *model.Actor, in catalogin.CommentListInp) (*catalogout.CommentListOut, error) {
	return s.listForTarget(ctx, actor, consts.CatalogCommentTargetTypeCatalogTorrent, in.Id, in.Page, in.Size)
}

func (s *sCatalogCommentUsecase) ListForRequest(ctx context.Context, actor *model.Actor, in catalogin.RequestCommentListInp) (*catalogout.CommentListOut, error) {
	return s.listForTarget(ctx, actor, consts.CatalogCommentTargetTypeCatalogRequest, in.Id, in.Page, in.Size)
}

func (s *sCatalogCommentUsecase) listForTarget(ctx context.Context, actor *model.Actor, targetType string, targetId uint64, page int, size int) (*catalogout.CommentListOut, error) {
	if _, err := s.loadTarget(ctx, actor, targetType, targetId); err != nil {
		return nil, err
	}
	comments, total, err := service.CatalogCommentDomain().QueryCommentsByTarget(ctx, targetType, targetId, page, size)
	if err != nil {
		return nil, err
	}
	list := make([]catalogout.CommentListItem, 0, len(comments))
	if len(comments) == 0 {
		return &catalogout.CommentListOut{List: list, Total: total}, nil
	}

	userIds := make([]uint64, 0, len(comments))
	commentIds := make([]uint64, 0, len(comments))
	for _, comment := range comments {
		userIds = append(userIds, comment.UserId)
		commentIds = append(commentIds, comment.Id)
	}
	authorMap := s.loadCommentAuthorMap(ctx, userIds)
	likedMap := s.loadCommentLikedMap(ctx, actor, commentIds)
	for _, comment := range comments {
		list = append(list, catalogout.CommentListItem{
			Id:          comment.Id,
			Author:      authorMap[comment.UserId],
			Content:     comment.Content,
			LikeCount:   comment.LikeCount,
			RewardCount: comment.RewardCount,
			CreatedAt:   comment.CreatedAt.String(),
			IsLiked:     likedMap[comment.Id],
		})
	}
	return &catalogout.CommentListOut{List: list, Total: total}, nil
}

func (s *sCatalogCommentUsecase) ToggleLike(ctx context.Context, actor *model.Actor, in catalogin.CommentToggleLikeInp) (*catalogout.CommentToggleLikeOut, error) {
	return s.toggleLikeForTarget(ctx, actor, consts.CatalogCommentTargetTypeCatalogTorrent, in.Id, in.Cid)
}

func (s *sCatalogCommentUsecase) ToggleLikeForRequest(ctx context.Context, actor *model.Actor, in catalogin.RequestCommentActionInp) (*catalogout.CommentToggleLikeOut, error) {
	return s.toggleLikeForTarget(ctx, actor, consts.CatalogCommentTargetTypeCatalogRequest, in.Id, in.Cid)
}

func (s *sCatalogCommentUsecase) toggleLikeForTarget(ctx context.Context, actor *model.Actor, targetType string, targetId uint64, commentId uint64) (*catalogout.CommentToggleLikeOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	if _, _, err := s.loadTargetComment(ctx, actor, targetType, targetId, commentId); err != nil {
		return nil, err
	}
	var isLiked bool
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var toggleErr error
		isLiked, toggleErr = service.CatalogCommentDomain().ToggleLike(ctx, actor.Id, commentId)
		return toggleErr
	})
	if err != nil {
		return nil, err
	}
	return &catalogout.CommentToggleLikeOut{IsLiked: isLiked}, nil
}

func (s *sCatalogCommentUsecase) Reward(ctx context.Context, actor *model.Actor, in catalogin.CommentRewardInp) error {
	return s.rewardForTarget(ctx, actor, consts.CatalogCommentTargetTypeCatalogTorrent, in.Id, in.Cid, in.Amount)
}

func (s *sCatalogCommentUsecase) RewardForRequest(ctx context.Context, actor *model.Actor, in catalogin.RequestCommentRewardInp) error {
	return s.rewardForTarget(ctx, actor, consts.CatalogCommentTargetTypeCatalogRequest, in.Id, in.Cid, in.Amount)
}

func (s *sCatalogCommentUsecase) rewardForTarget(ctx context.Context, actor *model.Actor, targetType string, targetId uint64, commentId uint64, amount float64) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	target, comment, err := s.loadTargetComment(ctx, actor, targetType, targetId, commentId)
	if err != nil {
		return err
	}
	if comment.UserId == actor.Id {
		return gerror.New(gi18n.T(ctx, "catalog.comment.reward_self_denied"))
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if transferErr := service.EconomyBonusUsecase().TransferBonus(ctx, actor.Id, comment.UserId, amount, consts.EconomyBonusTargetTypeCatalogComment, comment.Id); transferErr != nil {
			return transferErr
		}
		if insertErr := service.EconomyRewardDomain().InsertRewardRecord(ctx, entity.EconomyRewardRecord{
			TargetType: consts.EconomyBonusTargetTypeCatalogComment,
			TargetId:   comment.Id,
			FromUserId: actor.Id,
			ToUserId:   comment.UserId,
			Amount:     amount,
		}); insertErr != nil {
			return insertErr
		}
		return service.CatalogCommentDomain().IncrementCommentRewardStats(ctx, comment.Id)
	})
	if err != nil {
		return err
	}
	service.SiteMessageUsecase().Notify(ctx, sitein.MessageNotifyInp{
		ActorId:     actor.Id,
		ReceiverId:  comment.UserId,
		TitleKey:    "site.message.reward.catalog_comment.title",
		ContentKey:  "site.message.reward.content",
		ContentArgs: []any{amount},
		TargetType:  target.MessageTargetType,
		TargetId:    target.Id,
	})
	return nil
}

func (s *sCatalogCommentUsecase) Report(ctx context.Context, actor *model.Actor, in catalogin.CommentReportInp) error {
	return s.reportForTarget(ctx, actor, consts.CatalogCommentTargetTypeCatalogTorrent, in.Id, in.Cid, in.Reason)
}

func (s *sCatalogCommentUsecase) ReportForRequest(ctx context.Context, actor *model.Actor, in catalogin.RequestCommentReportInp) error {
	return s.reportForTarget(ctx, actor, consts.CatalogCommentTargetTypeCatalogRequest, in.Id, in.Cid, in.Reason)
}

func (s *sCatalogCommentUsecase) reportForTarget(ctx context.Context, actor *model.Actor, targetType string, targetId uint64, commentId uint64, reason string) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	if _, _, err := s.loadTargetComment(ctx, actor, targetType, targetId, commentId); err != nil {
		return err
	}
	return service.ModReportUsecase().Create(ctx, actor, modin.CreateReportInp{
		TargetType: consts.ModReportTargetTypeCatalogComment,
		TargetId:   commentId,
		Reason:     reason,
	})
}

func (s *sCatalogCommentUsecase) loadTargetComment(ctx context.Context, actor *model.Actor, targetType string, targetId uint64, commentId uint64) (*catalogCommentTarget, *entity.CatalogComment, error) {
	target, err := s.loadTarget(ctx, actor, targetType, targetId)
	if err != nil {
		return nil, nil, err
	}
	comment, err := service.CatalogCommentDomain().GetCommentById(ctx, commentId)
	if err != nil {
		return nil, nil, err
	}
	if comment.TargetType != target.Type || comment.TargetId != target.Id {
		return nil, nil, gerror.New(gi18n.T(ctx, "catalog.comment.not_found"))
	}
	return target, comment, nil
}

func (s *sCatalogCommentUsecase) loadTarget(ctx context.Context, actor *model.Actor, targetType string, targetId uint64) (*catalogCommentTarget, error) {
	switch targetType {
	case consts.CatalogCommentTargetTypeCatalogTorrent:
		torrent, err := service.CatalogTorrentDomain().LoadViewableTorrent(ctx, actor, targetId)
		if err != nil {
			return nil, err
		}
		return &catalogCommentTarget{
			Type:              targetType,
			Id:                torrent.Id,
			OwnerId:           torrent.OwnerId,
			Title:             torrent.Name,
			MessageTargetType: consts.SiteMessageTargetTypeCatalogTorrent,
		}, nil
	case consts.CatalogCommentTargetTypeCatalogRequest:
		request, err := service.CatalogRequestDomain().GetRequestById(ctx, targetId)
		if err != nil {
			return nil, err
		}
		return &catalogCommentTarget{
			Type:              targetType,
			Id:                request.Id,
			OwnerId:           request.RequesterId,
			Title:             request.Title,
			MessageTargetType: consts.SiteMessageTargetTypeCatalogRequest,
		}, nil
	default:
		return nil, gerror.New(gi18n.T(ctx, "catalog.comment.not_found"))
	}
}

func (s *sCatalogCommentUsecase) loadCommentLikedMap(ctx context.Context, actor *model.Actor, commentIds []uint64) map[uint64]bool {
	likedMap := make(map[uint64]bool)
	if actor == nil || actor.Id == 0 || len(commentIds) == 0 {
		return likedMap
	}
	likes, err := service.CatalogCommentDomain().GetCommentLikesByUser(ctx, actor.Id, commentIds)
	if err != nil {
		return likedMap
	}
	for _, like := range likes {
		likedMap[like.CommentId] = true
	}
	return likedMap
}

func (s *sCatalogCommentUsecase) loadCommentAuthorMap(ctx context.Context, userIds []uint64) map[uint64]model.IamUserSummary {
	authorMap := make(map[uint64]model.IamUserSummary)
	uniqueIds := make([]uint64, 0, len(userIds))
	seen := make(map[uint64]struct{}, len(userIds))
	for _, id := range userIds {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIds = append(uniqueIds, id)
		authorMap[id] = model.IamUserSummary{Id: id}
	}
	if len(uniqueIds) == 0 {
		return authorMap
	}
	users, err := service.IamUserDomain().GetUsersByIds(ctx, uniqueIds)
	if err == nil {
		for _, user := range users {
			author := authorMap[user.Id]
			author.Username = user.Username
			authorMap[user.Id] = author
		}
	}
	profiles, err := service.IamUserDomain().GetUserProfilesByUserIds(ctx, uniqueIds)
	if err == nil {
		for _, profile := range profiles {
			author := authorMap[profile.UserId]
			author.Avatar = profile.Avatar
			authorMap[profile.UserId] = author
		}
	}
	return authorMap
}
