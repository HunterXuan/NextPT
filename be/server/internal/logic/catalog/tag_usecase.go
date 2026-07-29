package catalog

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"
	"server/internal/service"
)

type sCatalogTagUsecase struct{}

func init() {
	service.RegisterCatalogTagUsecase(NewCatalogTagUsecase())
}

func NewCatalogTagUsecase() *sCatalogTagUsecase {
	return &sCatalogTagUsecase{}
}

func (s *sCatalogTagUsecase) ListTagGroups(ctx context.Context, actor *model.Actor, in catalogin.TagGroupListInp) (*catalogout.TagGroupListOut, error) {
	groups, tags, err := service.CatalogTagDomain().ListTagGroups(ctx)
	if err != nil {
		return nil, err
	}

	tagMap := make(map[uint][]catalogout.TagItem)
	for _, tag := range tags {
		var name map[string]any
		_ = tag.NameI18N.Scan(&name)
		tagMap[tag.GroupId] = append(tagMap[tag.GroupId], catalogout.TagItem{
			Id:      tag.Id,
			GroupId: tag.GroupId,
			Name:    name,
			Value:   tag.Value,
		})
	}

	var list []catalogout.TagGroupItem
	for _, group := range groups {
		var name map[string]any
		var categories []uint
		_ = group.NameI18N.Scan(&name)
		_ = group.CategoryIds.Scan(&categories)

		list = append(list, catalogout.TagGroupItem{
			Id:         group.Id,
			Name:       name,
			Slug:       group.Slug,
			Categories: categories,
			Tags:       tagMap[group.Id],
		})
	}

	return &catalogout.TagGroupListOut{List: list}, nil
}
