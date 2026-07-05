package admin

import (
	"context"
	"fmt"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type sAdminForumTopicUsecase struct{}

func NewAdminForumTopicUsecase() *sAdminForumTopicUsecase {
	return &sAdminForumTopicUsecase{}
}

func init() {
	service.RegisterAdminForumTopicUsecase(NewAdminForumTopicUsecase())
}

func (s *sAdminForumTopicUsecase) Lock(ctx context.Context, actor *model.Actor, in adminin.ForumTopicLockInp) error {
	return service.ForumTopicDomain().AdminSetTopicLock(ctx, in.Id, true)
}

func (s *sAdminForumTopicUsecase) Unlock(ctx context.Context, actor *model.Actor, in adminin.ForumTopicUnlockInp) error {
	return service.ForumTopicDomain().AdminSetTopicLock(ctx, in.Id, false)
}

func (s *sAdminForumTopicUsecase) Pin(ctx context.Context, actor *model.Actor, in adminin.ForumTopicPinInp) error {
	return service.ForumTopicDomain().AdminSetTopicSticky(ctx, in.Id, true)
}

func (s *sAdminForumTopicUsecase) Unpin(ctx context.Context, actor *model.Actor, in adminin.ForumTopicUnpinInp) error {
	return service.ForumTopicDomain().AdminSetTopicSticky(ctx, in.Id, false)
}

func (s *sAdminForumTopicUsecase) Move(ctx context.Context, actor *model.Actor, in adminin.ForumTopicMoveInp) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return service.ForumTopicDomain().AdminMoveTopic(ctx, in.Id, in.NodeId)
	})
}

func (s *sAdminForumTopicUsecase) Delete(ctx context.Context, actor *model.Actor, in adminin.ForumTopicDeleteInp) error {
	topic, err := service.ForumTopicDomain().GetTopicById(ctx, in.Id)
	if err != nil {
		return err
	}
	replyIds, err := service.ForumReplyDomain().QueryReplyIdsByTopic(ctx, in.Id)
	if err != nil {
		return err
	}

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := service.EconomyRewardDomain().DeleteRewardRecordsByTarget(ctx, consts.EconomyBonusTargetTypeForumTopic, topic.Id); err != nil {
			return err
		}
		if err := service.EconomyRewardDomain().DeleteRewardRecordsByTargets(ctx, consts.EconomyBonusTargetTypeForumReply, replyIds); err != nil {
			return err
		}
		if err := service.ModReportDomain().DeleteReportsByTarget(ctx, consts.ModReportTargetTypeForumTopic, topic.Id); err != nil {
			return err
		}
		if err := service.ModReportDomain().DeleteReportsByTargets(ctx, consts.ModReportTargetTypeForumReply, replyIds); err != nil {
			return err
		}
		if err := service.ForumTopicDomain().DeleteTopic(ctx, topic, replyIds); err != nil {
			return err
		}
		return service.IamPermissionDomain().RevokeUserPermission(ctx, topic.UserId, fmt.Sprintf("update:forum/topic:%d", topic.Id), false)
	})
}
