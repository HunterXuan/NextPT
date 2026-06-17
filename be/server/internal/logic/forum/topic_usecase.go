package forum

import (
	"context"
	"fmt"

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
)

type sForumTopicUsecase struct{}

func init() {
	service.RegisterForumTopicUsecase(NewForumTopicUsecase())
}

func NewForumTopicUsecase() *sForumTopicUsecase {
	return &sForumTopicUsecase{}
}

func (s *sForumTopicUsecase) List(ctx context.Context, actor *model.Actor, in forumin.TopicListInp) (*forumout.TopicListOut, error) {
	nodePtr, err := service.ForumNodeDomain().GetNodeBySlug(ctx, in.Slug)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "forum.topic.query_node_failed"))
	}
	if nodePtr == nil || nodePtr.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "forum.node.not_found"))
	}
	node := *nodePtr

	if err := service.ForumNodeDomain().CheckNodeReadPolicy(ctx, actor, &node); err != nil {
		return nil, err
	}

	topics, total, err := service.ForumTopicDomain().QueryTopicsByNode(ctx, node.Id, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	userIds := make([]uint64, 0, len(topics))
	for _, t := range topics {
		userIds = append(userIds, t.UserId)
	}
	usernameMap := s.loadUsernameMap(ctx, userIds)

	var list []forumout.TopicListItem
	for _, t := range topics {
		list = append(list, forumout.TopicListItem{
			Id:          t.Id,
			Subject:     t.Subject,
			UserId:      t.UserId,
			Username:    usernameMap[t.UserId],
			IsLocked:    t.IsLocked,
			IsSticky:    t.IsSticky,
			Views:       t.Views,
			ReplyCount:  t.ReplyCount,
			LastReplyAt: t.LastReplyAt.String(),
			CreatedAt:   t.CreatedAt.String(),
		})
	}

	return &forumout.TopicListOut{
		List:  list,
		Total: total,
		Node: forumout.NodeItem{
			Id:         node.Id,
			Slug:       node.Slug,
			NameI18N:   node.NameI18N,
			DescI18N:   node.DescI18N,
			TopicCount: node.TopicCount,
			ReplyCount: node.ReplyCount,
		},
	}, nil
}

func (s *sForumTopicUsecase) Detail(ctx context.Context, actor *model.Actor, in forumin.TopicDetailInp) (*forumout.TopicDetailOut, error) {
	topic, err := service.ForumTopicDomain().GetTopicById(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	nodePtr, err := service.ForumNodeDomain().GetNodeById(ctx, topic.NodeId)
	if err != nil || nodePtr == nil || nodePtr.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "forum.node.not_found"))
	}
	node := *nodePtr

	if err := service.ForumNodeDomain().CheckNodeReadPolicy(ctx, actor, &node); err != nil {
		return nil, err
	}

	if err := service.ForumTopicDomain().IncrementTopicViews(ctx, in.Id); err != nil {
		return nil, err
	}

	var user entity.IamUser
	userPtr, _ := service.IamUserDomain().GetUserById(ctx, topic.UserId)
	if userPtr != nil {
		user = *userPtr
	}

	isLiked := false
	isBookmarked := false
	if actor != nil && actor.Id > 0 {
		isLiked, _ = service.ForumTopicDomain().CheckTopicLiked(ctx, in.Id, actor.Id)
		isBookmarked, _ = service.ForumTopicDomain().CheckTopicBookmarked(ctx, in.Id, actor.Id)
	}

	return &forumout.TopicDetailOut{
		TopicListItem: forumout.TopicListItem{
			Id:          topic.Id,
			Subject:     topic.Subject,
			UserId:      topic.UserId,
			Username:    user.Username,
			IsLocked:    topic.IsLocked,
			IsSticky:    topic.IsSticky,
			Views:       topic.Views + 1,
			ReplyCount:  topic.ReplyCount,
			LastReplyAt: topic.LastReplyAt.String(),
			CreatedAt:   topic.CreatedAt.String(),
		},
		Content:      topic.Content,
		Appends:      topic.Appends,
		NodeId:       topic.NodeId,
		IsLiked:      isLiked,
		IsBookmarked: isBookmarked,
	}, nil
}

func (s *sForumTopicUsecase) Create(ctx context.Context, actor *model.Actor, in forumin.TopicCreateInp) (uint64, error) {
	if actor == nil {
		return 0, gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}

	nodePtr, err := service.ForumNodeDomain().GetNodeById(ctx, in.NodeId)
	if err != nil || nodePtr == nil || nodePtr.Id == 0 {
		return 0, gerror.New(gi18n.T(ctx, "forum.node.not_found"))
	}
	node := *nodePtr

	if err := service.ForumNodeDomain().CheckNodeCreatePolicy(ctx, actor, &node); err != nil {
		return 0, err
	}

	var topicId uint64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := service.ForumTopicDomain().InsertTopic(ctx, actor, in)
		if err != nil {
			return err
		}
		topicId = id

		// Update Node
		return service.ForumNodeDomain().IncrementNodeStats(ctx, in.NodeId, topicId)
	})

	if err != nil {
		return 0, gerror.Wrap(err, gi18n.T(ctx, "forum.topic.create_failed"))
	}
	// Assign update permission to creator
	_ = service.IamPermissionDomain().GrantUserPermission(ctx, actor.Id, fmt.Sprintf("update:forum/topic:%d", topicId), false)

	return topicId, nil
}

