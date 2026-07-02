package catalog

import (
	"context"

	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sCatalogCommentDomain struct{}

func init() {
	service.RegisterCatalogCommentDomain(NewCatalogCommentDomain())
}

func NewCatalogCommentDomain() *sCatalogCommentDomain {
	return &sCatalogCommentDomain{}
}

func (s *sCatalogCommentDomain) CreateComment(ctx context.Context, targetType string, targetId uint64, userId uint64, content string) (uint64, error) {
	result, err := dao.CatalogComment.Ctx(ctx).Insert(&entity.CatalogComment{
		TargetType: targetType,
		TargetId:   targetId,
		UserId:     userId,
		Content:    content,
	})
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return uint64(id), nil
}

func (s *sCatalogCommentDomain) GetCommentById(ctx context.Context, id uint64) (*entity.CatalogComment, error) {
	var comment entity.CatalogComment
	if err := dao.CatalogComment.Ctx(ctx).Where(dao.CatalogComment.Columns().Id, id).Scan(&comment); err != nil || comment.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "catalog.comment.not_found"))
	}
	return &comment, nil
}

func (s *sCatalogCommentDomain) QueryCommentsByTarget(ctx context.Context, targetType string, targetId uint64, page int, size int) ([]entity.CatalogComment, int, error) {
	m := dao.CatalogComment.Ctx(ctx).Where(dao.CatalogComment.Columns().TargetType, targetType).Where(dao.CatalogComment.Columns().TargetId, targetId)

	total, err := m.Count()
	if err != nil || total == 0 {
		return nil, 0, err
	}

	var comments []entity.CatalogComment
	err = m.OrderAsc(dao.CatalogComment.Columns().Id).Page(page, size).Scan(&comments)
	return comments, total, err
}

func (s *sCatalogCommentDomain) ToggleLike(ctx context.Context, userId uint64, commentId uint64) (bool, error) {
	likeColumns := dao.CatalogCommentLike.Columns()
	commentColumns := dao.CatalogComment.Columns()
	likeModel := dao.CatalogCommentLike.Ctx(ctx).
		Where(likeColumns.UserId, userId).
		Where(likeColumns.CommentId, commentId)

	count, err := likeModel.Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		_, err = likeModel.Delete()
		if err != nil {
			return false, err
		}
		_, err = dao.CatalogComment.Ctx(ctx).Where(commentColumns.Id, commentId).Decrement(commentColumns.LikeCount, 1)
		return false, err
	}

	_, err = dao.CatalogCommentLike.Ctx(ctx).Data(&entity.CatalogCommentLike{
		UserId:    userId,
		CommentId: commentId,
	}).Insert()
	if err != nil {
		return false, err
	}
	_, err = dao.CatalogComment.Ctx(ctx).Where(commentColumns.Id, commentId).Increment(commentColumns.LikeCount, 1)
	return true, err
}

func (s *sCatalogCommentDomain) GetCommentLikesByUser(ctx context.Context, userId uint64, commentIds []uint64) ([]entity.CatalogCommentLike, error) {
	var likes []entity.CatalogCommentLike
	err := dao.CatalogCommentLike.Ctx(ctx).Where(dao.CatalogCommentLike.Columns().UserId, userId).WhereIn(dao.CatalogCommentLike.Columns().CommentId, commentIds).Scan(&likes)
	return likes, err
}

func (s *sCatalogCommentDomain) IncrementCommentRewardStats(ctx context.Context, commentId uint64) error {
	_, err := dao.CatalogComment.Ctx(ctx).Where(dao.CatalogComment.Columns().Id, commentId).Increment(dao.CatalogComment.Columns().RewardCount, 1)
	return err
}

func (s *sCatalogCommentDomain) QueryCommentIdsByTarget(ctx context.Context, targetType string, targetId uint64) ([]uint64, error) {
	var commentIds []uint64
	err := dao.CatalogComment.Ctx(ctx).
		Where(dao.CatalogComment.Columns().TargetType, targetType).
		Where(dao.CatalogComment.Columns().TargetId, targetId).
		ScanList(&commentIds, "Id")
	return commentIds, err
}

func (s *sCatalogCommentDomain) DeleteCommentsByTarget(ctx context.Context, targetType string, targetId uint64) error {
	commentIds, err := s.QueryCommentIdsByTarget(ctx, targetType, targetId)
	if err != nil {
		return err
	}
	if len(commentIds) > 0 {
		if _, err := dao.CatalogCommentLike.Ctx(ctx).WhereIn(dao.CatalogCommentLike.Columns().CommentId, commentIds).Delete(); err != nil {
			return err
		}
	}
	_, err = dao.CatalogComment.Ctx(ctx).Where(dao.CatalogComment.Columns().TargetType, targetType).Where(dao.CatalogComment.Columns().TargetId, targetId).Delete()
	return err
}
