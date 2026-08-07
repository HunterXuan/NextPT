package model

import "testing"

func TestSiteMaintenanceDisplayMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{name: "custom message", message: "Scheduled maintenance", want: "Scheduled maintenance"},
		{name: "trimmed custom message", message: "  Scheduled maintenance  ", want: "Scheduled maintenance"},
		{name: "default message", message: "  ", want: "Please try again later."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maintenance := SiteMaintenance{Message: tt.message}
			if got := maintenance.DisplayMessage("Please try again later."); got != tt.want {
				t.Fatalf("DisplayMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}
