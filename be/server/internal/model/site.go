package model

import "strings"

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
