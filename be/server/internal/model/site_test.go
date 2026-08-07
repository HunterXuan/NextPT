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

func TestSiteAdvertisementsNormalizeAndValidate(t *testing.T) {
	advertisements := SiteAdvertisements{
		"home": {
			Enabled: true,
			Title:   "  Partner  ",
			Image:   " https://cdn.example.com/banner.webp ",
			URL:     " /catalog/torrents ",
		},
	}

	advertisements = advertisements.Normalized()
	if err := advertisements.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if got := advertisements["home"]; got.Title != "Partner" || got.Image != "https://cdn.example.com/banner.webp" || got.URL != "/catalog/torrents" {
		t.Fatalf("Normalized() home = %#v", got)
	}
	if len(advertisements.Enabled()) != 1 {
		t.Fatalf("Enabled() length = %d, want 1", len(advertisements.Enabled()))
	}
}

func TestSiteAdvertisementsRejectInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name string
		ads  SiteAdvertisements
	}{
		{
			name: "unsupported placement",
			ads:  SiteAdvertisements{"topic": {Enabled: true}},
		},
		{
			name: "missing enabled fields",
			ads:  SiteAdvertisements{"home": {Enabled: true, Title: "Partner"}},
		},
		{
			name: "unsafe URL",
			ads: SiteAdvertisements{"home": {
				Enabled: true,
				Title:   "Partner",
				Image:   "https://cdn.example.com/banner.webp",
				URL:     "javascript:alert(1)",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ads.Validate(); err == nil {
				t.Fatal("Validate() error = nil, want invalid configuration error")
			}
		})
	}
}
