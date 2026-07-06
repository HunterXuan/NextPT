package admin

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
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
	id, err := service.ForumNodeDomain().AdminCreateNode(ctx, in.CategoryId, in.Slug, in.NameI18N, in.DescI18N, in.SortOrder, in.MinRoleRead, in.MinRoleWrite, in.MinRoleCreate, moderators)
	if err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionCreate,
		TargetType: consts.SiteAuditTargetTypeForumNode,
		TargetId:   uint64(id),
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"categoryId":    in.CategoryId,
			"slug":          in.Slug,
			"minRoleRead":   in.MinRoleRead,
			"minRoleWrite":  in.MinRoleWrite,
			"minRoleCreate": in.MinRoleCreate,
		},
	})
	return nil
}

func (s *sAdminForumNodeUsecase) Update(ctx context.Context, actor *model.Actor, in adminin.ForumNodeUpdateInp) error {
	var moderators []byte
	if in.Moderators != nil {
		moderators, _ = gjson.Encode(*in.Moderators)
	}
	if err := service.ForumNodeDomain().AdminUpdateNode(ctx, in.Id, in.CategoryId, in.Slug, in.NameI18N, in.DescI18N, in.SortOrder, in.MinRoleRead, in.MinRoleWrite, in.MinRoleCreate, moderators); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeForumNode,
		TargetId:   uint64(in.Id),
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"categoryChanged":      in.CategoryId != nil,
			"slugChanged":          in.Slug != nil,
			"nameChanged":          in.NameI18N != nil,
			"descriptionChanged":   in.DescI18N != nil,
			"sortOrderChanged":     in.SortOrder != nil,
			"minRoleReadChanged":   in.MinRoleRead != nil,
			"minRoleWriteChanged":  in.MinRoleWrite != nil,
			"minRoleCreateChanged": in.MinRoleCreate != nil,
			"moderatorsChanged":    in.Moderators != nil,
		},
	})
	return nil
}

func (s *sAdminForumNodeUsecase) Delete(ctx context.Context, actor *model.Actor, in adminin.ForumNodeDeleteInp) error {
	count, err := service.ForumNodeDomain().AdminDeleteNode(ctx, in.Id)
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New(gi18n.T(ctx, "admin.forum.node_has_topics"))
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionDelete,
		TargetType: consts.SiteAuditTargetTypeForumNode,
		TargetId:   uint64(in.Id),
		Level:      consts.SiteAuditLevelCritical,
	})
	return nil
}

func (s *sAdminForumNodeUsecase) List(ctx context.Context, actor *model.Actor, in adminin.ForumNodeListInp) (*adminout.ForumNodeListOut, error) {
	list, err := service.ForumNodeDomain().AdminListNodes(ctx)
	if err != nil {
		return nil, err
	}
	return &adminout.ForumNodeListOut{Nodes: list}, nil
}
