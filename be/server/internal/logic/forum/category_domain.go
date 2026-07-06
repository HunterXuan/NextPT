package forum

import (
	"context"

	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sForumCategoryDomain struct{}

func init() {
	service.RegisterForumCategoryDomain(NewForumCategoryDomain())
}

func NewForumCategoryDomain() *sForumCategoryDomain {
	return &sForumCategoryDomain{}
}

func (s *sForumCategoryDomain) AdminCreateCategory(ctx context.Context, nameI18N, descI18N string, sortOrder, minRoleView int) (uint, error) {
	id, err := dao.ForumCategory.Ctx(ctx).Data(g.Map{
		dao.ForumCategory.Columns().NameI18N:    nameI18N,
		dao.ForumCategory.Columns().DescI18N:    descI18N,
		dao.ForumCategory.Columns().SortOrder:   sortOrder,
		dao.ForumCategory.Columns().MinRoleView: minRoleView,
	}).InsertAndGetId()
	return uint(id), err
}

func (s *sForumCategoryDomain) AdminUpdateCategory(ctx context.Context, id uint, nameI18N, descI18N *string, sortOrder, minRoleView *int) error {
	data := g.Map{}
	if nameI18N != nil {
		data[dao.ForumCategory.Columns().NameI18N] = *nameI18N
	}
	if descI18N != nil {
		data[dao.ForumCategory.Columns().DescI18N] = *descI18N
	}
	if sortOrder != nil {
		data[dao.ForumCategory.Columns().SortOrder] = *sortOrder
	}
	if minRoleView != nil {
		data[dao.ForumCategory.Columns().MinRoleView] = *minRoleView
	}
	if len(data) == 0 {
		return nil
	}
	_, err := dao.ForumCategory.Ctx(ctx).Where(dao.ForumCategory.Columns().Id, id).Data(data).Update()
	return err
}

func (s *sForumCategoryDomain) AdminDeleteCategory(ctx context.Context, id uint) (int, error) {
	count, err := dao.ForumNode.Ctx(ctx).Where(dao.ForumNode.Columns().CategoryId, id).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return count, nil
	}
	_, err = dao.ForumCategory.Ctx(ctx).Where(dao.ForumCategory.Columns().Id, id).Delete()
	return 0, err
}

func (s *sForumCategoryDomain) AdminListCategories(ctx context.Context) ([]entity.ForumCategory, error) {
	var list []entity.ForumCategory
	err := dao.ForumCategory.Ctx(ctx).OrderDesc(dao.ForumCategory.Columns().SortOrder).Scan(&list)
	return list, err
}
