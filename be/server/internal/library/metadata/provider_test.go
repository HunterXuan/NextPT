package metadata

import "testing"

func TestParseIMDb(t *testing.T) {
	body := []byte(`<!doctype html><html><head>
<script type="application/ld+json">{
  "@context": "https://schema.org",
  "@type": "Movie",
  "name": "The Shawshank Redemption",
  "alternateName": "肖申克的救赎",
  "datePublished": "1994-10-14",
  "description": "Two imprisoned men bond over a number of years.",
  "image": "https://m.media-amazon.com/images/example.jpg",
  "genre": ["Drama"],
  "aggregateRating": {"ratingValue": 9.3}
}</script></head></html>`)

	item, err := parseIMDb("tt0111161", body)
	if err != nil {
		t.Fatalf("parseIMDb() error = %v", err)
	}
	if item.Provider != ProviderIMDb || item.ProviderId != "tt0111161" || item.ImdbId != "tt0111161" {
		t.Fatalf("unexpected provider identity: %+v", item)
	}
	if item.Title != "The Shawshank Redemption" || item.OriginalTitle != "肖申克的救赎" || item.Year != "1994" {
		t.Fatalf("unexpected title metadata: %+v", item)
	}
	if item.TmdbType != "movie" || item.Rating != 9.3 || len(item.Genres) != 1 {
		t.Fatalf("unexpected media metadata: %+v", item)
	}
}

func TestParseIMDbFromGraph(t *testing.T) {
	body := []byte(`<!doctype html><html><head>
<script type="application/ld+json">{
  "@graph": [
    {"@type": "WebSite", "name": "IMDb"},
    {
      "@type": ["TVSeries", "CreativeWork"],
      "name": "Example Series",
      "datePublished": "2025",
      "aggregateRating": {"ratingValue": "8.6"}
    }
  ]
}</script></head></html>`)

	item, err := parseIMDb("tt1234567", body)
	if err != nil {
		t.Fatalf("parseIMDb() error = %v", err)
	}
	if item.TmdbType != "tv" || item.Rating != 8.6 || item.Year != "2025" {
		t.Fatalf("unexpected graph metadata: %+v", item)
	}
}

func TestParseIMDbGraphQL(t *testing.T) {
	body := []byte(`{
  "data": {
    "title": {
      "id": "tt0111161",
      "titleText": {"text": "The Shawshank Redemption"},
      "originalTitleText": {"text": "The Shawshank Redemption"},
      "releaseYear": {"year": 1994},
      "releaseDate": {"day": 14, "month": 10, "year": 1994},
      "ratingsSummary": {"aggregateRating": 9.3},
      "plot": {"plotText": {"plainText": "A banker is sent to prison."}},
      "primaryImage": {"url": "https://m.media-amazon.com/example.jpg"},
      "titleGenres": {"genres": [{"genre": {"text": "Drama"}}]},
      "titleType": {"id": "movie", "text": "Movie"}
    }
  }
}`)

	item, err := parseIMDbGraphQL("tt0111161", body)
	if err != nil {
		t.Fatalf("parseIMDbGraphQL() error = %v", err)
	}
	if item.Title != "The Shawshank Redemption" || item.Year != "1994" || item.ReleaseDate != "1994-10-14" {
		t.Fatalf("unexpected title metadata: %+v", item)
	}
	if item.Rating != 9.3 || item.TmdbType != "movie" || len(item.Genres) != 1 {
		t.Fatalf("unexpected media metadata: %+v", item)
	}
}

func TestParseIMDbRejectsMissingOrMalformedJsonLd(t *testing.T) {
	for _, body := range [][]byte{
		[]byte(`<html><head></head><body>empty</body></html>`),
		[]byte(`<html><script type="application/ld+json">{</script></html>`),
		[]byte(`<html><script type="application/ld+json">{"@type":"Movie"}</script></html>`),
	} {
		if _, err := parseIMDb("tt0000001", body); err == nil {
			t.Fatalf("parseIMDb() expected error for %q", body)
		}
	}
}

func TestParseDoubanObject(t *testing.T) {
	body := []byte(`<!doctype html><html><head>
<script TYPE="Application/LD+JSON">{
  "@type": "Movie",
  "name": "测试电影",
  "alternateName": "Test Movie",
  "datePublished": "2024-06-01",
  "description": "简介",
  "image": "https://img.example/poster.jpg",
  "genre": ["剧情", "悬疑"],
  "aggregateRating": {"ratingValue": "8.7"},
  "sameAs": ["https://www.imdb.com/title/tt1234567/"]
}</script></head></html>`)

	item, err := parseDouban("1292052", body)
	if err != nil {
		t.Fatalf("parseDouban() error = %v", err)
	}
	if item.Provider != ProviderDouban || item.ProviderId != "1292052" {
		t.Fatalf("unexpected provider identity: %+v", item)
	}
	if item.Title != "测试电影" || item.OriginalTitle != "Test Movie" || item.Year != "2024" {
		t.Fatalf("unexpected title metadata: %+v", item)
	}
	if item.Rating != 8.7 || item.ImdbId != "tt1234567" {
		t.Fatalf("unexpected rating or imdb id: %+v", item)
	}
	if len(item.Genres) != 2 {
		t.Fatalf("unexpected genres: %+v", item.Genres)
	}
}

