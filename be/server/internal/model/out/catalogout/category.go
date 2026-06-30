package catalogout

import "server/internal/model"

type CategoryItem struct {
	Id           uint                       `json:"id"`
	Name         map[string]any             `json:"name"`
	Slug         string                     `json:"slug"`
	UploadConfig *model.CatalogUploadConfig `json:"uploadConfig"`
}

type CategoryListOut struct {
	List []CategoryItem `json:"list"`
}

type TagItem struct {
	Id    uint           `json:"id"`
	Name  map[string]any `json:"name"`
	Value string         `json:"value"`
}

type TagGroupItem struct {
	Id         uint           `json:"id"`
	Name       map[string]any `json:"name"`
	Slug       string         `json:"slug"`
	Categories []uint         `json:"categories"`
	Tags       []TagItem      `json:"tags"`
}

type TagGroupListOut struct {
	List []TagGroupItem `json:"list"`
}
