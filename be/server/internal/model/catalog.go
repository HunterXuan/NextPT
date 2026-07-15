package model

import (
	"math/rand"
	"sort"
	"strings"
	"time"

	"server/internal/consts"

	"github.com/gogf/gf/v2/os/gtime"
)

const catalogBytesPerGiB = 1024 * 1024 * 1024

type CatalogTorrentUpdate struct {
	Name          string
	SubTitle      string
	CategoryId    uint
	Description   string
	Anonymous     *bool
	ReleaseFields any
}

type CatalogTorrentSummary struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Size  uint64 `json:"size"`
	Exist bool   `json:"exist"`
}

type CatalogTorrentListOptions struct {
	Keyword             string
	CategoryIds         []uint
	Promotion           string
	SeedStatus          string
	FeaturedOnly        bool
	MinSize             uint64
	MaxSize             uint64
	PublishedWithinDays int
	Sort                string
	Page                int
	Size                int
}

func (o CatalogTorrentListOptions) Normalized() CatalogTorrentListOptions {
	o.Keyword = strings.TrimSpace(o.Keyword)
	o.CategoryIds = normalizeCatalogTorrentCategoryIds(o.CategoryIds)
	if o.Page <= 0 {
		o.Page = 1
	}
	if o.Size <= 0 {
		o.Size = 50
	}
	if o.Size > 100 {
		o.Size = 100
	}
	if o.PublishedWithinDays < 0 {
		o.PublishedWithinDays = 0
	}

	switch o.Promotion {
	case consts.CatalogTorrentPromotionFilterPromoted,
		consts.ResourceTorrentPromotionStateNormal,
		consts.ResourceTorrentPromotionStateFree,
		consts.ResourceTorrentPromotionState2x,
		consts.ResourceTorrentPromotionState2xFree,
		consts.ResourceTorrentPromotionState50Percent,
		consts.ResourceTorrentPromotionState2x50Percent,
		consts.ResourceTorrentPromotionState30Percent:
	default:
		o.Promotion = consts.CatalogTorrentPromotionFilterAll
	}

	switch o.SeedStatus {
	case consts.CatalogTorrentSeedStatusSeeded, consts.CatalogTorrentSeedStatusUnseeded:
	default:
		o.SeedStatus = consts.CatalogTorrentSeedStatusAll
	}

	switch o.Sort {
	case consts.CatalogTorrentSortOldest,
		consts.CatalogTorrentSortSeeders,
		consts.CatalogTorrentSortLeechers,
		consts.CatalogTorrentSortComplete,
		consts.CatalogTorrentSortSizeAsc,
		consts.CatalogTorrentSortSizeDesc:
	default:
		o.Sort = consts.CatalogTorrentSortNewest
	}
	return o
}

func (o CatalogTorrentListOptions) HasInvalidSizeRange() bool {
	return o.MinSize > 0 && o.MaxSize > 0 && o.MinSize > o.MaxSize
}

func normalizeCatalogTorrentCategoryIds(categoryIds []uint) []uint {
	if len(categoryIds) == 0 {
		return nil
	}

	seen := make(map[uint]struct{}, len(categoryIds))
	list := make([]uint, 0, len(categoryIds))
	for _, id := range categoryIds {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		list = append(list, id)
	}
	return list
}

type CatalogRequestListOptions struct {
	ActorId     uint64
	Keyword     string
	RequestType uint
	Status      *uint
	CategoryId  uint
	View        string
	Page        int
	Size        int
}

type CatalogRequestActions struct {
	CanClaim    bool `json:"canClaim"`
	CanAbandon  bool `json:"canAbandon"`
	CanSubmit   bool `json:"canSubmit"`
	CanComplete bool `json:"canComplete"`
	CanCancel   bool `json:"canCancel"`
}

type CatalogTorrentPromotion struct {
	SpState    int
	SpExpireAt *gtime.Time
}

type CatalogTorrentPromotionFactor struct {
	Upload   float64
	Download float64
}

type CatalogTorrentGlobalPromotionConfig struct {
	Enabled  bool   `json:"enabled"`
	State    string `json:"state"`
	ExpireAt string `json:"expireAt"`
}

type CatalogTorrentNewPromotionConfig struct {
	Enabled bool                          `json:"enabled"`
	Rules   []CatalogTorrentPromotionRule `json:"rules"`
}

type CatalogTorrentPromotionRule struct {
	MinGiB        float64                         `json:"minGiB"`
	DurationHours int                             `json:"durationHours"`
	Options       []CatalogTorrentPromotionOption `json:"options"`
}

type CatalogTorrentPromotionOption struct {
	State  string `json:"state"`
	Weight int    `json:"weight"`
}

func CatalogTorrentPromotionStateToSp(state string) int {
	state = strings.TrimSpace(state)
	if state == "" {
		return consts.ResourceTorrentSpNormal
	}
	if sp, ok := consts.ResourceTorrentPromotionStateToSp[state]; ok {
		return sp
	}
	return consts.ResourceTorrentSpNormal
}

func CatalogTorrentPromotionFactorBySp(spState int) CatalogTorrentPromotionFactor {
	switch spState {
	case consts.ResourceTorrentSpFree:
		return CatalogTorrentPromotionFactor{Upload: 1, Download: 0}
	case consts.ResourceTorrentSp2x:
		return CatalogTorrentPromotionFactor{Upload: 2, Download: 1}
	case consts.ResourceTorrentSp2xFree:
		return CatalogTorrentPromotionFactor{Upload: 2, Download: 0}
	case consts.ResourceTorrentSp50Off:
		return CatalogTorrentPromotionFactor{Upload: 1, Download: 0.5}
	case consts.ResourceTorrentSp2x50Off:
		return CatalogTorrentPromotionFactor{Upload: 2, Download: 0.5}
	case consts.ResourceTorrentSp30Off:
		return CatalogTorrentPromotionFactor{Upload: 1, Download: 0.3}
	default:
		return CatalogTorrentPromotionFactor{Upload: 1, Download: 1}
	}
}

