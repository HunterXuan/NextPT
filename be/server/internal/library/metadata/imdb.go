package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"server/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/net/html"
)

const imdbMaxBodySize = 3 * 1024 * 1024

const imdbGraphQLQuery = `query TitleMetadata($id: ID!) {
  title(id: $id) {
    id
    titleText { text }
    originalTitleText { text }
    releaseYear { year }
    releaseDate { day month year }
    ratingsSummary { aggregateRating }
    plot { plotText { plainText } }
    primaryImage { url }
    titleGenres { genres { genre { text } } }
    titleType { id text }
  }
}`

type IMDbProvider struct{}

type imdbGraphQLResponse struct {
	Data struct {
		Title *imdbGraphQLTitle `json:"title"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type imdbGraphQLTitle struct {
	Id                string                 `json:"id"`
	TitleText         imdbGraphQLText        `json:"titleText"`
	OriginalTitleText imdbGraphQLText        `json:"originalTitleText"`
	ReleaseYear       imdbGraphQLYear        `json:"releaseYear"`
	ReleaseDate       imdbGraphQLReleaseDate `json:"releaseDate"`
	RatingsSummary    imdbGraphQLRating      `json:"ratingsSummary"`
	Plot              imdbGraphQLPlot        `json:"plot"`
	PrimaryImage      imdbGraphQLImage       `json:"primaryImage"`
	TitleGenres       imdbGraphQLGenres      `json:"titleGenres"`
	TitleType         imdbGraphQLTitleType   `json:"titleType"`
}

type imdbGraphQLText struct {
	Text string `json:"text"`
}

type imdbGraphQLYear struct {
	Year int `json:"year"`
}

type imdbGraphQLReleaseDate struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type imdbGraphQLRating struct {
	AggregateRating float64 `json:"aggregateRating"`
}

type imdbGraphQLPlot struct {
	PlotText struct {
		PlainText string `json:"plainText"`
	} `json:"plotText"`
}

type imdbGraphQLImage struct {
	Url string `json:"url"`
}

type imdbGraphQLGenres struct {
	Genres []struct {
		Genre imdbGraphQLText `json:"genre"`
	} `json:"genres"`
}

type imdbGraphQLTitleType struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

func (p *IMDbProvider) Name() string {
	return ProviderIMDb
}

func (p *IMDbProvider) Get(ctx context.Context, in GetInput) (*Item, error) {
	externalId := strings.ToLower(strings.TrimSpace(in.ExternalId))
	if !imdbIdPattern.MatchString(externalId) {
		return nil, gerror.New("invalid imdb id")
	}
	item, graphQLErr := p.fetchGraphQL(ctx, externalId)
	if graphQLErr == nil {
		return item, nil
	}
	body, err := p.fetchPage(ctx, externalId)
	if err != nil {
		return nil, gerror.Newf("imdb graphql request failed: %v; page fallback failed: %v", graphQLErr, err)
	}
	pageItem, err := parseIMDb(externalId, body)
	if err != nil {
		return nil, gerror.Newf("imdb graphql request failed: %v; page fallback failed: %v", graphQLErr, err)
	}
	return &pageItem, nil
}

func (p *IMDbProvider) fetchGraphQL(ctx context.Context, externalId string) (*Item, error) {
	endpoint := strings.TrimSpace(g.Cfg().MustGet(ctx, "catalog.metadata.imdb.graphqlUrl", "https://api.graphql.imdb.com/").String())
	response, err := g.Client().SetTimeout(8*time.Second).
		SetHeader("Accept", "application/json").
		SetHeader("Origin", "https://www.imdb.com").
		SetHeader("Referer", "https://www.imdb.com/").
		SetHeader("User-Agent", browserUserAgent(ctx)).
		ContentJson().
		Post(ctx, endpoint, g.Map{
			"query":     imdbGraphQLQuery,
			"variables": g.Map{"id": externalId},
		})
	if err != nil {
		return nil, err
	}
	defer response.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, gerror.Newf("imdb graphql request failed with status %d", response.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, imdbMaxBodySize+1))
	if err != nil {
		return nil, err
	}
	if len(body) > imdbMaxBodySize {
		return nil, gerror.New("imdb graphql response is too large")
	}
	return parseIMDbGraphQL(externalId, body)
}

func (p *IMDbProvider) fetchPage(ctx context.Context, externalId string) ([]byte, error) {
	baseURL := strings.TrimRight(g.Cfg().MustGet(ctx, "catalog.metadata.imdb.baseUrl", "https://www.imdb.com").String(), "/")
	response, err := g.Client().SetTimeout(8*time.Second).
		SetHeader("Accept", "text/html,application/xhtml+xml").
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		SetHeader("User-Agent", browserUserAgent(ctx)).
		Get(ctx, fmt.Sprintf("%s/title/%s/", baseURL, url.PathEscape(externalId)))
	if err != nil {
		return nil, err
	}
	defer response.Close()
	if response.StatusCode != 200 {
		return nil, gerror.Newf("imdb request failed with status %d", response.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, imdbMaxBodySize+1))
	if err != nil {
		return nil, err
	}
	if len(body) > imdbMaxBodySize {
		return nil, gerror.New("imdb response is too large")
	}
	return body, nil
}

func parseIMDbGraphQL(externalId string, body []byte) (*Item, error) {
	var response imdbGraphQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	if len(response.Errors) > 0 {
		return nil, gerror.New(response.Errors[0].Message)
	}
	if response.Data.Title == nil || response.Data.Title.Id == "" {
		return nil, gerror.New("imdb graphql metadata is empty")
	}
	data := response.Data.Title
	mediaType, ok := imdbMediaType(data.TitleType.Id)
	if !ok {
		return nil, gerror.New("unsupported imdb media type")
	}
	title := strings.TrimSpace(data.TitleText.Text)
	if title == "" {
		return nil, gerror.New("imdb title is empty")
	}
	originalTitle := strings.TrimSpace(data.OriginalTitleText.Text)
	if originalTitle == title {
		originalTitle = ""
	}
	genres := make([]string, 0, len(data.TitleGenres.Genres))
	for _, item := range data.TitleGenres.Genres {
		if genre := strings.TrimSpace(item.Genre.Text); genre != "" {
			genres = append(genres, genre)
		}
	}
	releaseDate := imdbReleaseDate(data.ReleaseDate)
	yearValue := data.ReleaseYear.Year
	if yearValue == 0 {
		yearValue = data.ReleaseDate.Year
	}
	return &Item{
		Provider:      ProviderIMDb,
		ProviderId:    externalId,
		TmdbType:      mediaType,
		Title:         title,
		OriginalTitle: originalTitle,
		Year:          imdbYear(yearValue),
		ReleaseDate:   releaseDate,
		Overview:      strings.TrimSpace(data.Plot.PlotText.PlainText),
		PosterUrl:     strings.TrimSpace(data.PrimaryImage.Url),
		Rating:        data.RatingsSummary.AggregateRating,
		Genres:        genres,
		ImdbId:        externalId,
	}, nil
}

func imdbReleaseDate(value imdbGraphQLReleaseDate) string {
	if value.Year <= 0 {
		return ""
	}
	if value.Month <= 0 {
		return fmt.Sprintf("%04d", value.Year)
	}
	if value.Day <= 0 {
		return fmt.Sprintf("%04d-%02d", value.Year, value.Month)
	}
	return fmt.Sprintf("%04d-%02d-%02d", value.Year, value.Month, value.Day)
}

func imdbYear(value int) string {
	if value <= 0 {
		return ""
	}
	return fmt.Sprintf("%04d", value)
}

func parseIMDb(externalId string, body []byte) (Item, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return Item{}, err
	}
	for _, raw := range metadataJsonLdScripts(doc) {
		var value any
		decoder := json.NewDecoder(strings.NewReader(raw))
		decoder.UseNumber()
		if decoder.Decode(&value) != nil {
			continue
		}
		if item, ok := imdbFromJsonLd(externalId, value); ok {
			return item, nil
		}
	}
	return Item{}, gerror.New("imdb json-ld metadata is missing")
}

func imdbFromJsonLd(externalId string, value any) (Item, bool) {
	switch typed := value.(type) {
	case []any:
		for _, item := range typed {
			if metadata, ok := imdbFromJsonLd(externalId, item); ok {
				return metadata, true
			}
		}
	case map[string]any:
		if graph, ok := typed["@graph"]; ok {
			if metadata, found := imdbFromJsonLd(externalId, graph); found {
				return metadata, true
			}
		}
		mediaType, ok := imdbMediaType(typed["@type"])
		if !ok {
			return Item{}, false
		}
		title := metadataString(typed["name"])
		if title == "" {
			return Item{}, false
		}
		originalTitle := metadataString(typed["alternateName"])
		if originalTitle == title {
			originalTitle = ""
		}
		aggregateRating, _ := typed["aggregateRating"].(map[string]any)
		releaseDate := metadataString(typed["datePublished"])
		return Item{
			Provider:      ProviderIMDb,
			ProviderId:    externalId,
			TmdbType:      mediaType,
			Title:         title,
			OriginalTitle: originalTitle,
			Year:          year(releaseDate),
			ReleaseDate:   releaseDate,
			Overview:      metadataString(typed["description"]),
			PosterUrl:     metadataImageUrl(typed["image"]),
			Rating:        metadataFloat64(aggregateRating["ratingValue"]),
			Genres:        metadataStrings(typed["genre"]),
			ImdbId:        externalId,
		}, true
	}
	return Item{}, false
}

func imdbMediaType(value any) (string, bool) {
	for _, typeName := range metadataStrings(value) {
		switch strings.ToLower(typeName) {
		case "movie", "tvmovie", "video", "videoobject", "short":
			return consts.CatalogMetadataTmdbTypeMovie, true
		case "tvseries", "tvminiseries", "tvepisode", "tvspecial":
			return consts.CatalogMetadataTmdbTypeTv, true
		}
	}
	return "", false
}

var _ Provider = (*IMDbProvider)(nil)
