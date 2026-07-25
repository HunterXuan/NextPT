package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"server/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type TMDBProvider struct{}

type tmdbSearchResponse struct {
	Page         int              `json:"page"`
	TotalPages   int              `json:"total_pages"`
	TotalResults int              `json:"total_results"`
	Results      []tmdbSearchItem `json:"results"`
}

type tmdbSearchItem struct {
	Id            int64   `json:"id"`
	Title         string  `json:"title"`
	Name          string  `json:"name"`
	OriginalTitle string  `json:"original_title"`
	OriginalName  string  `json:"original_name"`
	ReleaseDate   string  `json:"release_date"`
	FirstAirDate  string  `json:"first_air_date"`
	Overview      string  `json:"overview"`
	PosterPath    string  `json:"poster_path"`
	BackdropPath  string  `json:"backdrop_path"`
	VoteAverage   float64 `json:"vote_average"`
}

type tmdbDetailResponse struct {
	Id            int64           `json:"id"`
	Title         string          `json:"title"`
	Name          string          `json:"name"`
	OriginalTitle string          `json:"original_title"`
	OriginalName  string          `json:"original_name"`
	ReleaseDate   string          `json:"release_date"`
	FirstAirDate  string          `json:"first_air_date"`
	Overview      string          `json:"overview"`
	PosterPath    string          `json:"poster_path"`
	BackdropPath  string          `json:"backdrop_path"`
	VoteAverage   float64         `json:"vote_average"`
	Genres        []tmdbGenre     `json:"genres"`
	ExternalIds   tmdbExternalIds `json:"external_ids"`
}

type tmdbGenre struct {
	Name string `json:"name"`
}

type tmdbExternalIds struct {
	ImdbId string `json:"imdb_id"`
}

type tmdbFindResponse struct {
	MovieResults []tmdbFindItem `json:"movie_results"`
	TvResults    []tmdbFindItem `json:"tv_results"`
}

type tmdbFindItem struct {
	Id int64 `json:"id"`
}

func (p *TMDBProvider) Name() string {
	return ProviderTMDB
}

func (p *TMDBProvider) Search(ctx context.Context, query string, mediaType string, locale string, page int) (*SearchResult, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("language", providerLocale(ctx, locale))
	params.Set("page", strconv.Itoa(page))
	params.Set("include_adult", "false")

	var response tmdbSearchResponse
	path := fmt.Sprintf("/search/%s?%s", mediaType, params.Encode())
	if err := p.request(ctx, path, &response); err != nil {
		return nil, err
	}
	items := make([]Item, 0, len(response.Results))
	for _, item := range response.Results {
		items = append(items, p.mapSearchItem(ctx, mediaType, item))
	}
	return &SearchResult{
		List:         items,
		Page:         response.Page,
		TotalPages:   response.TotalPages,
		TotalResults: response.TotalResults,
	}, nil
}

func (p *TMDBProvider) FindByIMDb(ctx context.Context, imdbId string, locale string) (string, string, error) {
	var response tmdbFindResponse
	path := fmt.Sprintf("/find/%s?external_source=imdb_id&language=%s", url.PathEscape(imdbId), url.QueryEscape(providerLocale(ctx, locale)))
	if err := p.request(ctx, path, &response); err != nil {
		return "", "", err
	}
	if len(response.MovieResults) > 0 {
		return strconv.FormatInt(response.MovieResults[0].Id, 10), consts.CatalogMetadataTmdbTypeMovie, nil
	}
	if len(response.TvResults) > 0 {
		return strconv.FormatInt(response.TvResults[0].Id, 10), consts.CatalogMetadataTmdbTypeTv, nil
	}
	return "", "", nil
}

func (p *TMDBProvider) Get(ctx context.Context, in GetInput) (*Item, error) {
	var response tmdbDetailResponse
	path := fmt.Sprintf("/%s/%s?append_to_response=external_ids&language=%s", in.MediaType, url.PathEscape(in.ExternalId), url.QueryEscape(providerLocale(ctx, in.Locale)))
	if err := p.request(ctx, path, &response); err != nil {
		return nil, err
	}
	if response.Id == 0 {
		return nil, gerror.New("tmdb metadata is empty")
	}
	metadata := p.mapDetail(ctx, in.MediaType, response)
	return &metadata, nil
}

