package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"

	"server/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/net/html"
)

const doubanMaxBodySize = 2 * 1024 * 1024

var (
	doubanImdbUrlPattern   = regexp.MustCompile(`(?i)imdb\.com/title/(tt\d+)`)
	doubanImdbLabelPattern = regexp.MustCompile(`(?is)IMDb.{0,120}?(tt\d{5,})`)
	doubanOriginalPattern  = regexp.MustCompile(`^(.+?)[（(](\d{4})[）)]$`)
)

type DoubanProvider struct{}

type doubanSubjectResponse struct {
	Id            string   `json:"id"`
	Type          string   `json:"type"`
	Subtype       string   `json:"subtype"`
	Title         string   `json:"title"`
	OriginalTitle string   `json:"original_title"`
	Year          string   `json:"year"`
	Intro         string   `json:"intro"`
	CoverUrl      string   `json:"cover_url"`
	Genres        []string `json:"genres"`
	IsTv          bool     `json:"is_tv"`
	Rating        struct {
		Value float64 `json:"value"`
	} `json:"rating"`
}

func (p *DoubanProvider) Name() string {
	return ProviderDouban
}

func (p *DoubanProvider) Get(ctx context.Context, in GetInput) (*Item, error) {
	if !numericIdPattern.MatchString(in.ExternalId) {
		return nil, gerror.New("invalid douban id")
	}
	jsonBody, jsonErr := p.fetchSubjectJSON(ctx, in.ExternalId)
	if jsonErr == nil {
		if item, err := parseDoubanJSON(in.ExternalId, jsonBody); err == nil {
			return item, nil
		} else {
			jsonErr = err
		}
	}
	body, err := p.fetchPage(ctx, in.ExternalId)
	if err != nil {
		return nil, gerror.Newf("douban json request failed: %v; page fallback failed: %v", jsonErr, err)
	}
	metadata, err := parseDouban(in.ExternalId, body)
	if err != nil {
		return nil, gerror.Newf("douban json request failed: %v; page fallback failed: %v", jsonErr, err)
	}
	return &metadata, nil
}

func (p *DoubanProvider) fetchSubjectJSON(ctx context.Context, externalId string) ([]byte, error) {
	apiBaseURL := strings.TrimRight(g.Cfg().MustGet(ctx, "catalog.metadata.douban.apiBaseUrl", "https://m.douban.com/rexxar/api/v2").String(), "/")
	pageBaseURL := strings.TrimRight(g.Cfg().MustGet(ctx, "catalog.metadata.douban.baseUrl", "https://m.douban.com/movie").String(), "/")
	response, err := g.Client().SetTimeout(8*time.Second).
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9").
		SetHeader("Referer", fmt.Sprintf("%s/subject/%s/", pageBaseURL, url.PathEscape(externalId))).
		SetHeader("User-Agent", browserUserAgent(ctx)).
		Get(ctx, fmt.Sprintf("%s/movie/%s", apiBaseURL, url.PathEscape(externalId)))
	if err != nil {
		return nil, err
	}
	defer response.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, gerror.Newf("douban json request failed with status %d", response.StatusCode)
	}
	return readDoubanBody(response.Body, "douban json")
}

func (p *DoubanProvider) fetchPage(ctx context.Context, externalId string) ([]byte, error) {
	baseURL := strings.TrimRight(g.Cfg().MustGet(ctx, "catalog.metadata.douban.baseUrl", "https://m.douban.com/movie").String(), "/")
	response, err := g.Client().SetTimeout(8*time.Second).
		SetHeader("Accept", "text/html,application/xhtml+xml").
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9").
		SetHeader("User-Agent", browserUserAgent(ctx)).
		Get(ctx, fmt.Sprintf("%s/subject/%s/", baseURL, url.PathEscape(externalId)))
	if err != nil {
		return nil, err
	}
	defer response.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, gerror.Newf("douban request failed with status %d", response.StatusCode)
	}
	return readDoubanBody(response.Body, "douban page")
}

func readDoubanBody(reader io.Reader, source string) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(reader, doubanMaxBodySize+1))
	if err != nil {
		return nil, err
	}
	if len(body) > doubanMaxBodySize {
		return nil, gerror.Newf("%s response is too large", source)
	}
	return body, nil
}

func parseDoubanJSON(externalId string, body []byte) (*Item, error) {
	var subject doubanSubjectResponse
	if err := json.Unmarshal(body, &subject); err != nil {
		return nil, err
	}
	if subject.Id == "" || strings.TrimSpace(subject.Title) == "" {
		return nil, gerror.New("douban json metadata is empty")
	}
	mediaType := consts.CatalogMetadataTmdbTypeMovie
	if subject.IsTv || strings.Contains(strings.ToLower(subject.Type), "tv") || strings.Contains(strings.ToLower(subject.Subtype), "tv") {
		mediaType = consts.CatalogMetadataTmdbTypeTv
	}
	return &Item{
		Provider:      ProviderDouban,
		ProviderId:    externalId,
		TmdbType:      mediaType,
		Title:         strings.TrimSpace(subject.Title),
		OriginalTitle: strings.TrimSpace(subject.OriginalTitle),
		Year:          strings.TrimSpace(subject.Year),
		Overview:      strings.TrimSpace(subject.Intro),
		PosterUrl:     strings.TrimSpace(subject.CoverUrl),
		Rating:        subject.Rating.Value,
		Genres:        append([]string(nil), subject.Genres...),
	}, nil
}

