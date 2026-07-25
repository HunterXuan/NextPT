package metadata

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/net/html"
)

const (
	ProviderTMDB    = "tmdb"
	ProviderIMDb    = "imdb"
	ProviderDouban  = "douban"
	ProviderBangumi = "bangumi"

	defaultBrowserUserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"
)

var (
	numericIdPattern = regexp.MustCompile(`^\d+$`)
	imdbIdPattern    = regexp.MustCompile(`^tt\d+$`)
)

type GetInput struct {
	ExternalId string
	MediaType  string
	Locale     string
}

type Item struct {
	Provider      string   `json:"provider"`
	ProviderId    string   `json:"providerId"`
	TmdbType      string   `json:"tmdbType"`
	Title         string   `json:"title"`
	OriginalTitle string   `json:"originalTitle"`
	Year          string   `json:"year"`
	ReleaseDate   string   `json:"releaseDate"`
	Overview      string   `json:"overview"`
	PosterUrl     string   `json:"posterUrl"`
	BackdropUrl   string   `json:"backdropUrl"`
	Rating        float64  `json:"rating"`
	Genres        []string `json:"genres"`
	ImdbId        string   `json:"imdbId"`
}

type SearchResult struct {
	List         []Item
	Page         int
	TotalPages   int
	TotalResults int
}

type Provider interface {
	Name() string
	Get(ctx context.Context, in GetInput) (*Item, error)
}

func year(value string) string {
	if len(value) < 4 {
		return ""
	}
	return value[:4]
}

func browserUserAgent(ctx context.Context) string {
	return strings.TrimSpace(g.Cfg().MustGet(ctx, "catalog.metadata.userAgent", defaultBrowserUserAgent).String())
}

func metadataJsonLdScripts(root *html.Node) []string {
	var scripts []string
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "script" && strings.EqualFold(htmlAttr(node, "type"), "application/ld+json") {
			var content strings.Builder
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				if child.Type == html.TextNode {
					content.WriteString(child.Data)
				}
			}
			if value := strings.TrimSpace(content.String()); value != "" {
				scripts = append(scripts, value)
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(root)
	return scripts
}

func htmlAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if strings.EqualFold(attr.Key, key) {
			return strings.TrimSpace(attr.Val)
		}
	}
	return ""
}

func metadataString(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case json.Number:
		return typed.String()
	case []any:
		if len(typed) > 0 {
			return metadataString(typed[0])
		}
	}
	return ""
}

func metadataStrings(value any) []string {
	switch typed := value.(type) {
	case string:
		if value := strings.TrimSpace(typed); value != "" {
			return []string{value}
		}
	case []any:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			if value := metadataString(item); value != "" {
				result = append(result, value)
			}
		}
		return result
	}
	return nil
}

func metadataFloat64(value any) float64 {
	result, _ := strconv.ParseFloat(metadataString(value), 64)
	return result
}

func metadataImageUrl(value any) string {
	if image := metadataString(value); image != "" {
		return image
	}
	if image, ok := value.(map[string]any); ok {
		return metadataString(image["url"])
	}
	return ""
}

func htmlHasClass(node *html.Node, className string) bool {
	for _, value := range strings.Fields(htmlAttr(node, "class")) {
		if value == className {
			return true
		}
	}
	return false
}

func htmlFirstByClass(root *html.Node, className string) *html.Node {
	if root.Type == html.ElementNode && htmlHasClass(root, className) {
		return root
	}
	for child := root.FirstChild; child != nil; child = child.NextSibling {
		if node := htmlFirstByClass(child, className); node != nil {
			return node
		}
	}
	return nil
}

func htmlText(node *html.Node) string {
	if node == nil {
		return ""
	}
	var content strings.Builder
	var walk func(*html.Node)
	walk = func(current *html.Node) {
		if current.Type == html.TextNode {
			content.WriteString(current.Data)
			content.WriteByte(' ')
		}
		for child := current.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(node)
	return strings.Join(strings.Fields(content.String()), " ")
}

func htmlMetaContent(root *html.Node, key string, value string) string {
	if root.Type == html.ElementNode && root.Data == "meta" && strings.EqualFold(htmlAttr(root, key), value) {
		return htmlAttr(root, "content")
	}
	for child := root.FirstChild; child != nil; child = child.NextSibling {
		if content := htmlMetaContent(child, key, value); content != "" {
			return content
		}
	}
	return ""
}
