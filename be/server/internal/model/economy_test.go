package model

import (
	"testing"

	"server/internal/consts"
)

func TestEconomyShopProductsNormalizedAndValidate(t *testing.T) {
	products := EconomyShopProducts{{
		Key:     " invite ",
		Type:    " invite ",
		Enabled: true,
		Price:   1000.126,
		Options: map[string]any{"amount": 1},
	}}.Normalized()

	if products[0].Key != "invite" || products[0].Type != "invite" {
		t.Fatalf("product key/type were not normalized: %#v", products[0])
	}
	if products[0].Price != 1000.13 {
		t.Fatalf("product price was not normalized: %v", products[0].Price)
	}
	if err := products.Validate(); err != nil {
		t.Fatalf("expected valid product config, got %v", err)
	}
}

func TestEconomyShopProductsRejectUnsupportedType(t *testing.T) {
	products := EconomyShopProducts{{
		Key:     "unknown",
		Type:    "unknown",
		Enabled: true,
		Price:   1000,
	}}.Normalized()

	if err := products.Validate(); err == nil {
		t.Fatal("expected unsupported product type to be rejected")
	}
}

func TestEconomyShopProductsValidateImplementedOptions(t *testing.T) {
	products := EconomyShopProducts{
		{
			Key:     "vip_30d",
			Type:    consts.EconomyShopProductTypeVip,
			Enabled: true,
			Price:   3000,
			Options: map[string]any{"durationDays": 30},
		},
		{
			Key:     "upload_100_gib",
			Type:    consts.EconomyShopProductTypeUpload,
			Enabled: true,
			Price:   500,
			Options: map[string]any{"amountGiB": 100},
		},
		{
			Key:     "download_50_gib",
			Type:    consts.EconomyShopProductTypeDownload,
			Enabled: true,
			Price:   800,
			Options: map[string]any{"amountGiB": 50},
		},
	}.Normalized()

	if err := products.Validate(); err != nil {
		t.Fatalf("expected implemented product options to be valid, got %v", err)
	}
	if products[0].VipDurationDays() != 30 {
		t.Fatalf("unexpected vip duration: %d", products[0].VipDurationDays())
	}
	if products[1].TrafficBytes() != 100*1024*1024*1024 {
		t.Fatalf("unexpected traffic bytes: %d", products[1].TrafficBytes())
	}
}

func TestEconomyShopProductsRejectInvalidOptions(t *testing.T) {
	tests := []struct {
		name        string
		productType string
		options     map[string]any
	}{
		{name: "invite amount", productType: consts.EconomyShopProductTypeInvite, options: map[string]any{"amount": 0}},
		{name: "vip duration", productType: consts.EconomyShopProductTypeVip, options: map[string]any{"durationDays": 0}},
		{name: "upload amount", productType: consts.EconomyShopProductTypeUpload, options: map[string]any{"amountGiB": 1.5}},
		{name: "download amount", productType: consts.EconomyShopProductTypeDownload, options: map[string]any{"amountGiB": -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			products := EconomyShopProducts{{
				Key:     tt.name,
				Type:    tt.productType,
				Enabled: true,
				Price:   100,
				Options: tt.options,
			}}.Normalized()
			if err := products.Validate(); err == nil {
				t.Fatal("expected invalid product options to be rejected")
			}
		})
	}
}

func TestNewEconomyShopProductSnapshot(t *testing.T) {
	product := NewEconomyShopProductSnapshot("invite", consts.EconomyShopProductTypeInvite, 1000.126)

	if product.Key != "invite" || product.Type != consts.EconomyShopProductTypeInvite {
		t.Fatalf("unexpected fallback product identity: %#v", product)
	}
	if product.Price != 1000.13 {
		t.Fatalf("fallback product price was not normalized: %v", product.Price)
	}
	if !product.Enabled || product.Options == nil {
		t.Fatalf("fallback product should remain renderable: %#v", product)
	}
}
