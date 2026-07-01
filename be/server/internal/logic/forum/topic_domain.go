package forum

import (
	"context"
	"fmt"
	"time"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/forumin"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

type sForumTopicDomain struct{}

func init() {
	service.RegisterForumTopicDomain(NewForumTopicDomain())
}

func NewForumTopicDomain() *sForumTopicDomain {
	return &sForumTopicDomain{}
}

func (s *sForumTopicDomain) GetTopicById(ctx context.Context, topicId uint64) (*entity.ForumTopic, error) {
	var topic entity.ForumTopic
	if err := dao.ForumTopic.Ctx(ctx).Where(dao.ForumTopic.Columns().Id, topicId).Scan(&topic); err != nil || topic.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "forum.topic.not_found"))
	}
	return &topic, nil
}

func (s *sForumTopicDomain) CheckTopicWritePolicy(ctx context.Context, actor *model.Actor, topic *entity.ForumTopic) error {
	if topic.IsLocked {
		return gerror.New(gi18n.T(ctx, "forum.topic.is_locked"))
	}
	return nil
}

func (s *sForumTopicDomain) CheckTopicEditPolicy(ctx context.Context, actor *model.Actor, topic *entity.ForumTopic) error {
	if actor == nil || actor.Id == 0 {
		return gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}
	if topic.UserId != actor.Id {
		return gerror.New(gi18n.T(ctx, "forum.topic.edit_author_only"))
	}
	if topic.IsLocked {
		return gerror.New(gi18n.T(ctx, "forum.topic.is_locked"))
	}
	if topic.CreatedAt == nil || topic.CreatedAt.Add(time.Duration(consts.ForumTopicEditWindowSeconds)*time.Second).Before(gtime.Now()) {
		return gerror.New(gi18n.T(ctx, "forum.topic.edit_window_expired"))
	}
	return nil
}

func (s *sForumTopicDomain) CheckTopicAppendPolicy(ctx context.Context, actor *model.Actor, topic *entity.ForumTopic) error {
	if topic.UserId != actor.Id {
		return gerror.New(gi18n.T(ctx, "forum.topic.append_author_only"))
	}
	if topic.IsLocked {
		return gerror.New(gi18n.T(ctx, "forum.topic.is_locked"))
	}
	return nil
}

