package site

import (
	"reflect"
	"testing"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func TestEncodeConfigValueWrapsValAndPreservesTypedValue(t *testing.T) {
	s := NewSiteConfigDomain()
	tests := []struct {
		name  string
		value any
		want  any
	}{
		{
			name:  "numeric string stays string",
			value: "1800",
			want:  "1800",
		},
		{
			name:  "int",
			value: int64(1800),
			want:  float64(1800),
		},
		{
			name:  "bool",
			value: true,
			want:  true,
		},
		{
			name:  "object",
			value: map[string]any{"enabled": true},
			want: map[string]any{
				"enabled": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, err := gjson.LoadJson([]byte(s.encodeConfigValue(tt.value)))
			if err != nil {
				t.Fatalf("LoadJson error = %v", err)
			}

			got := j.Get(siteConfigValueField).Val()
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("encoded val = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestInferConfigValueTypeFromDefaults(t *testing.T) {
	s := NewSiteConfigDomain()
	tests := []struct {
		path string
		want string
	}{
		{path: consts.SiteConfigTrackerAnnounceInterval, want: "int"},
		{path: consts.SiteConfigTrackerBonusBase, want: "float"},
		{path: consts.SiteConfigIamRegisterEnabled, want: "boolean"},
		{path: consts.SiteConfigCatalogTorrentSource, want: "string"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			group, key := splitSiteConfigPathForTest(tt.path)
			got := string(s.getConfigValueType(group, key))
			if got != tt.want {
				t.Fatalf("value type = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNormalizeConfigValueUsesDefaultType(t *testing.T) {
	s := NewSiteConfigDomain()
	tests := []struct {
		name    string
		path    string
		value   any
		want    any
		wantErr bool
	}{
		{
			name:  "string default keeps numeric-looking string",
			path:  consts.SiteConfigCatalogTorrentSource,
			value: "123",
			want:  "123",
		},
		{
			name:  "int default accepts integer number",
			path:  consts.SiteConfigTrackerAnnounceInterval,
			value: float64(1800),
			want:  int64(1800),
		},
		{
			name:    "int default rejects decimal number",
			path:    consts.SiteConfigTrackerAnnounceInterval,
			value:   float64(1800.5),
			wantErr: true,
		},
		{
			name:    "int default rejects non-number string",
			path:    consts.SiteConfigTrackerAnnounceInterval,
			value:   "abc",
			wantErr: true,
		},
		{
			name:    "int default rejects overflow string",
			path:    consts.SiteConfigTrackerAnnounceInterval,
			value:   "9223372036854775808",
			wantErr: true,
		},
		{
			name:  "float default accepts decimal number",
			path:  consts.SiteConfigTrackerBonusBase,
			value: "0.4",
			want:  float64(0.4),
		},
		{
			name:  "bool default accepts boolean",
			path:  consts.SiteConfigIamRegisterEnabled,
			value: false,
			want:  false,
		},
		{
			name:  "bool default accepts trimmed string",
			path:  consts.SiteConfigIamRegisterEnabled,
			value: " false ",
			want:  false,
		},
		{
			name:    "bool default rejects arbitrary string",
			path:    consts.SiteConfigIamRegisterEnabled,
			value:   "abc",
			wantErr: true,
		},
		{
			name:  "invite registration email pattern accepts valid expression",
			path:  consts.SiteConfigIamInviteBypassEmailPattern,
			value: " ^[A-Z0-9._%+-]+@example\\.com$ ",
			want:  "^[A-Z0-9._%+-]+@example\\.com$",
		},
		{
			name:  "invite registration email pattern accepts empty expression",
			path:  consts.SiteConfigIamInviteBypassEmailPattern,
			value: "  ",
			want:  "",
		},
		{
			name:    "invite registration email pattern rejects invalid expression",
			path:    consts.SiteConfigIamInviteBypassEmailPattern,
			value:   "[invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group, key := splitSiteConfigPathForTest(tt.path)
			got, err := s.normalizeConfigValue(group, key, tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatal("normalizeConfigValue error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeConfigValue error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("normalized value = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNormalizeSiteAdvertisements(t *testing.T) {
	s := NewSiteConfigDomain()
	group, key := splitSiteConfigPathForTest(consts.SiteConfigSiteAdvertisements)
	value, err := s.normalizeConfigValue(group, key, map[string]any{
		"home": map[string]any{
			"enabled": true,
			"title":   " Partner ",
			"image":   "https://cdn.example.com/banner.webp",
			"url":     "/catalog/torrents",
		},
	})
	if err != nil {
		t.Fatalf("normalizeConfigValue() error = %v", err)
	}
	advertisements, ok := value.(model.SiteAdvertisements)
	if !ok {
		t.Fatalf("normalized value type = %T, want model.SiteAdvertisements", value)
	}
	if advertisements[consts.SiteAdvertisementPlacementHome].Title != "Partner" {
		t.Fatalf("home title = %q, want Partner", advertisements[consts.SiteAdvertisementPlacementHome].Title)
	}
	if advertisements[consts.SiteAdvertisementPlacementHome].AspectRatio != "8:1" {
		t.Fatalf("home aspect ratio = %q, want 8:1", advertisements[consts.SiteAdvertisementPlacementHome].AspectRatio)
	}
	if _, ok := advertisements[consts.SiteAdvertisementPlacementForumList]; !ok {
		t.Fatal("normalized advertisements should include all known placements")
	}
}

func TestDecodeConfigValueUnwrapsStoredValue(t *testing.T) {
	s := NewSiteConfigDomain()
	tests := []struct {
		name string
		cfg  *entity.SiteConfig
		def  any
		want any
	}{
		{
			name: "string",
			cfg: &entity.SiteConfig{
				Value: gjson.New(map[string]any{siteConfigValueField: "NextPT"}),
			},
			want: "NextPT",
		},
		{
			name: "object",
			cfg: &entity.SiteConfig{
				Value: gjson.New(map[string]any{
					siteConfigValueField: map[string]any{"enabled": true},
				}),
			},
			want: map[string]any{"enabled": true},
		},
		{
			name: "missing value uses default",
			cfg:  &entity.SiteConfig{Value: gjson.New(map[string]any{})},
			def:  "fallback",
			want: "fallback",
		},
		{
			name: "nil config uses default",
			def:  "fallback",
			want: "fallback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.getConfigValue(tt.cfg, tt.def)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("config value = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func splitSiteConfigPathForTest(path string) (string, string) {
	for i, r := range path {
		if r == '.' {
			return path[:i], path[i+1:]
		}
	}
	return path, ""
}
