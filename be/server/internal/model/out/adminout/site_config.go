package adminout

import "server/internal/model/entity"

type SiteConfigListOut struct {
	Configs []*entity.SiteConfig `json:"configs"`
}