func TestParseDoubanJSON(t *testing.T) {
	body := []byte(`{
  "id": "1292052",
  "type": "movie",
  "subtype": "movie",
  "title": "肖申克的救赎",
  "original_title": "The Shawshank Redemption",
  "year": "1994",
  "intro": "一段剧情简介。",
  "cover_url": "https://img.example/poster.jpg",
  "genres": ["剧情", "犯罪"],
  "is_tv": false,
  "rating": {"value": 9.7}
}`)

	item, err := parseDoubanJSON("1292052", body)
	if err != nil {
		t.Fatalf("parseDoubanJSON() error = %v", err)
	}
	if item.Title != "肖申克的救赎" || item.OriginalTitle != "The Shawshank Redemption" || item.Year != "1994" {
		t.Fatalf("unexpected title metadata: %+v", item)
	}
	if item.Rating != 9.7 || item.Overview != "一段剧情简介。" || item.PosterUrl != "https://img.example/poster.jpg" || len(item.Genres) != 2 {
		t.Fatalf("unexpected json metadata: %+v", item)
	}
}

func TestParseDoubanGraphAndNumericRating(t *testing.T) {
	body := []byte(`<!doctype html><html><head>
<script type="application/ld+json">{
  "@graph": [
    {"@type": "WebSite", "name": "豆瓣电影", "url": "https://movie.douban.com"},
    {
      "@type": "TVSeries",
      "name": "测试剧集",
      "datePublished": "2023",
      "image": {"url": "https://img.example/tv.jpg"},
      "aggregateRating": {"ratingValue": 9.1}
    }
  ]
}</script></head><body><span>IMDb: tt7654321</span></body></html>`)

	item, err := parseDouban("36000000", body)
	if err != nil {
		t.Fatalf("parseDouban() error = %v", err)
	}
	if item.TmdbType != "tv" || item.Rating != 9.1 || item.ImdbId != "tt7654321" {
		t.Fatalf("unexpected graph metadata: %+v", item)
	}
	if item.PosterUrl != "https://img.example/tv.jpg" {
		t.Fatalf("unexpected poster: %q", item.PosterUrl)
	}
}

func TestParseDoubanMobilePage(t *testing.T) {
	body := []byte(`<!doctype html><html><head>
<meta itemprop="image" content="https://img.example/poster.jpg">
<meta itemprop="ratingValue" content="9.7">
</head><body>
<div class="sub-title">肖申克的救赎</div>
<div class="sub-original-title">The Shawshank Redemption（1994）</div>
<section class="subject-mark mark-movie"></section>
<section class="subject-intro"><div class="bd"><p>一段剧情简介。</p></div></section>
</body></html>`)

	item, err := parseDouban("1292052", body)
	if err != nil {
		t.Fatalf("parseDouban() error = %v", err)
	}
	if item.Title != "肖申克的救赎" || item.OriginalTitle != "The Shawshank Redemption" || item.Year != "1994" {
		t.Fatalf("unexpected title metadata: %+v", item)
	}
	if item.Rating != 9.7 || item.Overview != "一段剧情简介。" || item.PosterUrl != "https://img.example/poster.jpg" {
		t.Fatalf("unexpected mobile metadata: %+v", item)
	}
}

func TestParseDoubanRejectsMissingOrMalformedJsonLd(t *testing.T) {
	tests := [][]byte{
		[]byte(`<html><head></head><body>empty</body></html>`),
		[]byte(`<html><script type="application/ld+json">{</script></html>`),
		[]byte(`<html><script type="application/ld+json">{"@type":"Movie"}</script></html>`),
	}
	for _, body := range tests {
		if _, err := parseDouban("1", body); err == nil {
			t.Fatalf("parseDouban() expected error for %q", body)
		}
	}
}

func TestParseBangumi(t *testing.T) {
	body := []byte(`{
  "id": 326,
  "type": 2,
  "name": "Cowboy Bebop",
  "name_cn": "星际牛仔",
  "summary": "未来世界的赏金猎人故事。",
  "date": "1998-04-03",
  "images": {"large": "//lain.bgm.tv/pic/cover/l/example.jpg"},
  "rating": {"score": 9.1},
  "tags": [
    {"name": "科幻"},
    {"name": "太空"},
    {"name": "科幻"}
  ],
  "infobox": [
    {"key": "别名", "value": [{"v": "カウボーイビバップ"}]},
    {"key": "IMDb", "value": [{"v": "https://www.imdb.com/title/tt0213338/"}]}
  ]
}`)

	item, err := parseBangumi("326", body)
	if err != nil {
		t.Fatalf("parseBangumi() error = %v", err)
	}
	if item.Provider != ProviderBangumi || item.ProviderId != "326" {
		t.Fatalf("unexpected provider identity: %+v", item)
	}
	if item.Title != "星际牛仔" || item.OriginalTitle != "Cowboy Bebop" || item.Year != "1998" {
		t.Fatalf("unexpected title metadata: %+v", item)
	}
	if item.Rating != 9.1 || item.ImdbId != "tt0213338" {
		t.Fatalf("unexpected rating or imdb id: %+v", item)
	}
	if item.PosterUrl != "https://lain.bgm.tv/pic/cover/l/example.jpg" {
		t.Fatalf("unexpected poster url: %q", item.PosterUrl)
	}
	if len(item.Genres) != 2 {
		t.Fatalf("unexpected tags: %+v", item.Genres)
	}
}

func TestParseBangumiRejectsInvalidResponse(t *testing.T) {
	for _, body := range [][]byte{
		[]byte(`{`),
		[]byte(`{"id": 0, "name": ""}`),
	} {
		if _, err := parseBangumi("1", body); err == nil {
			t.Fatalf("parseBangumi() expected error for %q", body)
		}
	}
}
