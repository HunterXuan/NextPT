package adminout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

type CatalogCategoryItem struct {
	Id           uint                       `json:"id"`
	NameI18N     *gjson.Json                `json:"nameI18N"`
	Slug         string                     `json:"slug"`
	SortOrder    int                        `json:"sortOrder"`
	Enabled      bool                       `json:"enabled"`
	UploadConfig *model.CatalogUploadConfig `json:"uploadConfig"`
	CreatedAt    *gtime.Time                `json:"createdAt"`
	UpdatedAt    *gtime.Time                `json:"updatedAt"`
}

type CatalogCategoryListOut struct {
	Categories []CatalogCategoryItem `json:"categories"`
}