func (p CatalogTorrentPromotion) ApplyTraffic(rawUploaded, rawDownloaded int64) (int64, int64) {
	factor := CatalogTorrentPromotionFactorBySp(p.SpState)
	return int64(float64(rawUploaded) * factor.Upload), int64(float64(rawDownloaded) * factor.Download)
}

func ResolveCatalogTorrentPromotion(torrentSpState int, torrentSpExpireAt *gtime.Time, globalConfig CatalogTorrentGlobalPromotionConfig, now *gtime.Time) CatalogTorrentPromotion {
	if now == nil {
		now = gtime.Now()
	}
	global := CatalogTorrentPromotion{
		SpState:    CatalogTorrentPromotionStateToSp(globalConfig.State),
		SpExpireAt: globalConfig.expireTime(),
	}
	if globalConfig.Enabled && global.activeAt(now) {
		return global
	}

	torrent := CatalogTorrentPromotion{
		SpState:    torrentSpState,
		SpExpireAt: torrentSpExpireAt,
	}
	if torrent.activeAt(now) {
		return torrent
	}

	return CatalogTorrentPromotion{
		SpState: consts.ResourceTorrentSpNormal,
	}
}

func PickCatalogNewTorrentPromotion(config CatalogTorrentNewPromotionConfig, size uint64, now *gtime.Time) CatalogTorrentPromotion {
	if now == nil {
		now = gtime.Now()
	}
	if !config.Enabled {
		return CatalogTorrentPromotion{SpState: consts.ResourceTorrentSpNormal}
	}

	rule, ok := config.matchRule(size)
	if !ok {
		return CatalogTorrentPromotion{SpState: consts.ResourceTorrentSpNormal}
	}

	spState := pickCatalogTorrentPromotionState(rule.Options)
	if spState == consts.ResourceTorrentSpNormal {
		return CatalogTorrentPromotion{SpState: consts.ResourceTorrentSpNormal}
	}
	if rule.DurationHours < 0 {
		return CatalogTorrentPromotion{SpState: consts.ResourceTorrentSpNormal}
	}

	promotion := CatalogTorrentPromotion{SpState: spState}
	if rule.DurationHours > 0 {
		promotion.SpExpireAt = now.Add(time.Duration(rule.DurationHours) * time.Hour)
	}
	return promotion
}

func (p CatalogTorrentPromotion) activeAt(now *gtime.Time) bool {
	if p.SpState == consts.ResourceTorrentSpNormal {
		return false
	}
	if now == nil {
		now = gtime.Now()
	}
	return p.SpExpireAt == nil || p.SpExpireAt.After(now)
}

func (c CatalogTorrentGlobalPromotionConfig) expireTime() *gtime.Time {
	expireAt := strings.TrimSpace(c.ExpireAt)
	if expireAt == "" {
		return nil
	}
	t, err := gtime.StrToTime(expireAt)
	if err != nil {
		return nil
	}
	return t
}

func (c CatalogTorrentNewPromotionConfig) matchRule(size uint64) (CatalogTorrentPromotionRule, bool) {
	if len(c.Rules) == 0 {
		return CatalogTorrentPromotionRule{}, false
	}

	sizeGiB := float64(size) / catalogBytesPerGiB
	rules := append([]CatalogTorrentPromotionRule(nil), c.Rules...)
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].MinGiB > rules[j].MinGiB
	})
	for _, rule := range rules {
		if sizeGiB >= rule.MinGiB && len(rule.Options) > 0 {
			return rule, true
		}
	}
	return CatalogTorrentPromotionRule{}, false
}

func pickCatalogTorrentPromotionState(options []CatalogTorrentPromotionOption) int {
	total := 0
	for _, option := range options {
		if option.Weight <= 0 {
			continue
		}
		total += option.Weight
	}
	if total <= 0 {
		return consts.ResourceTorrentSpNormal
	}

	point := rand.Intn(total)
	for _, option := range options {
		if option.Weight <= 0 {
			continue
		}
		if point < option.Weight {
			return CatalogTorrentPromotionStateToSp(option.State)
		}
		point -= option.Weight
	}
	return consts.ResourceTorrentSpNormal
}

type CatalogUploadConfig struct {
	Title  CatalogUploadTitleConfig   `json:"title"`
	Fields []CatalogUploadFieldConfig `json:"fields"`
}

type CatalogUploadTitleConfig struct {
	Mode                string                   `json:"mode"`
	AllowManualOverride bool                     `json:"allowManualOverride"`
	Parts               []CatalogUploadTitlePart `json:"parts"`
}

type CatalogUploadTitlePart struct {
	Field     string `json:"field"`
	Prefix    string `json:"prefix"`
	Suffix    string `json:"suffix"`
	Separator string `json:"separator"`
}

type CatalogUploadFieldConfig struct {
	Key         string                     `json:"key"`
	Type        string                     `json:"type"`
	Label       map[string]any             `json:"label"`
	Description map[string]any             `json:"description,omitempty"`
	Placeholder map[string]any             `json:"placeholder,omitempty"`
	Required    bool                       `json:"required"`
	Options     *CatalogUploadFieldOptions `json:"options,omitempty"`
}

type CatalogUploadFieldOptions struct {
	Source string                    `json:"source"`
	Slug   string                    `json:"slug,omitempty"`
	Items  []CatalogUploadOptionItem `json:"items,omitempty"`
}

type CatalogUploadOptionItem struct {
	Value string         `json:"value"`
	Label map[string]any `json:"label"`
}
