package catalog

import (
	"context"

	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"
)

type sCatalogCategoryDomain struct{}

func init() {
	service.RegisterCatalogCategoryDomain(NewCatalogCategoryDomain())
}

func NewCatalogCategoryDomain() *sCatalogCategoryDomain {
	return &sCatalogCategoryDomain{}
}

func (s *sCatalogCategoryDomain) ListCategories(ctx context.Context) ([]entity.CatalogCategory, error) {
	var categories []entity.CatalogCategory
	err := dao.CatalogCategory.Ctx(ctx).Where(dao.CatalogCategory.Columns().Enabled, 1).OrderAsc(dao.CatalogCategory.Columns().SortOrder).Scan(&categories)
	return categories, err
}

func (s *sCatalogCategoryDomain) ListTagGroups(ctx context.Context) ([]entity.CatalogTagGroup, []entity.CatalogTag, error) {
	var groups []entity.CatalogTagGroup
	err := dao.CatalogTagGroup.Ctx(ctx).OrderAsc(dao.CatalogTagGroup.Columns().SortOrder).Scan(&groups)
	if err != nil {
		return nil, nil, err
	}
	var tags []entity.CatalogTag
	err = dao.CatalogTag.Ctx(ctx).OrderAsc(dao.CatalogTag.Columns().SortOrder).Scan(&tags)
	return groups, tags, err
}
