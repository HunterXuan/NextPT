package adminout

import "server/internal/model/out/siteout"

type SiteConfigItem = siteout.SiteConfigItem

type SiteConfigListOut struct {
	Configs []*SiteConfigItem `json:"configs"`
}
