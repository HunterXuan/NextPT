package model

import (
	"testing"

	"server/internal/consts"
)

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
			Enabled:     true,
			Title:       "  Partner  ",
			Image:       " https://cdn.example.com/banner.webp ",
			URL:         " /catalog/torrents ",
			AspectRatio: " 16 : 9 ",
		},
	}

	advertisements = advertisements.Normalized()
	if err := advertisements.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if got := advertisements["home"]; got.Title != "Partner" || got.Image != "https://cdn.example.com/banner.webp" || got.URL != "/catalog/torrents" || got.AspectRatio != "16:9" {
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
		{
			name: "invalid aspect ratio",
			ads:  SiteAdvertisements{"home": {AspectRatio: "wide"}},
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

func TestSiteTasksNormalizeAndValidate(t *testing.T) {
	tasks := SiteTasks{{
		Key:      "  weekly_seed  ",
		Enabled:  true,
		Cycle:    consts.SiteTaskCycleWeekly,
		NameI18N: map[string]string{"zh-CN": "  每周保种  "},
		Rule:     SiteTaskRule{Type: consts.SiteTaskRuleTypeSeedDuration, Target: 12},
		Rewards:  []SiteTaskReward{{Type: consts.SiteTaskRewardTypeBonus, Amount: 100.04}},
	}}

	tasks = tasks.Normalized()
	if err := tasks.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if tasks[0].Key != "weekly_seed" || tasks[0].NameI18N["zh-CN"] != "每周保种" || tasks[0].Rewards[0].Amount != 100 {
		t.Fatalf("Normalized() task = %#v", tasks[0])
	}
}

func TestSiteTasksRejectInvalidRuleAndDuplicateReward(t *testing.T) {
	tasks := SiteTasks{{
		Key:      "bad_task",
		Cycle:    consts.SiteTaskCycleOnce,
		NameI18N: map[string]string{"zh-CN": "无效任务"},
		Rule:     SiteTaskRule{Type: consts.SiteTaskRuleTypeUploaded, Target: 1},
		Rewards: []SiteTaskReward{
			{Type: consts.SiteTaskRewardTypeInvite, Amount: 1},
			{Type: consts.SiteTaskRewardTypeInvite, Amount: 1},
		},
	}}

	if err := tasks.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want invalid task error")
	}
}
