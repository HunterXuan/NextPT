package forum

import (
	"context"

	"server/internal/dao"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/forumin"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

type sForumReplyDomain struct{}

func init() {
	service.RegisterForumReplyDomain(NewForumReplyDomain())
}

func NewForumReplyDomain() *sForumReplyDomain {
	return &sForumReplyDomain{}
}

func (s *sForumReplyDomain) InsertReply(ctx context.Context, actor *model.Actor, in forumin.ReplyCreateInp) (uint64, error) {
	id, err := dao.ForumReply.Ctx(ctx).Data(entity.ForumReply{
		TopicId: in.Id,
		UserId:  actor.Id,
		Content: in.Content,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (s *sForumReplyDomain) QueryRepliesByTopic(ctx context.Context, topicId uint64, page, size int) ([]entity.ForumReply, int, error) {
	m := dao.ForumReply.Ctx(ctx).Where(dao.ForumReply.Columns().TopicId, topicId)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var replies []entity.ForumReply
	err = m.Page(page, size).OrderAsc(dao.ForumReply.Columns().Id).Scan(&replies)
	if err != nil {
		return nil, 0, err
	}
	return replies, total, nil
}

func (s *sForumReplyDomain) QueryReplyIdsByTopic(ctx context.Context, topicId uint64) ([]uint64, error) {
	var replyIds []uint64
	err := dao.ForumReply.Ctx(ctx).
		Where(dao.ForumReply.Columns().TopicId, topicId).
		ScanList(&replyIds, "Id")
	return replyIds, err
}

func (s *sForumReplyDomain) GetReplyById(ctx context.Context, replyId uint64) (*entity.ForumReply, error) {
	var reply entity.ForumReply
	err := dao.ForumReply.Ctx(ctx).Where(dao.ForumReply.Columns().Id, replyId).Scan(&reply)
	if err != nil {
		return nil, err
	} else if reply.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "forum.reply.not_found"))
	}

	return &reply, nil
}

func (s *sForumReplyDomain) GetRepliesByIds(ctx context.Context, ids []uint64) ([]entity.ForumReply, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var replies []entity.ForumReply
	err := dao.ForumReply.Ctx(ctx).WhereIn(dao.ForumReply.Columns().Id, ids).Scan(&replies)
	return replies, err
}

func (s *sForumReplyDomain) ToggleLike(ctx context.Context, actor *model.Actor, replyId uint64) (bool, error) {
	likeColumns := dao.ForumReplyLike.Columns()
	replyColumns := dao.ForumReply.Columns()
	likeModel := dao.ForumReplyLike.Ctx(ctx).Where(g.Map{
		likeColumns.UserId:  actor.Id,
		likeColumns.ReplyId: replyId,
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
		_, err = dao.ForumReply.Ctx(ctx).Where(replyColumns.Id, replyId).Decrement(replyColumns.LikeCount, 1)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	_, err = dao.ForumReplyLike.Ctx(ctx).Data(g.Map{
		likeColumns.UserId:    actor.Id,
		likeColumns.ReplyId:   replyId,
		likeColumns.CreatedAt: gtime.Now(),
	}).Insert()
	if err != nil {
		return false, err
	}
	_, err = dao.ForumReply.Ctx(ctx).Where(replyColumns.Id, replyId).Increment(replyColumns.LikeCount, 1)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *sForumReplyDomain) GetReplyLikesByUser(ctx context.Context, userId uint64, replyIds []uint64) ([]entity.ForumReplyLike, error) {
	if len(replyIds) == 0 {
		return nil, nil
	}
	var likes []entity.ForumReplyLike
	err := dao.ForumReplyLike.Ctx(ctx).Where(dao.ForumReplyLike.Columns().UserId, userId).WhereIn(dao.ForumReplyLike.Columns().ReplyId, replyIds).Scan(&likes)
	return likes, err
}

func (s *sForumReplyDomain) IncrementRewardStats(ctx context.Context, replyId uint64) error {
	_, err := dao.ForumReply.Ctx(ctx).Where(dao.ForumReply.Columns().Id, replyId).Increment(dao.ForumReply.Columns().RewardCount, 1)
	return err
}
