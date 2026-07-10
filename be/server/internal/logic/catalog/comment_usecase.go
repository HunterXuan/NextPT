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

func init() {
	service.RegisterCatalogCommentUsecase(NewCatalogCommentUsecase())
}

func NewCatalogCommentUsecase() *sCatalogCommentUsecase {
	return &sCatalogCommentUsecase{}
}

func (s *sCatalogCommentUsecase) Create(ctx context.Context, actor *model.Actor, in catalogin.CommentCreateInp) (*catalogout.CommentCreateOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}

	// Ensure torrent exists and is visible
	torrent, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return nil, err
	}

	id, err := service.CatalogCommentDomain().CreateComment(ctx, consts.CatalogCommentTargetTypeCatalogTorrent, in.Id, actor.Id, in.Content)
	if err != nil {
		return nil, err
	}
	service.SiteMessageUsecase().Notify(ctx, sitein.MessageNotifyInp{
		ActorId:    actor.Id,
		ReceiverId: torrent.OwnerId,
		TitleKey:   "site.message.catalog_comment.title",
		TitleArgs:  []any{torrent.Name},
		Content:    in.Content,
		TargetType: consts.SiteMessageTargetTypeCatalogTorrent,
		TargetId:   in.Id,
	})

	return &catalogout.CommentCreateOut{Id: id}, nil
}

func (s *sCatalogCommentUsecase) List(ctx context.Context, actor *model.Actor, in catalogin.CommentListInp) (*catalogout.CommentListOut, error) {
	// Ensure torrent exists and is visible
	_, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return nil, err
	}

	comments, total, err := service.CatalogCommentDomain().QueryCommentsByTarget(ctx, consts.CatalogCommentTargetTypeCatalogTorrent, in.Id, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	var list []catalogout.CommentListItem
	if len(comments) == 0 {
		return &catalogout.CommentListOut{List: list, Total: total}, nil
	}

	userIds := make([]uint64, 0, len(comments))
	for _, c := range comments {
		userIds = append(userIds, c.UserId)
	}
	authorMap := s.loadCommentAuthorMap(ctx, userIds)

	likedMap := make(map[uint64]bool)
	if actor != nil && actor.Id > 0 {
		var cids []uint64
		for _, c := range comments {
			cids = append(cids, c.Id)
		}
		var likes []entity.CatalogCommentLike
		likes, _ = service.CatalogCommentDomain().GetCommentLikesByUser(ctx, actor.Id, cids)
		for _, l := range likes {
			likedMap[l.CommentId] = true
		}
	}

	for _, c := range comments {
		list = append(list, catalogout.CommentListItem{
			Id:          c.Id,
			Author:      authorMap[c.UserId],
			Content:     c.Content,
			LikeCount:   c.LikeCount,
			RewardCount: c.RewardCount,
			CreatedAt:   c.CreatedAt.String(),
			IsLiked:     likedMap[c.Id],
		})
	}

	return &catalogout.CommentListOut{
		List:  list,
		Total: total,
	}, nil
}

func (s *sCatalogCommentUsecase) loadCommentAuthorMap(ctx context.Context, userIds []uint64) map[uint64]model.IamUserSummary {
	authorMap := make(map[uint64]model.IamUserSummary)
	if len(userIds) == 0 {
		return authorMap
	}

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
			author.Id = user.Id
			author.Username = user.Username
			authorMap[user.Id] = author
		}
	}

	profiles, err := service.IamUserDomain().GetUserProfilesByUserIds(ctx, uniqueIds)
	if err == nil {
		for _, profile := range profiles {
			author := authorMap[profile.UserId]
			author.Id = profile.UserId
			author.Avatar = profile.Avatar
			authorMap[profile.UserId] = author
		}
	}
	return authorMap
}

func (s *sCatalogCommentUsecase) ToggleLike(ctx context.Context, actor *model.Actor, in catalogin.CommentToggleLikeInp) (*catalogout.CommentToggleLikeOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	// Verify comment exists and belongs to target
	comment, err := service.CatalogCommentDomain().GetCommentById(ctx, in.Cid)
	if err != nil {
		return nil, err
	}
	if comment.TargetType != consts.CatalogCommentTargetTypeCatalogTorrent || comment.TargetId != in.Id {
		return nil, gerror.New(gi18n.T(ctx, "catalog.comment.not_found"))
	}

	// Verify target torrent is visible to the actor
	_, err = service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, comment.TargetId)
	if err != nil {
		return nil, err
	}

	var isLiked bool
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var txErr error
		isLiked, txErr = service.CatalogCommentDomain().ToggleLike(ctx, actor.Id, in.Cid)
		return txErr
	})

	if err != nil {
		return nil, err
	}
	return &catalogout.CommentToggleLikeOut{IsLiked: isLiked}, nil
}

func (s *sCatalogCommentUsecase) Reward(ctx context.Context, actor *model.Actor, in catalogin.CommentRewardInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	// Verify comment exists
	comment, err := service.CatalogCommentDomain().GetCommentById(ctx, in.Cid)
	if err != nil {
		return err
	}
	if comment.TargetType != consts.CatalogCommentTargetTypeCatalogTorrent || comment.TargetId != in.Id {
		return gerror.New(gi18n.T(ctx, "catalog.comment.not_found"))
	}

	// Verify target torrent is visible to the actor
	_, err = service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, comment.TargetId)
	if err != nil {
		return err
	}

	if comment.UserId == actor.Id {
		return gerror.New(gi18n.T(ctx, "catalog.comment.reward_self_denied"))
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err := service.EconomyBonusUsecase().TransferBonus(ctx, actor.Id, comment.UserId, in.Amount, consts.EconomyBonusTargetTypeCatalogComment, comment.Id)
		if err != nil {
			return err
		}

		if err = service.EconomyRewardDomain().InsertRewardRecord(ctx, entity.EconomyRewardRecord{
			TargetType: consts.EconomyBonusTargetTypeCatalogComment,
			TargetId:   comment.Id,
			FromUserId: actor.Id,
			ToUserId:   comment.UserId,
			Amount:     in.Amount,
		}); err != nil {
			return err
		}

		err = service.CatalogCommentDomain().IncrementCommentRewardStats(ctx, comment.Id)
		return err
	})
	if err != nil {
		return err
	}

	service.SiteMessageUsecase().Notify(ctx, sitein.MessageNotifyInp{
		ActorId:     actor.Id,
		ReceiverId:  comment.UserId,
		TitleKey:    "site.message.reward.catalog_comment.title",
		ContentKey:  "site.message.reward.content",
		ContentArgs: []any{in.Amount},
		TargetType:  consts.SiteMessageTargetTypeCatalogTorrent,
		TargetId:    comment.TargetId,
	})
	return nil
}

func (s *sCatalogCommentUsecase) Report(ctx context.Context, actor *model.Actor, in catalogin.CommentReportInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}

	comment, err := service.CatalogCommentDomain().GetCommentById(ctx, in.Cid)
	if err != nil {
		return err
	}
	if comment.TargetType != consts.CatalogCommentTargetTypeCatalogTorrent || comment.TargetId != in.Id {
		return gerror.New(gi18n.T(ctx, "catalog.comment.not_found"))
	}

	_, err = service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, comment.TargetId)
	if err != nil {
		return err
	}

	return service.ModReportUsecase().Create(ctx, actor, modin.CreateReportInp{
		TargetType: consts.ModReportTargetTypeCatalogComment,
		TargetId:   in.Cid,
		Reason:     in.Reason,
	})
}
