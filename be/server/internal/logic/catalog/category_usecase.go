package catalog

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"
	"server/internal/service"
)

type sCatalogCategoryUsecase struct{}

func init() {
	service.RegisterCatalogCategoryUsecase(NewCatalogCategoryUsecase())
}

func NewCatalogCategoryUsecase() *sCatalogCategoryUsecase {
	return &sCatalogCategoryUsecase{}
}

func (s *sCatalogCategoryUsecase) ListCategories(ctx context.Context, actor *model.Actor, in catalogin.CategoryListInp) (*catalogout.CategoryListOut, error) {
	entities, err := service.CatalogCategoryDomain().ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	var list []catalogout.CategoryItem
	for _, e := range entities {
		item := catalogout.CategoryItem{
			Id:   e.Id,
			Slug: e.Slug,
		}
		_ = e.NameI18N.Scan(&item.Name)
		list = append(list, item)
	}

	return &catalogout.CategoryListOut{List: list}, nil
}

func (s *sCatalogCategoryUsecase) ListTagGroups(ctx context.Context, actor *model.Actor, in catalogin.TagGroupListInp) (*catalogout.TagGroupListOut, error) {
	groups, tags, err := service.CatalogCategoryDomain().ListTagGroups(ctx)
	if err != nil {
		return nil, err
	}

	tagMap := make(map[uint][]catalogout.TagItem)
	for _, t := range tags {
		var name map[string]any
		_ = t.NameI18N.Scan(&name)
		tagMap[t.GroupId] = append(tagMap[t.GroupId], catalogout.TagItem{
			Id:   t.Id,
			Name: name,
		})
	}

	var list []catalogout.TagGroupItem
	for _, g := range groups {
		var name map[string]any
		var categories []uint
		_ = g.NameI18N.Scan(&name)
		_ = g.CategoryIds.Scan(&categories)

		list = append(list, catalogout.TagGroupItem{
			Id:         g.Id,
			Name:       name,
			Slug:       g.Slug,
			Categories: categories,
			Tags:       tagMap[g.Id],
		})
	}

	return &catalogout.TagGroupListOut{List: list}, nil
}
