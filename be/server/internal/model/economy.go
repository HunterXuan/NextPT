package model

import (
	"math"
	"slices"
	"strings"

	"server/internal/consts"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
)

type EconomyShopProductConfig struct {
	Key       string         `json:"key"`
	Type      string         `json:"type"`
	Enabled   bool           `json:"enabled"`
	Price     float64        `json:"price"`
	SortOrder int            `json:"sortOrder"`
	Options   map[string]any `json:"options"`
}

type EconomyShopProducts []EconomyShopProductConfig

func NewEconomyShopProductSnapshot(key string, productType string, price float64) EconomyShopProductConfig {
	return EconomyShopProductConfig{
		Key:     key,
		Type:    productType,
		Enabled: true,
		Price:   math.Round(price*100) / 100,
		Options: map[string]any{},
	}
}

func (products EconomyShopProducts) Normalized() EconomyShopProducts {
	result := make(EconomyShopProducts, len(products))
	for i, product := range products {
		product.Key = strings.TrimSpace(product.Key)
		product.Type = strings.TrimSpace(product.Type)
		product.Price = math.Round(product.Price*100) / 100
		if product.Options == nil {
			product.Options = map[string]any{}
		}
		result[i] = product
	}
	return result
}

func (products EconomyShopProducts) Find(key string) *EconomyShopProductConfig {
	key = strings.TrimSpace(key)
	for i := range products {
		if products[i].Key == key {
			return &products[i]
		}
	}
	return nil
}

func (products EconomyShopProducts) Enabled() EconomyShopProducts {
	result := make(EconomyShopProducts, 0, len(products))
	for _, product := range products {
		if product.Enabled {
			result = append(result, product)
		}
	}
	return result
}

func (products EconomyShopProducts) Validate() error {
	seen := make(map[string]struct{}, len(products))
	for _, product := range products {
		if product.Key == "" || product.Type == "" {
			return gerror.New("shop product key and type are required")
		}
		if _, ok := seen[product.Key]; ok {
			return gerror.New("shop product key must be unique")
		}
		seen[product.Key] = struct{}{}
		if !slices.Contains(consts.EconomyShopProductTypes, product.Type) {
			return gerror.New("unsupported shop product type")
		}
		if product.Enabled && !slices.Contains(consts.EconomyShopImplementedProductTypes, product.Type) {
			return gerror.New("enabled shop product type is not implemented")
		}
		if product.Price <= 0 {
			return gerror.New("shop product price must be greater than 0")
		}
		if product.Type == consts.EconomyShopProductTypeInvite && gvar.New(product.Options["amount"]).Int() <= 0 {
			return gerror.New("invite product amount must be greater than 0")
		}
	}
	return nil
}

func (product EconomyShopProductConfig) InviteAmount() int {
	return gvar.New(product.Options["amount"]).Int()
}
