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

func (s *sForumReplyDomain) ToggleLike(ctx context.Context, actor *model.Actor, replyId uint64) (bool, error) {
	var likeRecord entity.ForumReplyLike
	err := dao.ForumReplyLike.Ctx(ctx).Where(g.Map{
		dao.ForumReplyLike.Columns().UserId:  actor.Id,
		dao.ForumReplyLike.Columns().ReplyId: replyId,
	}).Scan(&likeRecord)

	if err != nil {
		return false, err
	}

	isLiked := false
	if likeRecord.Id > 0 {
		// Already liked, so unlike
		_, err := dao.ForumReplyLike.Ctx(ctx).Where(dao.ForumReplyLike.Columns().Id, likeRecord.Id).Delete()
		if err != nil {
			return false, err
		}
		_, err = dao.ForumReply.Ctx(ctx).Where(dao.ForumReply.Columns().Id, replyId).Decrement(dao.ForumReply.Columns().LikeCount, 1)
		if err != nil {
			return false, err
		}
		isLiked = false
	} else {
		// Not liked, so like
		_, err := dao.ForumReplyLike.Ctx(ctx).Data(entity.ForumReplyLike{
			UserId:  actor.Id,
			ReplyId: replyId,
		}).Insert()
		if err != nil {
			return false, err
		}
		_, err = dao.ForumReply.Ctx(ctx).Where(dao.ForumReply.Columns().Id, replyId).Increment(dao.ForumReply.Columns().LikeCount, 1)
		if err != nil {
			return false, err
		}
		isLiked = true
	}
	return isLiked, err
}

func (s *sForumReplyDomain) GetReplyLikesByUser(ctx context.Context, userId uint64, replyIds []uint64) ([]entity.ForumReplyLike, error) {
	if len(replyIds) == 0 {
		return nil, nil
	}
	var likes []entity.ForumReplyLike
	err := dao.ForumReplyLike.Ctx(ctx).Where(dao.ForumReplyLike.Columns().UserId, userId).WhereIn(dao.ForumReplyLike.Columns().ReplyId, replyIds).Scan(&likes)
	return likes, err
}
