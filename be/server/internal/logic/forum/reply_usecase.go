package forum

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/forumin"
	"server/internal/model/in/modin"
	"server/internal/model/out/forumout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

type sForumReplyUsecase struct{}

func init() {
	service.RegisterForumReplyUsecase(NewForumReplyUsecase())
}

func NewForumReplyUsecase() *sForumReplyUsecase {
	return &sForumReplyUsecase{}
}

func (s *sForumReplyUsecase) List(ctx context.Context, actor *model.Actor, in forumin.ReplyListInp) (*forumout.ReplyListOut, error) {
	// 1. Get Topic
	topic, err := service.ForumTopicDomain().GetTopicById(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 2. Get Node
	nodePtr, err := service.ForumNodeDomain().GetNodeById(ctx, topic.NodeId)
	if err != nil || nodePtr == nil || nodePtr.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "forum.node.not_found"))
	}
	node := *nodePtr

	// 3. Check Policy
	if err := service.ForumNodeDomain().CheckNodeReadPolicy(ctx, actor, &node); err != nil {
		return nil, err
	}

	// 4. Query
	replies, total, err := service.ForumReplyDomain().QueryRepliesByTopic(ctx, in.Id, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	// 5. Assemble
	var list []forumout.ReplyListItem
	if len(replies) == 0 {
		return &forumout.ReplyListOut{List: list, Total: total}, nil
	}

	userIds := make([]uint64, 0, len(replies))
	for _, r := range replies {
		userIds = append(userIds, r.UserId)
	}
	authorMap := s.loadUserSummaryMap(ctx, userIds)

	likedMap := make(map[uint64]bool)
	if actor != nil && actor.Id > 0 {
		var replyIds []uint64
		for _, r := range replies {
			replyIds = append(replyIds, r.Id)
		}
		likes, _ := service.ForumReplyDomain().GetReplyLikesByUser(ctx, actor.Id, replyIds)
		for _, l := range likes {
			likedMap[l.ReplyId] = true
		}
	}

	for _, r := range replies {
		list = append(list, forumout.ReplyListItem{
			Id:          r.Id,
			Author:      authorMap[r.UserId],
			Content:     r.Content,
			CreatedAt:   s.formatTime(r.CreatedAt),
			LikeCount:   r.LikeCount,
			RewardCount: r.RewardCount,
			IsLiked:     likedMap[r.Id],
		})
	}

	return &forumout.ReplyListOut{
		List:  list,
		Total: total,
	}, nil
}

func (s *sForumReplyUsecase) formatTime(value *gtime.Time) string {
	if value == nil {
		return ""
	}
	return value.String()
}

func (s *sForumReplyUsecase) loadUserSummaryMap(ctx context.Context, userIds []uint64) map[uint64]model.IamUserSummary {
	userMap := make(map[uint64]model.IamUserSummary)
	if len(userIds) == 0 {
		return userMap
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

func (s *sForumReplyUsecase) Create(ctx context.Context, actor *model.Actor, in forumin.ReplyCreateInp) (uint64, error) {
	if actor == nil {
		return 0, gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}

	topic, err := service.ForumTopicDomain().GetTopicById(ctx, in.Id)
	if err != nil {
		return 0, err
	}

	if err := service.ForumTopicDomain().CheckTopicWritePolicy(ctx, actor, topic); err != nil {
		return 0, err
	}

	nodePtr, err := service.ForumNodeDomain().GetNodeById(ctx, topic.NodeId)
	if err != nil || nodePtr == nil || nodePtr.Id == 0 {
		return 0, gerror.New(gi18n.T(ctx, "forum.node.not_found"))
	}
	node := *nodePtr

	if err := service.ForumNodeDomain().CheckNodeWritePolicy(ctx, actor, &node); err != nil {
		return 0, err
	}

	var replyId uint64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. Insert Reply
		id, err := service.ForumReplyDomain().InsertReply(ctx, actor, in)
		if err != nil {
			return err
		}
		replyId = id

		// 2. Update Topic Stats
		if err := service.ForumTopicDomain().UpdateTopicReplyStats(ctx, in.Id, replyId, actor.Id); err != nil {
			return err
		}

		// 3. Update Node Stats
		if err := service.ForumNodeDomain().UpdateStats(ctx, topic.NodeId, 0, 1); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, gerror.Wrap(err, gi18n.T(ctx, "forum.reply.create_failed"))
	}
	return replyId, nil
}

func (s *sForumReplyUsecase) ToggleReplyLike(ctx context.Context, actor *model.Actor, in forumin.ReplyToggleLikeInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}

	reply, err := service.ForumReplyDomain().GetReplyById(ctx, in.Id)
	if err != nil {
		return err
	}

	if err := s.checkReplyReadPolicy(ctx, actor, reply); err != nil {
		return err
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, txErr := service.ForumReplyDomain().ToggleLike(ctx, actor, in.Id)
		return txErr
	})
	return err
}

func (s *sForumReplyUsecase) checkReplyReadPolicy(ctx context.Context, actor *model.Actor, reply *entity.ForumReply) error {
	topic, err := service.ForumTopicDomain().GetTopicById(ctx, reply.TopicId)
	if err != nil {
		return err
	} else if topic == nil {
		return gerror.New(gi18n.T(ctx, "forum.topic.not_found"))
	}

	nodePtr, err := service.ForumNodeDomain().GetNodeById(ctx, topic.NodeId)
	if err != nil || nodePtr == nil || nodePtr.Id == 0 {
		return gerror.New(gi18n.T(ctx, "forum.node.not_found"))
	}
	node := *nodePtr
	return service.ForumNodeDomain().CheckNodeReadPolicy(ctx, actor, &node)
}

func (s *sForumReplyUsecase) RewardReply(ctx context.Context, actor *model.Actor, in forumin.ReplyRewardInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}

	reply, err := service.ForumReplyDomain().GetReplyById(ctx, in.Id)
	if err != nil {
		return err
	}

	if err := s.checkReplyReadPolicy(ctx, actor, reply); err != nil {
		return err
	}

	if reply.UserId == actor.Id {
		return gerror.New(gi18n.T(ctx, "forum.reply.reward_self_not_allowed"))
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := service.EconomyBonusUsecase().TransferBonus(ctx, actor.Id, reply.UserId, in.Amount, consts.EconomyBonusTargetTypeForumReply, in.Id, "Reward reply", "Reply rewarded"); err != nil {
			return err
		}
		if err := service.EconomyRewardDomain().InsertRewardRecord(ctx, entity.EconomyRewardRecord{
			TargetType: consts.EconomyBonusTargetTypeForumReply,
			TargetId:   in.Id,
			FromUserId: actor.Id,
			ToUserId:   reply.UserId,
			Amount:     in.Amount,
		}); err != nil {
			return err
		}
		return service.ForumReplyDomain().IncrementRewardStats(ctx, in.Id)
	})
}

func (s *sForumReplyUsecase) ReportReply(ctx context.Context, actor *model.Actor, in forumin.ReplyReportInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}
	reply, err := service.ForumReplyDomain().GetReplyById(ctx, in.Id)
	if err != nil {
		return err
	}

	if err := s.checkReplyReadPolicy(ctx, actor, reply); err != nil {
		return err
	}

	return service.ModReportUsecase().Create(ctx, actor, modin.CreateReportInp{
		TargetType: "forum_reply",
		TargetId:   in.Id,
		Reason:     in.Reason,
	})
}
