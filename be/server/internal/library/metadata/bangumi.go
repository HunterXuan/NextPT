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

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	bangumiMaxBodySize = 2 * 1024 * 1024
	bangumiMaxTags     = 8
)

var bangumiImdbIdPattern = regexp.MustCompile(`(?i)tt\d{5,}`)

type BangumiProvider struct{}

type bangumiSubjectResponse struct {
	Id      int64                 `json:"id"`
	Name    string                `json:"name"`
	NameCn  string                `json:"name_cn"`
	Summary string                `json:"summary"`
	Date    string                `json:"date"`
	Images  *bangumiSubjectImages `json:"images"`
	Rating  bangumiSubjectRating  `json:"rating"`
	Tags    []bangumiSubjectTag   `json:"tags"`
	Infobox []bangumiInfoboxItem  `json:"infobox"`
}

type bangumiSubjectImages struct {
	Large  string `json:"large"`
	Common string `json:"common"`
	Medium string `json:"medium"`
}

type bangumiSubjectRating struct {
	Score float64 `json:"score"`
}

type bangumiSubjectTag struct {
	Name string `json:"name"`
}

type bangumiInfoboxItem struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (p *BangumiProvider) Name() string {
	return ProviderBangumi
}

func (p *BangumiProvider) Get(ctx context.Context, in GetInput) (*Item, error) {
	if !numericIdPattern.MatchString(in.ExternalId) {
		return nil, gerror.New("invalid bangumi id")
	}
	body, err := p.fetchSubject(ctx, in.ExternalId)
	if err != nil {
		return nil, err
	}
	metadata, err := parseBangumi(in.ExternalId, body)
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

func (p *BangumiProvider) fetchSubject(ctx context.Context, externalId string) ([]byte, error) {
	baseURL := strings.TrimRight(g.Cfg().MustGet(ctx, "catalog.metadata.bangumi.baseUrl", "https://api.bgm.tv/v0").String(), "/")
	response, err := g.Client().SetTimeout(8*time.Second).
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", browserUserAgent(ctx)).
		Get(ctx, fmt.Sprintf("%s/subjects/%s", baseURL, url.PathEscape(externalId)))
	if err != nil {
		return nil, err
	}
	defer response.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, gerror.Newf("bangumi request failed with status %d", response.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, bangumiMaxBodySize+1))
	if err != nil {
		return nil, err
	}
	if len(body) > bangumiMaxBodySize {
		return nil, gerror.New("bangumi response is too large")
	}
	return body, nil
}

func parseBangumi(externalId string, body []byte) (Item, error) {
	var subject bangumiSubjectResponse
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(&subject); err != nil {
		return Item{}, err
	}
	title := strings.TrimSpace(subject.NameCn)
	originalTitle := strings.TrimSpace(subject.Name)
	if title == "" {
		title = originalTitle
	}
	if title == "" || subject.Id == 0 {
		return Item{}, gerror.New("bangumi metadata is empty")
	}
	if originalTitle == title {
		originalTitle = ""
	}
	return Item{
		Provider:      ProviderBangumi,
		ProviderId:    externalId,
		Title:         title,
		OriginalTitle: originalTitle,
		Year:          year(subject.Date),
		ReleaseDate:   strings.TrimSpace(subject.Date),
		Overview:      strings.TrimSpace(subject.Summary),
		PosterUrl:     bangumiPosterUrl(subject.Images),
		Rating:        subject.Rating.Score,
		Genres:        bangumiTagNames(subject.Tags),
		ImdbId:        bangumiImdbId(subject.Infobox),
	}, nil
}

func bangumiPosterUrl(images *bangumiSubjectImages) string {
	if images == nil {
		return ""
	}
	for _, value := range []string{images.Large, images.Common, images.Medium} {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if strings.HasPrefix(value, "//") {
			return "https:" + value
		}
		return value
	}
	return ""
}

func bangumiTagNames(tags []bangumiSubjectTag) []string {
	result := make([]string, 0, min(len(tags), bangumiMaxTags))
	seen := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		name := strings.TrimSpace(tag.Name)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, name)
		if len(result) == bangumiMaxTags {
			break
		}
	}
	return result
}

func bangumiImdbId(infobox []bangumiInfoboxItem) string {
	for _, item := range infobox {
		if !strings.Contains(strings.ToLower(strings.TrimSpace(item.Key)), "imdb") {
			continue
		}
		for _, value := range bangumiInfoboxValues(item.Value) {
			if imdbId := bangumiImdbIdPattern.FindString(value); imdbId != "" {
				return strings.ToLower(imdbId)
			}
		}
	}
	return ""
}

func bangumiInfoboxValues(value any) []string {
	switch typed := value.(type) {
	case string:
		if value := strings.TrimSpace(typed); value != "" {
			return []string{value}
		}
	case json.Number:
		return []string{typed.String()}
	case []any:
		var result []string
		for _, item := range typed {
			result = append(result, bangumiInfoboxValues(item)...)
		}
		return result
	case map[string]any:
		var result []string
		for _, item := range typed {
			result = append(result, bangumiInfoboxValues(item)...)
		}
		return result
	}
	return nil
}

var _ Provider = (*BangumiProvider)(nil)