func (s *sForumTopicDomain) InsertTopic(ctx context.Context, actor *model.Actor, in forumin.TopicCreateInp) (uint64, error) {
	id, err := dao.ForumTopic.Ctx(ctx).Data(entity.ForumTopic{
		NodeId:      in.NodeId,
		UserId:      actor.Id,
		Subject:     in.Subject,
		Content:     in.Content,
		Appends:     gjson.New([]any{}),
		LastReplyAt: gtime.Now(),
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (s *sForumTopicDomain) AppendContent(ctx context.Context, topic *entity.ForumTopic, content string) error {
	var appends []any
	if topic.Appends != nil {
		if err := topic.Appends.Scan(&appends); err != nil {
			appends = []any{}
		}
	}

	if len(appends) >= 3 {
		return gerror.New(gi18n.T(ctx, "forum.topic.append_limit_exceeded"))
	}

	appends = append(appends, map[string]any{
		"content":    content,
		"created_at": gtime.Now().String(),
	})

	_, err := dao.ForumTopic.Ctx(ctx).Where(dao.ForumTopic.Columns().Id, topic.Id).Data(g.Map{
		dao.ForumTopic.Columns().Appends: gjson.New(appends),
	}).Update()

	return err
}

func (s *sForumTopicDomain) UpdateTopic(ctx context.Context, id uint64, nodeId uint, subject string, content string) error {
	_, err := dao.ForumTopic.Ctx(ctx).Where(dao.ForumTopic.Columns().Id, id).Data(g.Map{
		dao.ForumTopic.Columns().NodeId:  nodeId,
		dao.ForumTopic.Columns().Subject: subject,
		dao.ForumTopic.Columns().Content: content,
	}).Update()
	return err
}

func (s *sForumTopicDomain) IncrementTopicViews(ctx context.Context, topicId uint64) error {
	_, err := dao.ForumTopic.Ctx(ctx).Where(dao.ForumTopic.Columns().Id, topicId).Increment(dao.ForumTopic.Columns().Views, 1)
	return err
}

func (s *sForumTopicDomain) UpdateTopicReplyStats(ctx context.Context, topicId uint64, replyId uint64, lastReplyBy uint64) error {
	_, err := dao.ForumTopic.Ctx(ctx).Where(dao.ForumTopic.Columns().Id, topicId).Update(g.Map{
		dao.ForumTopic.Columns().ReplyCount:  gdb.Raw("reply_count + 1"),
		dao.ForumTopic.Columns().LastReplyId: replyId,
		dao.ForumTopic.Columns().LastReplyAt: gtime.Now(),
		dao.ForumTopic.Columns().LastReplyBy: lastReplyBy,
	})
	return err
}

func (s *sForumTopicDomain) ToggleLike(ctx context.Context, actor *model.Actor, topicId uint64) (bool, error) {
	return s.toggleLikeInternal(ctx, actor.Id, topicId)
}

func (s *sForumTopicDomain) toggleLikeInternal(ctx context.Context, userId uint64, topicId uint64) (bool, error) {
	likeColumns := dao.ForumTopicLike.Columns()
	topicColumns := dao.ForumTopic.Columns()
	likeModel := dao.ForumTopicLike.Ctx(ctx).Where(g.Map{
		likeColumns.UserId:  userId,
		likeColumns.TopicId: topicId,
	})
	count, err := likeModel.Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		_, err := likeModel.Delete()
		if err != nil {
			return false, err
		}
		_, err = dao.ForumTopic.Ctx(ctx).Where(topicColumns.Id, topicId).Decrement(topicColumns.LikeCount, 1)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	_, err = dao.ForumTopicLike.Ctx(ctx).Data(g.Map{
		likeColumns.UserId:    userId,
		likeColumns.TopicId:   topicId,
		likeColumns.CreatedAt: gtime.Now(),
	}).Insert()
	if err != nil {
		return false, err
	}
	_, err = dao.ForumTopic.Ctx(ctx).Where(topicColumns.Id, topicId).Increment(topicColumns.LikeCount, 1)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *sForumTopicDomain) Bookmark(ctx context.Context, actor *model.Actor, topicId uint64) error {
	count, err := dao.ForumTopicBookmark.Ctx(ctx).Where(g.Map{
		dao.ForumTopicBookmark.Columns().UserId:  actor.Id,
		dao.ForumTopicBookmark.Columns().TopicId: topicId,
	}).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Already bookmarked
	}

	_, err = dao.ForumTopicBookmark.Ctx(ctx).Data(entity.ForumTopicBookmark{
		UserId:  actor.Id,
		TopicId: topicId,
	}).Insert()
	return err
}

func (s *sForumTopicDomain) Unbookmark(ctx context.Context, actor *model.Actor, topicId uint64) error {
	_, err := dao.ForumTopicBookmark.Ctx(ctx).Where(g.Map{
		dao.ForumTopicBookmark.Columns().UserId:  actor.Id,
		dao.ForumTopicBookmark.Columns().TopicId: topicId,
	}).Delete()
	return err
}

func (s *sForumTopicDomain) QueryBookmarkedTopics(ctx context.Context, actor *model.Actor, page, size int) ([]entity.ForumTopic, int, error) {
	bm := dao.ForumTopicBookmark.Ctx(ctx).Where(dao.ForumTopicBookmark.Columns().UserId, actor.Id)
	total, err := bm.Count()
	if err != nil {
		return nil, 0, err
	}

	var bookmarks []entity.ForumTopicBookmark
	err = bm.Page(page, size).OrderDesc(dao.ForumTopicBookmark.Columns().CreatedAt).Scan(&bookmarks)
	if err != nil || len(bookmarks) == 0 {
		return nil, total, err
	}

	var topicIds []uint64
	for _, b := range bookmarks {
		topicIds = append(topicIds, b.TopicId)
	}

	var unorderedTopics []entity.ForumTopic
	err = dao.ForumTopic.Ctx(ctx).WhereIn(dao.ForumTopic.Columns().Id, topicIds).Scan(&unorderedTopics)
	if err != nil {
		return nil, 0, err
	}

	topicMap := make(map[uint64]entity.ForumTopic)
	for _, t := range unorderedTopics {
		topicMap[t.Id] = t
	}

	var topics []entity.ForumTopic
	for _, b := range bookmarks {
		if t, ok := topicMap[b.TopicId]; ok {
			topics = append(topics, t)
		}
	}

	return topics, total, nil
}

func (s *sForumTopicDomain) AdminSetTopicLock(ctx context.Context, id uint64, isLocked bool) error {
	_, err := dao.ForumTopic.Ctx(ctx).Where(dao.ForumTopic.Columns().Id, id).Data(g.Map{dao.ForumTopic.Columns().IsLocked: isLocked}).Update()
	return err
}

func (s *sForumTopicDomain) AdminSetTopicSticky(ctx context.Context, id uint64, isSticky bool) error {
	_, err := dao.ForumTopic.Ctx(ctx).Where(dao.ForumTopic.Columns().Id, id).Data(g.Map{dao.ForumTopic.Columns().IsSticky: isSticky}).Update()
	return err
}

func (s *sForumTopicDomain) AdminMoveTopic(ctx context.Context, id uint64, newNodeId uint) error {
	var topic entity.ForumTopic
	err := dao.ForumTopic.Ctx(ctx).Where(dao.ForumTopic.Columns().Id, id).Scan(&topic)
	if err != nil || topic.Id == 0 {
		return gerror.New("Topic not found")
	}
	oldNodeId := topic.NodeId
	if oldNodeId == newNodeId {
		return nil
	}
	_, err = dao.ForumTopic.Ctx(ctx).Where(dao.ForumTopic.Columns().Id, id).Data(g.Map{dao.ForumTopic.Columns().NodeId: newNodeId}).Update()
	if err != nil {
		return err
	}
	nodeColumns := dao.ForumNode.Columns()
	_, err = dao.ForumNode.Ctx(ctx).Where(nodeColumns.Id, oldNodeId).Data(g.Map{
		nodeColumns.TopicCount: gdb.Raw("topic_count - 1"),
		nodeColumns.ReplyCount: gdb.Raw(fmt.Sprintf("reply_count - %d", topic.ReplyCount)),
	}).Update()
	if err != nil {
		return err
	}
	_, err = dao.ForumNode.Ctx(ctx).Where(nodeColumns.Id, newNodeId).Data(g.Map{
		nodeColumns.TopicCount: gdb.Raw("topic_count + 1"),
		nodeColumns.ReplyCount: gdb.Raw(fmt.Sprintf("reply_count + %d", topic.ReplyCount)),
	}).Update()
	return err
}

func (s *sForumTopicDomain) QueryTopicsByNode(ctx context.Context, nodeId uint, page, size int) ([]entity.ForumTopic, int, error) {
	m := dao.ForumTopic.Ctx(ctx).Where(dao.ForumTopic.Columns().NodeId, nodeId)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var topics []entity.ForumTopic
	err = m.Page(page, size).OrderDesc(dao.ForumTopic.Columns().IsSticky).OrderDesc(dao.ForumTopic.Columns().LastReplyAt).Scan(&topics)
	return topics, total, err
}

func (s *sForumTopicDomain) CheckTopicLiked(ctx context.Context, topicId, userId uint64) (bool, error) {
	count, err := dao.ForumTopicLike.Ctx(ctx).Where(dao.ForumTopicLike.Columns().TopicId, topicId).Where(dao.ForumTopicLike.Columns().UserId, userId).Count()
	return count > 0, err
}

func (s *sForumTopicDomain) CheckTopicBookmarked(ctx context.Context, topicId, userId uint64) (bool, error) {
	count, err := dao.ForumTopicBookmark.Ctx(ctx).Where(dao.ForumTopicBookmark.Columns().TopicId, topicId).Where(dao.ForumTopicBookmark.Columns().UserId, userId).Count()
	return count > 0, err
}
