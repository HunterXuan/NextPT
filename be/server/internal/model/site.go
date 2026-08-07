package model

import (
	"net/url"
	"slices"
	"strings"

	"server/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
)

type SiteMaintenance struct {
	Enabled bool
	Message string
}

func (m SiteMaintenance) DisplayMessage(defaultMessage string) string {
	if message := strings.TrimSpace(m.Message); message != "" {
		return message
	}
	return defaultMessage
}

type SiteMaintenanceError struct {
	Message string
}

func (e *SiteMaintenanceError) Error() string {
	return e.Message
}

type SiteAdvertisementConfig struct {
	Enabled bool   `json:"enabled"`
	Title   string `json:"title"`
	Image   string `json:"image"`
	URL     string `json:"url"`
}

type SiteAdvertisements map[string]SiteAdvertisementConfig

func (ads SiteAdvertisements) Normalized() SiteAdvertisements {
	result := make(SiteAdvertisements, len(consts.SiteAdvertisementPlacements))
	for _, placement := range consts.SiteAdvertisementPlacements {
		advertisement := ads[placement]
		advertisement.Title = strings.TrimSpace(advertisement.Title)
		advertisement.Image = strings.TrimSpace(advertisement.Image)
		advertisement.URL = strings.TrimSpace(advertisement.URL)
		result[placement] = advertisement
	}
	return result
}

func (ads SiteAdvertisements) Enabled() SiteAdvertisements {
	result := make(SiteAdvertisements)
	for _, placement := range consts.SiteAdvertisementPlacements {
		if advertisement := ads[placement]; advertisement.Enabled {
			result[placement] = advertisement
		}
	}
	return result
}

func (ads SiteAdvertisements) Validate() error {
	for placement := range ads {
		if !slices.Contains(consts.SiteAdvertisementPlacements, placement) {
			return gerror.New("unsupported advertisement placement")
		}
	}

	for _, placement := range consts.SiteAdvertisementPlacements {
		advertisement := ads[placement]
		if !advertisement.Enabled {
			continue
		}
		if advertisement.Title == "" || len(advertisement.Title) > 120 {
			return gerror.New("advertisement title is required and must not exceed 120 characters")
		}
		if !isSiteAdvertisementURL(advertisement.Image) {
			return gerror.New("advertisement image must be an internal path or http(s) URL")
		}
		if !isSiteAdvertisementURL(advertisement.URL) {
			return gerror.New("advertisement URL must be an internal path or http(s) URL")
		}
	}
	return nil
}

func isSiteAdvertisementURL(value string) bool {
	if strings.HasPrefix(value, "/") {
		return !strings.HasPrefix(value, "//")
	}
	parsed, err := url.ParseRequestURI(value)
	return err == nil && parsed.Host != "" && (parsed.Scheme == "http" || parsed.Scheme == "https")
}
