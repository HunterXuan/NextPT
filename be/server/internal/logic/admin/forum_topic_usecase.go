package admin

import (
	"context"
	"fmt"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
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
	if err := service.ForumTopicDomain().AdminSetTopicLock(ctx, in.Id, true); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeForumTopic,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"operation": consts.SiteAuditOperationLock,
		},
	})
	return nil
}

func (s *sAdminForumTopicUsecase) Unlock(ctx context.Context, actor *model.Actor, in adminin.ForumTopicUnlockInp) error {
	if err := service.ForumTopicDomain().AdminSetTopicLock(ctx, in.Id, false); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeForumTopic,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"operation": consts.SiteAuditOperationUnlock,
		},
	})
	return nil
}

func (s *sAdminForumTopicUsecase) Pin(ctx context.Context, actor *model.Actor, in adminin.ForumTopicPinInp) error {
	if err := service.ForumTopicDomain().AdminSetTopicSticky(ctx, in.Id, true); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeForumTopic,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"operation": consts.SiteAuditOperationPin,
		},
	})
	return nil
}

func (s *sAdminForumTopicUsecase) Unpin(ctx context.Context, actor *model.Actor, in adminin.ForumTopicUnpinInp) error {
	if err := service.ForumTopicDomain().AdminSetTopicSticky(ctx, in.Id, false); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeForumTopic,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"operation": consts.SiteAuditOperationUnpin,
		},
	})
	return nil
}

func (s *sAdminForumTopicUsecase) Move(ctx context.Context, actor *model.Actor, in adminin.ForumTopicMoveInp) error {
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return service.ForumTopicDomain().AdminMoveTopic(ctx, in.Id, in.NodeId)
	})
	if err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeForumTopic,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"operation": consts.SiteAuditOperationMove,
			"nodeId":    in.NodeId,
		},
	})
	return nil
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

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
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
	if err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionDelete,
		TargetType: consts.SiteAuditTargetTypeForumTopic,
		TargetId:   in.Id,
		Level:      consts.SiteAuditLevelImportant,
		Detail: map[string]any{
			"replyCount": len(replyIds),
			"snapshot": map[string]any{
				"id":         topic.Id,
				"nodeId":     topic.NodeId,
				"subject":    topic.Subject,
				"userId":     topic.UserId,
				"replyCount": topic.ReplyCount,
			},
		},
	})
	return nil
}
