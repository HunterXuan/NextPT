package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"
	"server/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sAdminForumNodeUsecase struct{}

func NewAdminForumNodeUsecase() *sAdminForumNodeUsecase {
	return &sAdminForumNodeUsecase{}
}

func init() {
	service.RegisterAdminForumNodeUsecase(NewAdminForumNodeUsecase())
}

func (s *sAdminForumNodeUsecase) Create(ctx context.Context, actor *model.Actor, in adminin.ForumNodeCreateInp) error {
	moderators, _ := gjson.Encode(in.Moderators)
	return service.ForumNodeDomain().AdminCreateNode(ctx, in.CategoryId, in.Slug, in.NameI18N, in.DescI18N, in.SortOrder, in.MinRoleRead, in.MinRoleWrite, in.MinRoleCreate, moderators)
}

func (s *sAdminForumNodeUsecase) Update(ctx context.Context, actor *model.Actor, in adminin.ForumNodeUpdateInp) error {
	var moderators []byte
	if in.Moderators != nil {
		moderators, _ = gjson.Encode(*in.Moderators)
	}
	return service.ForumNodeDomain().AdminUpdateNode(ctx, in.Id, in.CategoryId, in.Slug, in.NameI18N, in.DescI18N, in.SortOrder, in.MinRoleRead, in.MinRoleWrite, in.MinRoleCreate, moderators)
}

func (s *sAdminForumNodeUsecase) Delete(ctx context.Context, actor *model.Actor, in adminin.ForumNodeDeleteInp) error {
	count, err := service.ForumNodeDomain().AdminDeleteNode(ctx, in.Id)
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New(gi18n.T(ctx, "admin.forum.node_has_topics"))
	}
	return nil
}

func (s *sAdminForumNodeUsecase) List(ctx context.Context, actor *model.Actor, in adminin.ForumNodeListInp) (*adminout.ForumNodeListOut, error) {
	list, err := service.ForumNodeDomain().AdminListNodes(ctx)
	if err != nil {
		return nil, err
	}
	return &adminout.ForumNodeListOut{Nodes: list}, nil
}