func (p *TMDBProvider) request(ctx context.Context, path string, target any) error {
	token := strings.TrimSpace(g.Cfg().MustGet(ctx, "catalog.metadata.tmdb.token").String())
	if token == "" {
		return gerror.New("tmdb token is not configured")
	}
	baseURL := strings.TrimRight(g.Cfg().MustGet(ctx, "catalog.metadata.tmdb.baseUrl", "https://api.themoviedb.org/3").String(), "/")
	response, err := g.Client().SetTimeout(8*time.Second).
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", browserUserAgent(ctx)).
		Get(ctx, baseURL+path)
	if err != nil {
		return err
	}
	defer response.Close()
	body := response.ReadAll()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return gerror.Newf("tmdb request failed with status %d", response.StatusCode)
	}
	return json.Unmarshal(body, target)
}

func (p *TMDBProvider) mapSearchItem(ctx context.Context, mediaType string, item tmdbSearchItem) Item {
	title, originalTitle, releaseDate := item.Title, item.OriginalTitle, item.ReleaseDate
	if mediaType == consts.CatalogMetadataTmdbTypeTv {
		title, originalTitle, releaseDate = item.Name, item.OriginalName, item.FirstAirDate
	}
	return Item{
		Provider:      p.Name(),
		ProviderId:    strconv.FormatInt(item.Id, 10),
		TmdbType:      mediaType,
		Title:         title,
		OriginalTitle: originalTitle,
		Year:          year(releaseDate),
		ReleaseDate:   releaseDate,
		Overview:      item.Overview,
		PosterUrl:     tmdbImageUrl(ctx, item.PosterPath),
		BackdropUrl:   tmdbImageUrl(ctx, item.BackdropPath),
		Rating:        item.VoteAverage,
	}
}

func (p *TMDBProvider) mapDetail(ctx context.Context, mediaType string, detail tmdbDetailResponse) Item {
	title, originalTitle, releaseDate := detail.Title, detail.OriginalTitle, detail.ReleaseDate
	if mediaType == consts.CatalogMetadataTmdbTypeTv {
		title, originalTitle, releaseDate = detail.Name, detail.OriginalName, detail.FirstAirDate
	}
	genres := make([]string, 0, len(detail.Genres))
	for _, genre := range detail.Genres {
		if name := strings.TrimSpace(genre.Name); name != "" {
			genres = append(genres, name)
		}
	}
	return Item{
		Provider:      p.Name(),
		ProviderId:    strconv.FormatInt(detail.Id, 10),
		TmdbType:      mediaType,
		Title:         title,
		OriginalTitle: originalTitle,
		Year:          year(releaseDate),
		ReleaseDate:   releaseDate,
		Overview:      detail.Overview,
		PosterUrl:     tmdbImageUrl(ctx, detail.PosterPath),
		BackdropUrl:   tmdbImageUrl(ctx, detail.BackdropPath),
		Rating:        detail.VoteAverage,
		Genres:        genres,
		ImdbId:        detail.ExternalIds.ImdbId,
	}
}

func providerLocale(ctx context.Context, locale string) string {
	if locale = strings.TrimSpace(locale); locale != "" {
		return locale
	}
	return g.Cfg().MustGet(ctx, "i18n.default", "zh-CN").String()
}

func tmdbImageUrl(ctx context.Context, path string) string {
	if path == "" {
		return ""
	}
	baseURL := strings.TrimRight(g.Cfg().MustGet(ctx, "catalog.metadata.tmdb.imageBaseUrl", "https://image.tmdb.org/t/p/w500").String(), "/")
	return baseURL + "/" + strings.TrimLeft(path, "/")
}

var _ Provider = (*TMDBProvider)(nil)
