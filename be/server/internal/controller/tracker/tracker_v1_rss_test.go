package tracker

import (
	"encoding/xml"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestTrackerRssAbsoluteURL(t *testing.T) {
	got := trackerRssAbsoluteURL(
		"https://nextpt.example",
		"/api/tracker/download",
		url.Values{"id": []string{"42"}, "passkey": []string{"abc123"}},
	)
	want := "https://nextpt.example/api/tracker/download?id=42&passkey=abc123"
	if got != want {
		t.Fatalf("unexpected RSS URL: got %q, want %q", got, want)
	}
}

func TestTrackerRssDetailPath(t *testing.T) {
	if got, want := trackerRssDetailPath(42), "/catalog/torrents/42"; got != want {
		t.Fatalf("unexpected detail path: got %q, want %q", got, want)
	}
}

func TestTrackerRssFormatBytes(t *testing.T) {
	tests := []struct {
		size uint64
		want string
	}{
		{size: 512, want: "512 B"},
		{size: 1024, want: "1.00 KiB"},
		{size: 5 * 1024 * 1024 * 1024, want: "5.00 GiB"},
	}
	for _, tt := range tests {
		if got := trackerRssFormatBytes(tt.size); got != tt.want {
			t.Fatalf("unexpected formatted size for %d: got %q, want %q", tt.size, got, tt.want)
		}
	}
}

func TestTrackerRssXMLContract(t *testing.T) {
	document := trackerRssDocument{
		Version: "2.0",
		Channel: trackerRssChannel{
			Title:         "NextPT & RSS",
			Link:          "https://nextpt.example/catalog/torrents",
			Description:   "Latest torrents",
			Language:      "en-US",
			LastBuildDate: time.Date(2026, 7, 10, 6, 0, 0, 0, time.UTC).Format(time.RFC1123Z),
			Items: []trackerRssItem{{
				Title:       "A & B",
				Link:        "https://nextpt.example/catalog/torrents/42",
				Guid:        trackerRssGuid{IsPermaLink: true, Value: "https://nextpt.example/catalog/torrents/42"},
				Description: "5.00 GiB / 3 seeders",
				Category:    "Movies",
				PubDate:     time.Date(2026, 7, 10, 5, 0, 0, 0, time.UTC).Format(time.RFC1123Z),
				Enclosure: trackerRssEnclosure{
					URL:    "https://nextpt.example/api/tracker/download?id=42&passkey=abc123",
					Length: 0,
					Type:   "application/x-bittorrent",
				},
			}},
		},
	}

	payload, err := xml.Marshal(document)
	if err != nil {
		t.Fatalf("xml.Marshal() error = %v", err)
	}
	if !strings.Contains(string(payload), "NextPT &amp; RSS") || !strings.Contains(string(payload), "passkey=abc123") {
		t.Fatalf("unexpected RSS XML: %s", payload)
	}

	var decoded trackerRssDocument
	if err := xml.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("xml.Unmarshal() error = %v", err)
	}
	if decoded.Version != "2.0" || len(decoded.Channel.Items) != 1 {
		t.Fatalf("unexpected decoded RSS document: %+v", decoded)
	}
	item := decoded.Channel.Items[0]
	if item.Title != "A & B" || item.Enclosure.Type != "application/x-bittorrent" {
		t.Fatalf("unexpected decoded RSS item: %+v", item)
	}
}
