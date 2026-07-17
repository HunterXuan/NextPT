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

func TestEconomyShopProductsRejectEnabledUnimplementedType(t *testing.T) {
	products := EconomyShopProducts{{
		Key:     "vip_30d",
		Type:    consts.EconomyShopProductTypeVip,
		Enabled: true,
		Price:   1000,
	}}.Normalized()

	if err := products.Validate(); err == nil {
		t.Fatal("expected enabled unimplemented product type to be rejected")
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
