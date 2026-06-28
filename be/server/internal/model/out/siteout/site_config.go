package siteout

import "github.com/gogf/gf/v2/os/gtime"

type SiteConfigItem struct {
	Id          uint        `json:"id"`
	Group       string      `json:"group"`
	Key         string      `json:"key"`
	Value       any         `json:"value"`
	ValueType   string      `json:"valueType"`
	Description string      `json:"description"`
	CreatedAt   *gtime.Time `json:"createdAt"`
	UpdatedAt   *gtime.Time `json:"updatedAt"`
}

type SiteConfigListOut struct {
	Configs []*SiteConfigItem `json:"configs"`
}
