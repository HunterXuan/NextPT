package forum

import (
	"context"
	"fmt"
	"strings"
	"time"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/forumin"
	"server/internal/model/in/modin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/forumout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
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

	return &forumout.TopicListOut{
		List:  s.formatTopicListItems(ctx, topics),
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

func (s *sForumTopicUsecase) ListHot(ctx context.Context, actor *model.Actor, size int) (*forumout.TopicHotListOut, error) {
	if size <= 0 {
		size = 5
	}
	if size > 10 {
		size = 10
	}

	items, err := s.getHotTopicItemsCache(ctx, actor)
	if err != nil {
		return nil, err
	}
	if len(items) > size {
		items = items[:size]
	}

	return &forumout.TopicHotListOut{
		List:  items,
		Total: len(items),
	}, nil
}

func (s *sForumTopicUsecase) getHotTopicItemsCache(ctx context.Context, actor *model.Actor) ([]forumout.TopicHotItem, error) {
	roleLevel := s.hotTopicRoleLevel(actor)
	cacheKey := fmt.Sprintf("%s:role:%d", service.SysCache().KeyForumHotTopics(ctx), roleLevel)
	v, err := gcache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		return s.buildHotTopicItems(ctx, roleLevel)
	}, 5*time.Minute)
	if err != nil || v.IsNil() {
		return s.buildHotTopicItems(ctx, roleLevel)
	}

	if items, ok := v.Val().([]forumout.TopicHotItem); ok {
		return items, nil
	}
	var items []forumout.TopicHotItem
	if err := v.Scan(&items); err != nil {
		return s.buildHotTopicItems(ctx, roleLevel)
	}
	return items, nil
}

func (s *sForumTopicUsecase) buildHotTopicItems(ctx context.Context, roleLevel int) ([]forumout.TopicHotItem, error) {
	topics, err := service.ForumTopicDomain().QueryHotTopics(ctx, 100)
	if err != nil {
		return nil, err
	}

	nodes, err := service.ForumNodeDomain().GetNodes(ctx)
	if err != nil {
		return nil, err
	}
	nodeMap := make(map[uint]entity.ForumNode, len(nodes))
	for _, node := range nodes {
		nodeMap[node.Id] = node
	}

	selectedTopics := make([]entity.ForumTopic, 0, 10)
	selectedNodes := make([]entity.ForumNode, 0, 10)
	for _, topic := range topics {
		node, ok := nodeMap[topic.NodeId]
		if !ok || roleLevel < node.MinRoleRead {
			continue
		}
		selectedTopics = append(selectedTopics, topic)
		selectedNodes = append(selectedNodes, node)
		if len(selectedTopics) >= 10 {
			break
		}
	}

	items := s.formatTopicListItems(ctx, selectedTopics)
	list := make([]forumout.TopicHotItem, 0, len(items))
	for i, item := range items {
		list = append(list, forumout.TopicHotItem{
			TopicListItem: item,
			Node:          s.formatNodeItem(selectedNodes[i]),
		})
	}
	return list, nil
}

func (s *sForumTopicUsecase) hotTopicRoleLevel(actor *model.Actor) int {
	if actor == nil {
		return 0
	}
	return actor.RoleLevel
}

