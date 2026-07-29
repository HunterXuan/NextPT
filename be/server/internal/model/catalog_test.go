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

func TestCatalogTorrentListOptionsNormalized(t *testing.T) {
	got := (CatalogTorrentListOptions{
		Keyword:             "  test  ",
		CategoryIds:         []uint{3, 0, 3, 2},
		TagGroups:           [][]uint{{5, 0, 5, 6}, {}, {9}},
		Promotion:           "missing",
		SeedStatus:          "missing",
		PublishedWithinDays: -1,
		Sort:                "missing",
		Page:                0,
		Size:                500,
	}).Normalized()

	if got.Keyword != "test" {
		t.Fatalf("keyword = %q, want test", got.Keyword)
	}
	if len(got.CategoryIds) != 2 || got.CategoryIds[0] != 3 || got.CategoryIds[1] != 2 {
		t.Fatalf("category ids = %v, want [3 2]", got.CategoryIds)
	}
	if len(got.TagGroups) != 2 || len(got.TagGroups[0]) != 2 || got.TagGroups[0][0] != 5 || got.TagGroups[0][1] != 6 || len(got.TagGroups[1]) != 1 || got.TagGroups[1][0] != 9 {
		t.Fatalf("tag groups = %v, want [[5 6] [9]]", got.TagGroups)
	}
	if got.Promotion != consts.CatalogTorrentPromotionFilterAll {
		t.Fatalf("promotion = %q, want all", got.Promotion)
	}
	if got.SeedStatus != consts.CatalogTorrentSeedStatusAll {
		t.Fatalf("seed status = %q, want all", got.SeedStatus)
	}
	if got.PublishedWithinDays != 0 {
		t.Fatalf("published within = %d, want 0", got.PublishedWithinDays)
	}
	if got.Sort != consts.CatalogTorrentSortNewest {
		t.Fatalf("sort = %q, want newest", got.Sort)
	}
	if got.Page != 1 || got.Size != 100 {
		t.Fatalf("page/size = %d/%d, want 1/100", got.Page, got.Size)
	}
}

func TestCatalogTorrentListOptionsInvalidSizeRange(t *testing.T) {
	if !(CatalogTorrentListOptions{MinSize: 2, MaxSize: 1}).HasInvalidSizeRange() {
		t.Fatal("expected invalid size range")
	}
	if (CatalogTorrentListOptions{MinSize: 1, MaxSize: 2}).HasInvalidSizeRange() {
		t.Fatal("expected valid size range")
	}
}

func TestCatalogTorrentMetadataFilterNormalized(t *testing.T) {
	got := (CatalogTorrentMetadataFilter{
		ImdbId:    "  TT1234567 ",
		DoubanId:  " 1295644 ",
		BangumiId: " 2 ",
		TmdbId:    " 550 ",
		TmdbType:  " MOVIE ",
	}).Normalized()

	if got.ImdbId != "tt1234567" || got.DoubanId != "1295644" || got.BangumiId != "2" || got.TmdbId != "550" || got.TmdbType != "movie" {
		t.Fatalf("normalized metadata filter = %+v", got)
	}
	if !got.IsValid() || got.IsEmpty() {
		t.Fatalf("expected valid non-empty metadata filter: %+v", got)
	}
}

func TestCatalogTorrentMetadataFilterValidation(t *testing.T) {
	tests := []struct {
		name   string
		filter CatalogTorrentMetadataFilter
		valid  bool
	}{
		{name: "empty", filter: CatalogTorrentMetadataFilter{}, valid: true},
		{name: "imdb", filter: CatalogTorrentMetadataFilter{ImdbId: "tt1234567"}, valid: true},
		{name: "tmdb without type", filter: CatalogTorrentMetadataFilter{TmdbId: "550"}, valid: true},
		{name: "tmdb movie", filter: CatalogTorrentMetadataFilter{TmdbId: "550", TmdbType: "movie"}, valid: true},
		{name: "invalid imdb", filter: CatalogTorrentMetadataFilter{ImdbId: "1234567"}, valid: false},
		{name: "invalid numeric id", filter: CatalogTorrentMetadataFilter{DoubanId: "movie-1"}, valid: false},
		{name: "type without tmdb id", filter: CatalogTorrentMetadataFilter{TmdbType: "tv"}, valid: false},
		{name: "invalid tmdb type", filter: CatalogTorrentMetadataFilter{TmdbId: "550", TmdbType: "anime"}, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %t, want %t", got, tt.valid)
			}
		})
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

func TestPickCatalogNewTorrentPromotionAllowsNoExpire(t *testing.T) {
	now := gtime.NewFromStr("2026-07-08 12:00:00")
	got := PickCatalogNewTorrentPromotion(CatalogTorrentNewPromotionConfig{
		Enabled: true,
		Rules: []CatalogTorrentPromotionRule{
			{
				MinGiB:        0,
				DurationHours: 0,
				Options: []CatalogTorrentPromotionOption{
					{State: consts.ResourceTorrentPromotionStateFree, Weight: 100},
				},
			},
		},
	}, catalogBytesPerGiB, now)

	if got.SpState != consts.ResourceTorrentSpFree {
		t.Fatalf("new torrent sp = %d, want %d", got.SpState, consts.ResourceTorrentSpFree)
	}
	if got.SpExpireAt != nil {
		t.Fatalf("expire at = %v, want nil", got.SpExpireAt)
	}
}

func TestPickCatalogNewTorrentPromotionRejectsNegativeDuration(t *testing.T) {
	now := gtime.NewFromStr("2026-07-08 12:00:00")
	got := PickCatalogNewTorrentPromotion(CatalogTorrentNewPromotionConfig{
		Enabled: true,
		Rules: []CatalogTorrentPromotionRule{
			{
				MinGiB:        0,
				DurationHours: -1,
				Options: []CatalogTorrentPromotionOption{
					{State: consts.ResourceTorrentPromotionStateFree, Weight: 100},
				},
			},
		},
	}, catalogBytesPerGiB, now)

	if got.SpState != consts.ResourceTorrentSpNormal {
		t.Fatalf("new torrent sp = %d, want %d", got.SpState, consts.ResourceTorrentSpNormal)
	}
	if got.SpExpireAt != nil {
		t.Fatalf("expire at = %v, want nil", got.SpExpireAt)
	}
}