func (s *sForumTopicUsecase) Append(ctx context.Context, actor *model.Actor, in forumin.TopicAppendInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}

	topic, err := service.ForumTopicDomain().GetTopicById(ctx, in.Id)
	if err != nil {
		return err
	}

	if err := service.ForumTopicDomain().CheckTopicAppendPolicy(ctx, actor, topic); err != nil {
		return err
	}

	return service.ForumTopicDomain().AppendContent(ctx, topic, in.Content)
}

func (s *sForumTopicUsecase) ToggleTopicLike(ctx context.Context, actor *model.Actor, in forumin.TopicToggleLikeInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}
	topic, err := service.ForumTopicDomain().GetTopicById(ctx, in.Id)
	if err != nil {
		return err
	}
	nodePtr, err := service.ForumNodeDomain().GetNodeById(ctx, topic.NodeId)
	if err != nil || nodePtr == nil || nodePtr.Id == 0 {
		return gerror.New(gi18n.T(ctx, "forum.node.not_found"))
	}
	node := *nodePtr
	if err := service.ForumNodeDomain().CheckNodeReadPolicy(ctx, actor, &node); err != nil {
		return err
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, txErr := service.ForumTopicDomain().ToggleLike(ctx, actor, in.Id)
		return txErr
	})
	return err
}

func (s *sForumTopicUsecase) RewardTopic(ctx context.Context, actor *model.Actor, in forumin.TopicRewardInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}
	topic, err := service.ForumTopicDomain().GetTopicById(ctx, in.Id)
	if err != nil {
		return err
	}
	nodePtr, err := service.ForumNodeDomain().GetNodeById(ctx, topic.NodeId)
	if err != nil || nodePtr == nil || nodePtr.Id == 0 {
		return gerror.New(gi18n.T(ctx, "forum.node.not_found"))
	}
	node := *nodePtr
	if err := service.ForumNodeDomain().CheckNodeReadPolicy(ctx, actor, &node); err != nil {
		return err
	}
	if topic.UserId == actor.Id {
		return gerror.New(gi18n.T(ctx, "forum.topic.reward_self_not_allowed"))
	}
	err = service.EconomyBonusUsecase().TransferBonus(ctx, actor.Id, topic.UserId, in.Amount, "forum_topic", in.Id, "Reward topic", "Topic rewarded")
	return err
}

func (s *sForumTopicUsecase) ReportTopic(ctx context.Context, actor *model.Actor, in forumin.TopicReportInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}
	topic, err := service.ForumTopicDomain().GetTopicById(ctx, in.Id)
	if err != nil {
		return err
	}
	nodePtr, err := service.ForumNodeDomain().GetNodeById(ctx, topic.NodeId)
	if err != nil || nodePtr == nil || nodePtr.Id == 0 {
		return gerror.New(gi18n.T(ctx, "forum.node.not_found"))
	}
	node := *nodePtr
	if err := service.ForumNodeDomain().CheckNodeReadPolicy(ctx, actor, &node); err != nil {
		return err
	}
	return service.ModReportUsecase().Create(ctx, actor, modin.CreateReportInp{
		TargetType: "forum_topic",
		TargetId:   in.Id,
		Reason:     in.Reason,
	})
}

func (s *sForumTopicUsecase) BookmarkTopic(ctx context.Context, actor *model.Actor, in forumin.TopicBookmarkInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}
	topic, err := service.ForumTopicDomain().GetTopicById(ctx, in.Id)
	if err != nil {
		return err
	}
	nodePtr, err := service.ForumNodeDomain().GetNodeById(ctx, topic.NodeId)
	if err != nil || nodePtr == nil || nodePtr.Id == 0 {
		return gerror.New(gi18n.T(ctx, "forum.node.not_found"))
	}
	node := *nodePtr
	if err := service.ForumNodeDomain().CheckNodeReadPolicy(ctx, actor, &node); err != nil {
		return err
	}
	return service.ForumTopicDomain().Bookmark(ctx, actor, in.Id)
}

func (s *sForumTopicUsecase) UnbookmarkTopic(ctx context.Context, actor *model.Actor, in forumin.TopicUnbookmarkInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}
	return service.ForumTopicDomain().Unbookmark(ctx, actor, in.Id)
}

func (s *sForumTopicUsecase) ListBookmarkedTopics(ctx context.Context, actor *model.Actor, in forumin.TopicBookmarkListInp) (*forumout.TopicBookmarkListOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}
	topics, total, err := service.ForumTopicDomain().QueryBookmarkedTopics(ctx, actor, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	var list []forumout.TopicListItem
	userIds := make([]uint64, 0, len(topics))
	for _, t := range topics {
		userIds = append(userIds, t.UserId)
	}
	usernameMap := s.loadUsernameMap(ctx, userIds)

	for _, t := range topics {
		list = append(list, forumout.TopicListItem{
			Id:          t.Id,
			Subject:     t.Subject,
			UserId:      t.UserId,
			Username:    usernameMap[t.UserId],
			IsLocked:    t.IsLocked,
			IsSticky:    t.IsSticky,
			Views:       t.Views,
			ReplyCount:  t.ReplyCount,
			LastReplyAt: t.LastReplyAt.String(),
			CreatedAt:   t.CreatedAt.String(),
		})
	}
	return &forumout.TopicBookmarkListOut{
		List:  list,
		Total: total,
	}, nil
}

func (s *sForumTopicUsecase) loadUsernameMap(ctx context.Context, userIds []uint64) map[uint64]string {
	userMap := make(map[uint64]string)
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
	}
	if len(uniqueIds) == 0 {
		return userMap
	}

	users, err := service.IamUserDomain().GetUsersByIds(ctx, uniqueIds)
	if err != nil {
		return userMap
	}
	for _, user := range users {
		userMap[user.Id] = user.Username
	}
	return userMap
}