func (s *sForumTopicUsecase) formatNodeItem(node entity.ForumNode) forumout.NodeItem {
	return forumout.NodeItem{
		Id:         node.Id,
		Slug:       node.Slug,
		NameI18N:   node.NameI18N,
		DescI18N:   node.DescI18N,
		TopicCount: node.TopicCount,
		ReplyCount: node.ReplyCount,
	}
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

	isLiked := false
	isBookmarked := false
	if actor != nil && actor.Id > 0 {
		isLiked, _ = service.ForumTopicDomain().CheckTopicLiked(ctx, in.Id, actor.Id)
		isBookmarked, _ = service.ForumTopicDomain().CheckTopicBookmarked(ctx, in.Id, actor.Id)
	}

	listItems := s.formatTopicListItems(ctx, []entity.ForumTopic{*topic})
	var listItem forumout.TopicListItem
	if len(listItems) > 0 {
		listItem = listItems[0]
	}
	listItem.Views = topic.Views + 1

	return &forumout.TopicDetailOut{
		TopicListItem: listItem,
		Content:       topic.Content,
		Appends:       topic.Appends,
		NodeId:        topic.NodeId,
		IsLiked:       isLiked,
		IsBookmarked:  isBookmarked,
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

func (s *sForumTopicUsecase) Update(ctx context.Context, actor *model.Actor, in forumin.TopicUpdateInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}

	topic, err := service.ForumTopicDomain().GetTopicById(ctx, in.Id)
	if err != nil {
		return err
	}
	if err := service.ForumTopicDomain().CheckTopicEditPolicy(ctx, actor, topic); err != nil {
		return err
	}

	nodePtr, err := service.ForumNodeDomain().GetNodeById(ctx, in.NodeId)
	if err != nil || nodePtr == nil || nodePtr.Id == 0 {
		return gerror.New(gi18n.T(ctx, "forum.node.not_found"))
	}
	node := *nodePtr
	if err := service.ForumNodeDomain().CheckNodeCreatePolicy(ctx, actor, &node); err != nil {
		return err
	}

	subject := strings.TrimSpace(in.Subject)
	content := strings.TrimSpace(in.Content)
	oldNodeId := topic.NodeId
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := service.ForumTopicDomain().UpdateTopic(ctx, in.Id, in.NodeId, subject, content); err != nil {
			return err
		}
		if oldNodeId == in.NodeId {
			return nil
		}
		replyCount := int(topic.ReplyCount)
		if err := service.ForumNodeDomain().UpdateStats(ctx, oldNodeId, -1, -replyCount); err != nil {
			return err
		}
		return service.ForumNodeDomain().UpdateStats(ctx, in.NodeId, 1, replyCount)
	})
	if err != nil {
		return gerror.Wrap(err, gi18n.T(ctx, "forum.topic.update_failed"))
	}
	return nil
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
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := service.EconomyBonusUsecase().TransferBonus(ctx, actor.Id, topic.UserId, in.Amount, consts.EconomyBonusTargetTypeForumTopic, in.Id); err != nil {
			return err
		}
		return service.EconomyRewardDomain().InsertRewardRecord(ctx, entity.EconomyRewardRecord{
			TargetType: consts.EconomyBonusTargetTypeForumTopic,
			TargetId:   in.Id,
			FromUserId: actor.Id,
			ToUserId:   topic.UserId,
			Amount:     in.Amount,
		})
	})
	if err != nil {
		return err
	}

	service.SiteMessageUsecase().Notify(ctx, sitein.MessageNotifyInp{
		ActorId:     actor.Id,
		ReceiverId:  topic.UserId,
		TitleKey:    "site.message.reward.forum_topic.title",
		ContentKey:  "site.message.reward.content",
		ContentArgs: []any{in.Amount},
		TargetType:  consts.SiteMessageTargetTypeForumTopic,
		TargetId:    topic.Id,
	})
	return nil
}

func (s *sForumTopicUsecase) RewardList(ctx context.Context, actor *model.Actor, in forumin.TopicRewardListInp) (*forumout.TopicRewardListOut, error) {
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

	summaries, total, err := service.EconomyRewardDomain().QueryRewardSummaries(ctx, consts.EconomyBonusTargetTypeForumTopic, in.Id, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	userIds := make([]uint64, 0, len(summaries))
	for _, summary := range summaries {
		userIds = append(userIds, summary.UserId)
	}
	userMap := s.loadUserSummaryMap(ctx, userIds)

	list := make([]forumout.TopicRewardItem, 0, len(summaries))
	for _, summary := range summaries {
		list = append(list, forumout.TopicRewardItem{
			User:         userMap[summary.UserId],
			Amount:       summary.Amount,
			RewardCount:  summary.RewardCount,
			LastRewardAt: s.formatTime(summary.LastRewardAt),
		})
	}

	return &forumout.TopicRewardListOut{
		List:  list,
		Total: total,
	}, nil
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
		TargetType: consts.ModReportTargetTypeForumTopic,
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
	return &forumout.TopicBookmarkListOut{
		List:  s.formatTopicListItems(ctx, topics),
		Total: total,
	}, nil
}

func (s *sForumTopicUsecase) formatTopicListItems(ctx context.Context, topics []entity.ForumTopic) []forumout.TopicListItem {
	userIds := make([]uint64, 0, len(topics)*2)
	for _, topic := range topics {
		userIds = append(userIds, topic.UserId)
		if topic.ReplyCount > 0 && topic.LastReplyBy > 0 {
			userIds = append(userIds, topic.LastReplyBy)
		}
	}
	userMap := s.loadUserSummaryMap(ctx, userIds)

	list := make([]forumout.TopicListItem, 0, len(topics))
	for _, topic := range topics {
		lastReplyAt := ""
		lastReplyUser := model.IamUserSummary{}
		if topic.ReplyCount > 0 {
			lastReplyAt = s.formatTime(topic.LastReplyAt)
			lastReplyUser = userMap[topic.LastReplyBy]
		}

		list = append(list, forumout.TopicListItem{
			Id:            topic.Id,
			Subject:       topic.Subject,
			Author:        userMap[topic.UserId],
			IsLocked:      topic.IsLocked,
			IsSticky:      topic.IsSticky,
			Views:         topic.Views,
			ReplyCount:    topic.ReplyCount,
			LikeCount:     topic.LikeCount,
			LastReplyAt:   lastReplyAt,
			LastReplyUser: lastReplyUser,
			CreatedAt:     s.formatTime(topic.CreatedAt),
		})
	}
	return list
}

func (s *sForumTopicUsecase) formatTime(value *gtime.Time) string {
	if value == nil {
		return ""
	}
	return value.String()
}

func (s *sForumTopicUsecase) loadUserSummaryMap(ctx context.Context, userIds []uint64) map[uint64]model.IamUserSummary {
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
