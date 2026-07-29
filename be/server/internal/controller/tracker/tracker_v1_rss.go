package tracker

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
	"time"

	v1 "server/api/tracker/v1"
	"server/internal/library/contexts"
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
)

type trackerRssDocument struct {
	XMLName xml.Name          `xml:"rss"`
	Version string            `xml:"version,attr"`
	Channel trackerRssChannel `xml:"channel"`
}

type trackerRssChannel struct {
	Title         string           `xml:"title"`
	Link          string           `xml:"link"`
	Description   string           `xml:"description"`
	Language      string           `xml:"language,omitempty"`
	LastBuildDate string           `xml:"lastBuildDate,omitempty"`
	Items         []trackerRssItem `xml:"item"`
}

type trackerRssItem struct {
	Title       string              `xml:"title"`
	Link        string              `xml:"link"`
	Guid        trackerRssGuid      `xml:"guid"`
	Description string              `xml:"description"`
	Category    string              `xml:"category,omitempty"`
	PubDate     string              `xml:"pubDate,omitempty"`
	Enclosure   trackerRssEnclosure `xml:"enclosure"`
}

type trackerRssGuid struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

type trackerRssEnclosure struct {
	URL    string `xml:"url,attr"`
	Length int    `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

func (c *ControllerV1) Rss(ctx context.Context, req *v1.RssReq) (res *v1.RssRes, err error) {
	out, err := service.CatalogTorrentUsecase().ListRss(ctx, contexts.GetActor(ctx), catalogin.TorrentRssInp{
		Size:            req.Size,
		Keyword:         req.Keyword,
		CategoryIds:     req.CategoryIds,
		TagIds:          req.TagIds,
		Promotion:       req.Promotion,
		PromotionOnly:   req.PromotionOnly,
		SeedStatus:      req.SeedStatus,
		FeaturedOnly:    req.FeaturedOnly,
		MinSize:         req.MinSize,
		MaxSize:         req.MaxSize,
		PublishedWithin: req.PublishedWithin,
		TorrentMetadataFilterInp: catalogin.TorrentMetadataFilterInp{
			ImdbId:    req.ImdbId,
			DoubanId:  req.DoubanId,
			BangumiId: req.BangumiId,
			TmdbId:    req.TmdbId,
			TmdbType:  req.TmdbType,
		},
	})
	if err != nil {
		return nil, err
	}

	r := g.RequestFromCtx(ctx)
	origin := trackerRssRequestOrigin(r)
	passkey := r.GetQuery("passkey").String()
	lastBuildDate := ""
	if len(out.List) > 0 && out.List[0].CreatedAt != nil {
		lastBuildDate = out.List[0].CreatedAt.Time.Format(time.RFC1123Z)
	}
	document := trackerRssDocument{
		Version: "2.0",
		Channel: trackerRssChannel{
			Title:         out.Title,
			Link:          trackerRssAbsoluteURL(origin, "/catalog/torrents", nil),
			Description:   out.Description,
			Language:      out.Language,
			LastBuildDate: lastBuildDate,
			Items:         trackerRssItems(ctx, origin, passkey, out.List),
		},
	}
	payload, err := xml.MarshalIndent(document, "", "  ")
	if err != nil {
		return nil, err
	}

	r.Response.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	r.Response.Header().Set("Cache-Control", "private, no-store")
	r.Response.Write(xml.Header)
	r.Response.Write(payload)
	r.ExitAll()
	return nil, nil
}

func trackerRssItems(ctx context.Context, origin, passkey string, items []catalogout.TorrentRssItem) []trackerRssItem {
	result := make([]trackerRssItem, 0, len(items))
	for _, item := range items {
		detailURL := trackerRssAbsoluteURL(origin, trackerRssDetailPath(item.Id), nil)
		downloadURL := trackerRssAbsoluteURL(origin, "/api/tracker/download", url.Values{
			"id":      []string{fmt.Sprint(item.Id)},
			"passkey": []string{passkey},
		})
		description := fmt.Sprintf(
			gi18n.T(ctx, "catalog.torrent.rss.item_description"),
			trackerRssFormatBytes(item.Size),
			item.Seeders,
			item.Leechers,
			item.Snatched,
		)
		if strings.TrimSpace(item.SubTitle) != "" {
			description = item.SubTitle + "\n" + description
		}
		pubDate := ""
		if item.CreatedAt != nil {
			pubDate = item.CreatedAt.Time.Format(time.RFC1123Z)
		}
		result = append(result, trackerRssItem{
			Title:       item.Name,
			Link:        detailURL,
			Guid:        trackerRssGuid{IsPermaLink: true, Value: detailURL},
			Description: description,
			Category:    item.Category,
			PubDate:     pubDate,
			Enclosure: trackerRssEnclosure{
				URL:    downloadURL,
				Length: 0,
				Type:   "application/x-bittorrent",
			},
		})
	}
	return result
}

func trackerRssRequestOrigin(r *ghttp.Request) string {
	scheme := trackerRssForwardedValue(r.Header.Get("X-Forwarded-Proto"))
	if scheme == "" {
		scheme = r.GetSchema()
	}
	if scheme == "" {
		scheme = "http"
	}
	host := trackerRssForwardedValue(r.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = r.Host
	}
	return scheme + "://" + host
}

func trackerRssForwardedValue(value string) string {
	if index := strings.IndexByte(value, ','); index >= 0 {
		value = value[:index]
	}
	return strings.TrimSpace(value)
}

func trackerRssDetailPath(torrentId uint64) string {
	return fmt.Sprintf("/catalog/torrents/%d", torrentId)
}

func trackerRssAbsoluteURL(origin, requestPath string, query url.Values) string {
	u, err := url.Parse(origin)
	if err != nil {
		return origin + requestPath
	}
	u.Path = requestPath
	u.RawQuery = query.Encode()
	return u.String()
}

func trackerRssFormatBytes(size uint64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := uint64(unit), 0
	for n := size / unit; n >= unit && exp < 5; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}
