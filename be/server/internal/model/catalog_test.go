package model

import (
	"testing"

	"server/internal/consts"

	"github.com/gogf/gf/v2/os/gtime"
)

func TestCatalogTorrentPromotionStateToSp(t *testing.T) {
	tests := []struct {
		state string
		want  int
	}{
		{state: consts.ResourceTorrentPromotionStateNormal, want: consts.ResourceTorrentSpNormal},
		{state: consts.ResourceTorrentPromotionStateFree, want: consts.ResourceTorrentSpFree},
		{state: consts.ResourceTorrentPromotionState2x, want: consts.ResourceTorrentSp2x},
		{state: consts.ResourceTorrentPromotionState2xFree, want: consts.ResourceTorrentSp2xFree},
		{state: consts.ResourceTorrentPromotionState50Percent, want: consts.ResourceTorrentSp50Off},
		{state: consts.ResourceTorrentPromotionState2x50Percent, want: consts.ResourceTorrentSp2x50Off},
		{state: consts.ResourceTorrentPromotionState30Percent, want: consts.ResourceTorrentSp30Off},
		{state: "missing", want: consts.ResourceTorrentSpNormal},
	}

	for _, tt := range tests {
		if got := CatalogTorrentPromotionStateToSp(tt.state); got != tt.want {
			t.Fatalf("CatalogTorrentPromotionStateToSp(%q) = %d, want %d", tt.state, got, tt.want)
		}
	}
}

func TestResolveCatalogTorrentPromotionGlobalOverridesTorrent(t *testing.T) {
	now := gtime.NewFromStr("2026-07-08 12:00:00")
	torrentExpireAt := now.AddDate(0, 0, 1)

	got := ResolveCatalogTorrentPromotion(
		consts.ResourceTorrentSp2x,
		torrentExpireAt,
		CatalogTorrentGlobalPromotionConfig{
			Enabled:  true,
			State:    consts.ResourceTorrentPromotionStateFree,
			ExpireAt: "2026-07-09 12:00:00",
		},
		now,
	)

	if got.SpState != consts.ResourceTorrentSpFree {
		t.Fatalf("effective sp = %d, want %d", got.SpState, consts.ResourceTorrentSpFree)
	}
}

func TestApplyCatalogTorrentPromotionToTraffic(t *testing.T) {
	uploaded, downloaded := (CatalogTorrentPromotion{SpState: consts.ResourceTorrentSp2x50Off}).ApplyTraffic(100, 100)
	if uploaded != 200 || downloaded != 50 {
		t.Fatalf("credited traffic = %d/%d, want 200/50", uploaded, downloaded)
	}
}

func TestPickCatalogNewTorrentPromotionUsesLargestSizeRule(t *testing.T) {
	now := gtime.NewFromStr("2026-07-08 12:00:00")
	got := PickCatalogNewTorrentPromotion(CatalogTorrentNewPromotionConfig{
		Enabled: true,
		Rules: []CatalogTorrentPromotionRule{
			{
				MinGiB:        0,
				DurationHours: 12,
				Options: []CatalogTorrentPromotionOption{
					{State: consts.ResourceTorrentPromotionStateNormal, Weight: 100},
				},
			},
			{
				MinGiB:        10,
				DurationHours: 24,
				Options: []CatalogTorrentPromotionOption{
					{State: consts.ResourceTorrentPromotionStateFree, Weight: 100},
				},
			},
			{
				MinGiB:        50,
				DurationHours: 48,
				Options: []CatalogTorrentPromotionOption{
					{State: consts.ResourceTorrentPromotionState2xFree, Weight: 100},
				},
			},
		},
	}, 60*catalogBytesPerGiB, now)

	if got.SpState != consts.ResourceTorrentSp2xFree {
		t.Fatalf("new torrent sp = %d, want %d", got.SpState, consts.ResourceTorrentSp2xFree)
	}
	if got.SpExpireAt == nil {
		t.Fatal("new torrent promotion should have expire time")
	}
	if got.SpExpireAt.Unix()-now.Unix() != 48*60*60 {
		t.Fatalf("duration = %d seconds, want %d", got.SpExpireAt.Unix()-now.Unix(), 48*60*60)
	}
}
