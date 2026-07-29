package adminout

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

type CatalogTagItem struct {
	Id        uint        `json:"id"`
	GroupId   uint        `json:"groupId"`
	NameI18N  *gjson.Json `json:"nameI18N"`
	Value     string      `json:"value"`
	SortOrder int         `json:"sortOrder"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

type CatalogTagGroupItem struct {
	Id          uint             `json:"id"`
	NameI18N    *gjson.Json      `json:"nameI18N"`
	Slug        string           `json:"slug"`
	CategoryIds []uint           `json:"categoryIds"`
	SortOrder   int              `json:"sortOrder"`
	Tags        []CatalogTagItem `json:"tags"`
	CreatedAt   *gtime.Time      `json:"createdAt"`
	UpdatedAt   *gtime.Time      `json:"updatedAt"`
}

type CatalogTagGroupListOut struct {
	Groups []CatalogTagGroupItem `json:"groups"`
}
