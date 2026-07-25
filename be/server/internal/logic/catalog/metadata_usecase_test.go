package catalog

import (
	"testing"

	libmetadata "server/internal/library/metadata"
	"server/internal/model/out/catalogout"
)

func TestMergeMetadataUsesFixedProviderPriority(t *testing.T) {
	usecase := NewCatalogMetadataUsecase()
	tmdb := &catalogout.TorrentMetadataItem{
		Provider:    libmetadata.ProviderTMDB,
		Title:       "TMDB title",
		Rating:      7.5,
		PosterUrl:   "https://tmdb.example/poster.jpg",
		ReleaseDate: "",
	}
	imdb := &catalogout.TorrentMetadataItem{
		Provider:    libmetadata.ProviderIMDb,
		Title:       "IMDb title",
		ReleaseDate: "2024-01-01",
		Overview:    "IMDb overview",
		Rating:      8.5,
	}
	douban := &catalogout.TorrentMetadataItem{
		Provider:      libmetadata.ProviderDouban,
		Title:         "豆瓣标题",
		OriginalTitle: "Douban original title",
		Rating:        9.0,
		PosterUrl:     "https://douban.example/poster.jpg",
	}
	bangumi := &catalogout.TorrentMetadataItem{
		Provider:      libmetadata.ProviderBangumi,
		OriginalTitle: "Bangumi original title",
		Genres:        []string{"anime"},
		Rating:        9.2,
	}

	merged := usecase.mergeMetadata(tmdb, imdb, douban, bangumi)
	if merged == nil {
		t.Fatal("mergeMetadata() returned nil")
	}
	if merged.Provider != libmetadata.ProviderTMDB || merged.Title != "TMDB title" || merged.PosterUrl != tmdb.PosterUrl || merged.Rating != 7.5 {
		t.Fatalf("TMDB fields should keep priority: %+v", merged)
	}
	if merged.ReleaseDate != imdb.ReleaseDate || merged.Overview != imdb.Overview {
		t.Fatalf("IMDb should fill fields missing after TMDB: %+v", merged)
	}
	if merged.OriginalTitle != douban.OriginalTitle {
		t.Fatalf("Douban should fill fields still missing after IMDb: %+v", merged)
	}
	if len(merged.Genres) != 1 || merged.Genres[0] != "anime" {
		t.Fatalf("Bangumi should fill fields still missing after Douban: %+v", merged)
	}
}
