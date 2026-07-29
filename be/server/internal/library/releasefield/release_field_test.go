package releasefield

import (
	"context"
	"reflect"
	"testing"

	"server/internal/model"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func TestResolverResolveGeneratedTitleAndNormalizeFields(t *testing.T) {
	resolver := newTestResolver(t, model.CatalogUploadConfig{
		Title: model.CatalogUploadTitleConfig{
			Mode: "generated",
			Parts: []model.CatalogUploadTitlePart{
				{Field: "year", Suffix: " "},
				{Field: "quality"},
			},
		},
		Fields: []model.CatalogUploadFieldConfig{
			{Key: "year", Type: "text", Required: true},
			{
				Key:      "quality",
				Type:     "select",
				Required: true,
				Options: &model.CatalogUploadFieldOptions{
					Source: "static",
					Items:  []model.CatalogUploadOptionItem{{Value: "1080p"}},
				},
			},
			{
				Key:  "audio",
				Type: "multiSelect",
				Options: &model.CatalogUploadFieldOptions{
					Source: "static",
					Items: []model.CatalogUploadOptionItem{
						{Value: "DTS"},
						{Value: "Atmos"},
					},
				},
			},
		},
	})

	result, err := resolver.Resolve(context.Background(), Input{
		ManualName:    "ignored",
		FallbackName:  "fallback",
		ReleaseFields: `{"year":2026,"quality":"1080p","audio":["DTS","Atmos","DTS"]}`,
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if result.Name != "2026 1080p" {
		t.Fatalf("Resolve() name = %q, want %q", result.Name, "2026 1080p")
	}
	if got, want := result.Fields["audio"], []string{"DTS", "Atmos"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Resolve() audio = %#v, want %#v", got, want)
	}
}

func TestResolverResolveUsesStoredFieldsAndManualOverride(t *testing.T) {
	resolver := newTestResolver(t, model.CatalogUploadConfig{
		Title: model.CatalogUploadTitleConfig{
			Mode:                "generated",
			AllowManualOverride: true,
			Parts:               []model.CatalogUploadTitlePart{{Field: "edition"}},
		},
		Fields: []model.CatalogUploadFieldConfig{{Key: "edition", Type: "text"}},
	})

	result, err := resolver.Resolve(context.Background(), Input{
		ManualName:   "Manual title",
		FallbackName: "fallback",
		StoredFields: gjson.New(map[string]any{"edition": "Director's Cut"}),
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if result.Name != "Manual title" {
		t.Fatalf("Resolve() name = %q, want manual title", result.Name)
	}
	if result.Fields["edition"] != "Director's Cut" {
		t.Fatalf("Resolve() edition = %#v", result.Fields["edition"])
	}
}

func TestResolverResolveRejectsMissingAndInvalidOptions(t *testing.T) {
	resolver := newTestResolver(t, model.CatalogUploadConfig{
		Fields: []model.CatalogUploadFieldConfig{
			{
				Key:      "quality",
				Type:     "select",
				Required: true,
				Options: &model.CatalogUploadFieldOptions{
					Source: "static",
					Items:  []model.CatalogUploadOptionItem{{Value: "1080p"}},
				},
			},
		},
	})

	if _, err := resolver.Resolve(context.Background(), Input{}); err == nil {
		t.Fatal("Resolve() expected required-field error")
	}
	if _, err := resolver.Resolve(context.Background(), Input{ReleaseFields: `{"quality":"invalid"}`}); err == nil {
		t.Fatal("Resolve() expected invalid-option error")
	}
}

func TestResolverResolveTagIdsReplacesControlledGroup(t *testing.T) {
	resolver := newTestResolver(t, model.CatalogUploadConfig{
		Fields: []model.CatalogUploadFieldConfig{
			{
				Key:      "resolution",
				Type:     "select",
				Required: true,
				Options: &model.CatalogUploadFieldOptions{
					Source: "tagGroup",
					Slug:   "resolution",
				},
			},
		},
	})
	resolver.SetTagCatalog(2, []TagGroup{
		{Id: 1, Slug: "resolution", CategoryIds: gjson.New([]uint{2})},
		{Id: 2, Slug: "source"},
	}, []Tag{
		{Id: 11, GroupId: 1, Value: "4k"},
		{Id: 12, GroupId: 1, Value: "1080p"},
		{Id: 21, GroupId: 2, Value: "remux"},
	})

	result, err := resolver.Resolve(context.Background(), Input{ReleaseFields: `{"resolution":"4k"}`})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	tagIds, err := resolver.ResolveTagIds(context.Background(), []uint{12, 21, 21}, result)
	if err != nil {
		t.Fatalf("ResolveTagIds() error = %v", err)
	}
	if want := []uint{21, 11}; !reflect.DeepEqual(tagIds, want) {
		t.Fatalf("ResolveTagIds() = %#v, want %#v", tagIds, want)
	}
}

func TestResolverResolveTagIdsRejectsWrongCategory(t *testing.T) {
	resolver := newTestResolver(t, model.CatalogUploadConfig{})
	resolver.SetTagCatalog(2, []TagGroup{
		{Id: 1, Slug: "resolution", CategoryIds: gjson.New([]uint{3})},
	}, []Tag{{Id: 11, GroupId: 1, Value: "4k"}})

	if _, err := resolver.ResolveTagIds(context.Background(), []uint{11}, Result{}); err == nil {
		t.Fatal("ResolveTagIds() expected category error")
	}
}

func newTestResolver(t *testing.T, config model.CatalogUploadConfig) *Resolver {
	t.Helper()
	resolver, err := NewResolver(context.Background(), gjson.New(config))
	if err != nil {
		t.Fatalf("NewResolver() error = %v", err)
	}
	return resolver
}