func parseDouban(externalId string, body []byte) (Item, error) {
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
		if item, ok := doubanFromJsonLd(externalId, value); ok {
			if item.ImdbId == "" {
				item.ImdbId = doubanImdbId(string(body))
			}
			return item, nil
		}
	}
	if item, ok := doubanFromMobilePage(externalId, doc, string(body)); ok {
		return item, nil
	}
	return Item{}, gerror.New("douban metadata is missing")
}

func doubanFromMobilePage(externalId string, doc *html.Node, body string) (Item, bool) {
	title := htmlText(htmlFirstByClass(doc, "sub-title"))
	if title == "" {
		return Item{}, false
	}
	originalTitle := htmlText(htmlFirstByClass(doc, "sub-original-title"))
	yearValue := ""
	if match := doubanOriginalPattern.FindStringSubmatch(originalTitle); len(match) > 2 {
		originalTitle = strings.TrimSpace(match[1])
		yearValue = match[2]
	}
	mediaType := consts.CatalogMetadataTmdbTypeMovie
	if strings.Contains(body, "mark-tv") || strings.Contains(body, "mark-tvseries") {
		mediaType = consts.CatalogMetadataTmdbTypeTv
	}
	overview := ""
	if intro := htmlFirstByClass(doc, "subject-intro"); intro != nil {
		overview = htmlText(htmlFirstByClass(intro, "bd"))
	}
	return Item{
		Provider:      ProviderDouban,
		ProviderId:    externalId,
		TmdbType:      mediaType,
		Title:         title,
		OriginalTitle: originalTitle,
		Year:          yearValue,
		Overview:      overview,
		PosterUrl:     htmlMetaContent(doc, "itemprop", "image"),
		Rating:        metadataFloat64(htmlMetaContent(doc, "itemprop", "ratingValue")),
		ImdbId:        doubanImdbId(body),
	}, true
}

func doubanFromJsonLd(externalId string, value any) (Item, bool) {
	switch typed := value.(type) {
	case []any:
		for _, item := range typed {
			if metadata, ok := doubanFromJsonLd(externalId, item); ok {
				return metadata, true
			}
		}
	case map[string]any:
		if graph, ok := typed["@graph"]; ok {
			if metadata, found := doubanFromJsonLd(externalId, graph); found {
				return metadata, true
			}
		}
		typeName := strings.ToLower(metadataString(typed["@type"]))
		if !strings.Contains(typeName, "movie") && !strings.Contains(typeName, "tv") && !strings.Contains(typeName, "series") {
			return Item{}, false
		}
		name := metadataString(typed["name"])
		if name == "" {
			return Item{}, false
		}
		aggregateRating, _ := typed["aggregateRating"].(map[string]any)
		mediaType := consts.CatalogMetadataTmdbTypeMovie
		if strings.Contains(typeName, "tv") || strings.Contains(typeName, "series") {
			mediaType = consts.CatalogMetadataTmdbTypeTv
		}
		return Item{
			Provider:      ProviderDouban,
			ProviderId:    externalId,
			TmdbType:      mediaType,
			Title:         name,
			OriginalTitle: metadataString(typed["alternateName"]),
			Year:          year(metadataString(typed["datePublished"])),
			ReleaseDate:   metadataString(typed["datePublished"]),
			Overview:      metadataString(typed["description"]),
			PosterUrl:     metadataImageUrl(typed["image"]),
			Rating:        metadataFloat64(aggregateRating["ratingValue"]),
			Genres:        metadataStrings(typed["genre"]),
			ImdbId:        metadataImdbId(typed["sameAs"]),
		}, true
	}
	return Item{}, false
}

func metadataImdbId(value any) string {
	for _, item := range metadataStrings(value) {
		if match := doubanImdbUrlPattern.FindStringSubmatch(item); len(match) > 1 {
			return strings.ToLower(match[1])
		}
	}
	return ""
}

func doubanImdbId(body string) string {
	if match := doubanImdbUrlPattern.FindStringSubmatch(body); len(match) > 1 {
		return strings.ToLower(match[1])
	}
	if match := doubanImdbLabelPattern.FindStringSubmatch(body); len(match) > 1 {
		return strings.ToLower(match[1])
	}
	return ""
}

var _ Provider = (*DoubanProvider)(nil)
