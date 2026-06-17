package admin

import (
	"context"

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
