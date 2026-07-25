package metadata

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

func TestProvidersLive(t *testing.T) {
	if os.Getenv("NEXTP_METADATA_LIVE_TEST") != "1" {
		t.Skip("set NEXTP_METADATA_LIVE_TEST=1 to run external provider tests")
	}

	metadataConfig := map[string]any{
		"userAgent": defaultBrowserUserAgent,
		"imdb": map[string]any{
			"baseUrl":    "https://www.imdb.com",
			"graphqlUrl": "https://api.graphql.imdb.com/",
		},
		"douban": map[string]any{
			"baseUrl":    "https://m.douban.com/movie",
			"apiBaseUrl": "https://m.douban.com/rexxar/api/v2",
		},
		"bangumi": map[string]any{
			"baseUrl": "https://api.bgm.tv/v0",
		},
	}
	if token := os.Getenv("TMDB_API_TOKEN"); token != "" {
		metadataConfig["tmdb"] = map[string]any{
			"token":   token,
			"baseUrl": "https://api.themoviedb.org/3",
		}
	}
	content, err := json.Marshal(map[string]any{
		"catalog": map[string]any{"metadata": metadataConfig},
	})
	if err != nil {
		t.Fatal(err)
	}
	adapter, err := gcfg.NewAdapterContent(string(content))
	if err != nil {
		t.Fatal(err)
	}
	originalAdapter := g.Cfg().GetAdapter()
	g.Cfg().SetAdapter(adapter)
	t.Cleanup(func() { g.Cfg().SetAdapter(originalAdapter) })

	tests := []struct {
		name     string
		provider Provider
		input    GetInput
		expected string
	}{
		{name: "imdb", provider: &IMDbProvider{}, input: GetInput{ExternalId: "tt0111161"}, expected: ProviderIMDb},
		{name: "douban", provider: &DoubanProvider{}, input: GetInput{ExternalId: "1292052"}, expected: ProviderDouban},
		{name: "bangumi", provider: &BangumiProvider{}, input: GetInput{ExternalId: "326"}, expected: ProviderBangumi},
	}
	if os.Getenv("TMDB_API_TOKEN") != "" {
		tests = append(tests, struct {
			name     string
			provider Provider
			input    GetInput
			expected string
		}{name: "tmdb", provider: &TMDBProvider{}, input: GetInput{ExternalId: "278", MediaType: "movie", Locale: "en-US"}, expected: ProviderTMDB})
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			item, err := test.provider.Get(ctx, test.input)
			if err != nil {
				t.Fatalf("provider request failed: %v", err)
			}
			if item == nil || item.Provider != test.expected || item.ProviderId == "" || item.Title == "" {
				t.Fatalf("provider returned incomplete metadata: %+v", item)
			}
			if item.Rating <= 0 {
				t.Fatalf("provider returned invalid rating: %+v", item)
			}
		})
	}
}
