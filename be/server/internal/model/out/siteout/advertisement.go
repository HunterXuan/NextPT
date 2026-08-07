package siteout

import "server/internal/model"

type AdvertisementListOut struct {
	Placements model.SiteAdvertisements `json:"placements"`
}
