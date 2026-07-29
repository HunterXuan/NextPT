package catalog

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"
	"server/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
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
			Id:           e.Id,
			Slug:         e.Slug,
			UploadConfig: s.scanUploadConfig(e.UploadConfig),
		}
		_ = e.NameI18N.Scan(&item.Name)
		list = append(list, item)
	}

	return &catalogout.CategoryListOut{List: list}, nil
}

func (s *sCatalogCategoryUsecase) scanUploadConfig(value *gjson.Json) *model.CatalogUploadConfig {
	if value == nil {
		return nil
	}
	var config model.CatalogUploadConfig
	if err := value.Scan(&config); err != nil {
		return nil
	}
	return &config
}
