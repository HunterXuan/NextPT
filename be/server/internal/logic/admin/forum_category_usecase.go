package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"
	"server/internal/service"

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
	return service.ForumCategoryDomain().AdminCreateCategory(ctx, in.NameI18N, in.DescI18N, in.SortOrder, in.MinRoleView)
}

func (s *sAdminForumCategoryUsecase) Update(ctx context.Context, actor *model.Actor, in adminin.ForumCategoryUpdateInp) error {
	return service.ForumCategoryDomain().AdminUpdateCategory(ctx, in.Id, in.NameI18N, in.DescI18N, in.SortOrder, in.MinRoleView)
}

func (s *sAdminForumCategoryUsecase) Delete(ctx context.Context, actor *model.Actor, in adminin.ForumCategoryDeleteInp) error {
	count, err := service.ForumCategoryDomain().AdminDeleteCategory(ctx, in.Id)
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New(gi18n.T(ctx, "admin.forum.category_has_nodes"))
	}
	return nil
}

func (s *sAdminForumCategoryUsecase) List(ctx context.Context, actor *model.Actor, in adminin.ForumCategoryListInp) (*adminout.ForumCategoryListOut, error) {
	list, err := service.ForumCategoryDomain().AdminListCategories(ctx)
	if err != nil {
		return nil, err
	}
	return &adminout.ForumCategoryListOut{Categories: list}, nil
}
