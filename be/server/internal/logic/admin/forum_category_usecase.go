package admin

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/adminout"
	"server/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sAdminForumCategoryUsecase struct{}

func NewAdminForumCategoryUsecase() *sAdminForumCategoryUsecase {
	return &sAdminForumCategoryUsecase{}
}

func init() {
	service.RegisterAdminForumCategoryUsecase(NewAdminForumCategoryUsecase())
}

func (s *sAdminForumCategoryUsecase) Create(ctx context.Context, actor *model.Actor, in adminin.ForumCategoryCreateInp) error {
	id, err := service.ForumCategoryDomain().AdminCreateCategory(ctx, in.NameI18N, in.DescI18N, in.SortOrder, in.MinRoleView)
	if err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionCreate,
		TargetType: consts.SiteAuditTargetTypeForumCategory,
		TargetId:   uint64(id),
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"sortOrder":   in.SortOrder,
			"minRoleView": in.MinRoleView,
		},
	})
	return nil
}

func (s *sAdminForumCategoryUsecase) Update(ctx context.Context, actor *model.Actor, in adminin.ForumCategoryUpdateInp) error {
	if err := service.ForumCategoryDomain().AdminUpdateCategory(ctx, in.Id, in.NameI18N, in.DescI18N, in.SortOrder, in.MinRoleView); err != nil {
		return err
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionUpdate,
		TargetType: consts.SiteAuditTargetTypeForumCategory,
		TargetId:   uint64(in.Id),
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"nameChanged":        in.NameI18N != nil,
			"descriptionChanged": in.DescI18N != nil,
			"sortOrderChanged":   in.SortOrder != nil,
			"minRoleViewChanged": in.MinRoleView != nil,
		},
	})
	return nil
}

func (s *sAdminForumCategoryUsecase) Delete(ctx context.Context, actor *model.Actor, in adminin.ForumCategoryDeleteInp) error {
	category, err := service.ForumCategoryDomain().GetCategoryById(ctx, in.Id)
	if err != nil {
		return err
	}
	count, err := service.ForumCategoryDomain().AdminDeleteCategory(ctx, in.Id)
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New(gi18n.T(ctx, "admin.forum.category_has_nodes"))
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action:     consts.SiteAuditActionDelete,
		TargetType: consts.SiteAuditTargetTypeForumCategory,
		TargetId:   uint64(in.Id),
		Level:      consts.SiteAuditLevelCritical,
		Detail: map[string]any{
			"snapshot": s.forumCategoryAuditSnapshot(category),
		},
	})
	return nil
}

func (s *sAdminForumCategoryUsecase) List(ctx context.Context, actor *model.Actor, in adminin.ForumCategoryListInp) (*adminout.ForumCategoryListOut, error) {
	list, err := service.ForumCategoryDomain().AdminListCategories(ctx)
	if err != nil {
		return nil, err
	}
	return &adminout.ForumCategoryListOut{Categories: list}, nil
}

func (s *sAdminForumCategoryUsecase) forumCategoryAuditSnapshot(category *entity.ForumCategory) map[string]any {
	if category == nil {
		return nil
	}
	return map[string]any{
		"id":          category.Id,
		"nameI18N":    s.forumCategoryJSONMap(category.NameI18N),
		"sortOrder":   category.SortOrder,
		"minRoleView": category.MinRoleView,
	}
}

func (s *sAdminForumCategoryUsecase) forumCategoryJSONMap(value *gjson.Json) map[string]any {
	if value == nil {
		return nil
	}
	var data map[string]any
	if err := value.Scan(&data); err != nil {
		return nil
	}
	return data
}
