package catalog

import (
	"context"

	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
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

func (s *sCatalogCategoryDomain) GetCategoryById(ctx context.Context, id uint) (*entity.CatalogCategory, error) {
	var category *entity.CatalogCategory
	err := dao.CatalogCategory.Ctx(ctx).
		Where(dao.CatalogCategory.Columns().Id, id).
		Where(dao.CatalogCategory.Columns().Enabled, 1).
		Scan(&category)
	return category, err
}

func (s *sCatalogCategoryDomain) AdminCreateCategory(ctx context.Context, nameI18N []byte, slug string, sortOrder int, enabled bool, uploadConfig []byte) error {
	_, err := dao.CatalogCategory.Ctx(ctx).Data(g.Map{
		dao.CatalogCategory.Columns().NameI18N:     nameI18N,
		dao.CatalogCategory.Columns().Slug:         slug,
		dao.CatalogCategory.Columns().SortOrder:    sortOrder,
		dao.CatalogCategory.Columns().Enabled:      enabled,
		dao.CatalogCategory.Columns().UploadConfig: uploadConfig,
	}).Insert()
	return err
}

func (s *sCatalogCategoryDomain) AdminUpdateCategory(ctx context.Context, id uint, nameI18N []byte, slug *string, sortOrder *int, enabled *bool, uploadConfig *[]byte) error {
	data := g.Map{}
	if nameI18N != nil {
		data[dao.CatalogCategory.Columns().NameI18N] = nameI18N
	}
	if slug != nil {
		data[dao.CatalogCategory.Columns().Slug] = *slug
	}
	if sortOrder != nil {
		data[dao.CatalogCategory.Columns().SortOrder] = *sortOrder
	}
	if enabled != nil {
		data[dao.CatalogCategory.Columns().Enabled] = *enabled
	}
	if uploadConfig != nil {
		data[dao.CatalogCategory.Columns().UploadConfig] = *uploadConfig
	}
	if len(data) == 0 {
		return nil
	}
	_, err := dao.CatalogCategory.Ctx(ctx).Where(dao.CatalogCategory.Columns().Id, id).Data(data).Update()
	return err
}

func (s *sCatalogCategoryDomain) AdminDeleteCategory(ctx context.Context, id uint) error {
	_, err := dao.CatalogCategory.Ctx(ctx).Where(dao.CatalogCategory.Columns().Id, id).Delete()
	return err
}

func (s *sCatalogCategoryDomain) AdminListCategories(ctx context.Context) ([]entity.CatalogCategory, error) {
	var categories []entity.CatalogCategory
	err := dao.CatalogCategory.Ctx(ctx).
		OrderAsc(dao.CatalogCategory.Columns().SortOrder).
		OrderAsc(dao.CatalogCategory.Columns().Id).
		Scan(&categories)
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
