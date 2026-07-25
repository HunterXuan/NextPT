package catalog

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"time"

	"server/internal/consts"
	libmetadata "server/internal/library/metadata"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sCatalogMetadataUsecase struct {
	tmdb    *libmetadata.TMDBProvider
	imdb    *libmetadata.IMDbProvider
	douban  *libmetadata.DoubanProvider
	bangumi *libmetadata.BangumiProvider
}

func init() {
	service.RegisterCatalogMetadataUsecase(NewCatalogMetadataUsecase())
}

func NewCatalogMetadataUsecase() *sCatalogMetadataUsecase {
	return &sCatalogMetadataUsecase{
		tmdb:    &libmetadata.TMDBProvider{},
		imdb:    &libmetadata.IMDbProvider{},
		douban:  &libmetadata.DoubanProvider{},
		bangumi: &libmetadata.BangumiProvider{},
	}
}

func (s *sCatalogMetadataUsecase) Search(ctx context.Context, actor *model.Actor, in catalogin.TorrentMetadataSearchInp) (*catalogout.TorrentMetadataSearchOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	if !slices.Contains(consts.CatalogMetadataTmdbTypes, in.TmdbType) {
		return nil, gerror.New(gi18n.T(ctx, "catalog.metadata.invalid_type"))
	}
	query := strings.TrimSpace(in.Query)
	if query == "" {
		return nil, gerror.New(gi18n.T(ctx, "catalog.metadata.query_required"))
	}
	page := in.Page
	if page <= 0 {
		page = 1
	}
	out, err := s.tmdb.Search(ctx, query, in.TmdbType, s.metadataLanguage(ctx), page)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.metadata.provider_failed"))
	}
	items := make([]catalogout.TorrentMetadataItem, 0, len(out.List))
	for i := range out.List {
		items = append(items, *s.metadataItemOut(&out.List[i]))
	}
	return &catalogout.TorrentMetadataSearchOut{
		List:         items,
		Page:         out.Page,
		TotalPages:   out.TotalPages,
		TotalResults: out.TotalResults,
	}, nil
}

func (s *sCatalogMetadataUsecase) GetTorrentMetadata(ctx context.Context, torrentId uint64) (*catalogout.TorrentMetadataOut, error) {
	meta, err := service.CatalogMetadataDomain().GetTorrentMeta(ctx, torrentId)
	if err != nil {
		return nil, err
	}
	out := &catalogout.TorrentMetadataOut{Sources: []catalogout.TorrentMetadataItem{}}
	if meta == nil {
		return out, nil
	}

	binding := s.bindingFromEntity(meta)
	originalBinding := binding
	var (
		imdbData    *catalogout.TorrentMetadataItem
		doubanData  *catalogout.TorrentMetadataItem
		bangumiData *catalogout.TorrentMetadataItem
	)

	// IMDb is the preferred mapping identity. Lower-priority providers are only
	// consulted for identity when neither a TMDB nor IMDb binding is available.
	if binding.TmdbId == "" && binding.ImdbId == "" && binding.DoubanId != "" {
		doubanData = s.getDoubanMetadata(ctx, binding.DoubanId)
		if doubanData != nil {
			binding.ImdbId = doubanData.ImdbId
		}
	}
	if binding.TmdbId == "" && binding.ImdbId == "" && binding.BangumiId != "" {
		bangumiData = s.getBangumiMetadata(ctx, binding.BangumiId)
		if bangumiData != nil {
			binding.ImdbId = bangumiData.ImdbId
		}
	}
	if binding.TmdbId == "" && binding.ImdbId != "" {
		binding = s.resolveTmdbByImdb(ctx, binding)
	}

	var tmdbData *catalogout.TorrentMetadataItem
	if binding.TmdbId != "" && slices.Contains(consts.CatalogMetadataTmdbTypes, binding.TmdbType) {
		tmdbData = s.getTmdbMetadata(ctx, binding.TmdbId, binding.TmdbType)
		if tmdbData != nil && binding.ImdbId == "" {
			binding.ImdbId = tmdbData.ImdbId
		}
	}
	if binding.DoubanId != "" && doubanData == nil {
		doubanData = s.getDoubanMetadata(ctx, binding.DoubanId)
		if doubanData != nil && binding.ImdbId == "" {
			binding.ImdbId = doubanData.ImdbId
		}
	}
	if binding.BangumiId != "" && bangumiData == nil {
		bangumiData = s.getBangumiMetadata(ctx, binding.BangumiId)
		if bangumiData != nil && binding.ImdbId == "" {
			binding.ImdbId = bangumiData.ImdbId
		}
	}
	if binding.ImdbId != "" {
		imdbData = s.getImdbMetadata(ctx, binding.ImdbId)
	}

	if binding != originalBinding {
		if err := service.CatalogMetadataDomain().UpsertTorrentBinding(ctx, torrentId, binding); err != nil {
			g.Log().Warningf(ctx, "catalog metadata: update resolved binding for torrent %d failed: %v", torrentId, err)
		}
	}

	out.Binding = s.bindingOut(meta, binding)
	if tmdbData != nil {
		out.Sources = append(out.Sources, *tmdbData)
		out.Binding.TmdbRating = tmdbData.Rating
		if meta.TmdbRating != tmdbData.Rating {
			_ = service.CatalogMetadataDomain().UpdateTmdbRating(ctx, binding.TmdbId, binding.TmdbType, tmdbData.Rating)
		}
	}
	if imdbData != nil {
		out.Sources = append(out.Sources, *imdbData)
		out.Binding.ImdbRating = imdbData.Rating
		if meta.ImdbRating != imdbData.Rating {
			_ = service.CatalogMetadataDomain().UpdateImdbRating(ctx, binding.ImdbId, imdbData.Rating)
		}
	}
	if doubanData != nil {
		out.Sources = append(out.Sources, *doubanData)
		out.Binding.DoubanRating = doubanData.Rating
		if meta.DoubanRating != doubanData.Rating {
			_ = service.CatalogMetadataDomain().UpdateDoubanRating(ctx, binding.DoubanId, doubanData.Rating)
		}
	}
	if bangumiData != nil {
		out.Sources = append(out.Sources, *bangumiData)
		out.Binding.BangumiRating = bangumiData.Rating
		if meta.BangumiRating != bangumiData.Rating {
			_ = service.CatalogMetadataDomain().UpdateBangumiRating(ctx, binding.BangumiId, bangumiData.Rating)
		}
	}
	out.Data = s.mergeMetadata(tmdbData, imdbData, doubanData, bangumiData)
	return out, nil
}

func (s *sCatalogMetadataUsecase) ResolveBinding(ctx context.Context, binding model.CatalogTorrentMetadataBinding) (model.CatalogTorrentMetadataBinding, error) {
	binding = binding.Normalized()
	if !binding.IsValid() {
		return binding, gerror.New(gi18n.T(ctx, "catalog.metadata.invalid_binding"))
	}
	if binding.TmdbId == "" && binding.ImdbId == "" && binding.DoubanId != "" {
		if data := s.getDoubanMetadata(ctx, binding.DoubanId); data != nil {
			binding.ImdbId = data.ImdbId
		}
	}
	if binding.TmdbId == "" && binding.ImdbId == "" && binding.BangumiId != "" {
		if data := s.getBangumiMetadata(ctx, binding.BangumiId); data != nil {
			binding.ImdbId = data.ImdbId
		}
	}
	if binding.TmdbId == "" && binding.ImdbId != "" {
		binding = s.resolveTmdbByImdb(ctx, binding)
	}
	if !binding.IsValid() {
		return binding, gerror.New(gi18n.T(ctx, "catalog.metadata.invalid_binding"))
	}
	return binding, nil
}

func (s *sCatalogMetadataUsecase) resolveTmdbByImdb(ctx context.Context, binding model.CatalogTorrentMetadataBinding) model.CatalogTorrentMetadataBinding {
	tmdbId, tmdbType, err := s.tmdb.FindByIMDb(ctx, binding.ImdbId, s.metadataLanguage(ctx))
	if err != nil {
		g.Log().Warningf(ctx, "catalog metadata: resolve imdb %s failed: %v", binding.ImdbId, err)
		return binding
	}
	if tmdbId != "" {
		binding.TmdbId = tmdbId
		binding.TmdbType = tmdbType
	}
	return binding
}

func (s *sCatalogMetadataUsecase) getTmdbMetadata(ctx context.Context, tmdbId string, tmdbType string) *catalogout.TorrentMetadataItem {
	locale := s.metadataLanguage(ctx)
	data, err := s.loadMetadata(ctx,
		service.SysCache().KeyCatalogMetadataTmdb(ctx, tmdbType, tmdbId, locale),
		"catalog.metadata.tmdb.cacheTtl",
		7*24*time.Hour,
		func() (*libmetadata.Item, error) {
			return s.tmdb.Get(ctx, libmetadata.GetInput{ExternalId: tmdbId, MediaType: tmdbType, Locale: locale})
		},
	)
	if err != nil {
		g.Log().Warningf(ctx, "catalog metadata: resolve tmdb %s/%s failed: %v", tmdbType, tmdbId, err)
		return nil
	}
	return data
}

func (s *sCatalogMetadataUsecase) getImdbMetadata(ctx context.Context, imdbId string) *catalogout.TorrentMetadataItem {
	if !g.Cfg().MustGet(ctx, "catalog.metadata.imdb.enabled", true).Bool() {
		return nil
	}
	data, err := s.loadMetadata(ctx,
		service.SysCache().KeyCatalogMetadataImdb(ctx, imdbId),
		"catalog.metadata.imdb.cacheTtl",
		7*24*time.Hour,
		func() (*libmetadata.Item, error) {
			return s.imdb.Get(ctx, libmetadata.GetInput{ExternalId: imdbId})
		},
	)
	if err != nil {
		g.Log().Warningf(ctx, "catalog metadata: resolve imdb metadata %s failed: %v", imdbId, err)
		return nil
	}
	return data
}

func (s *sCatalogMetadataUsecase) getDoubanMetadata(ctx context.Context, doubanId string) *catalogout.TorrentMetadataItem {
	if !g.Cfg().MustGet(ctx, "catalog.metadata.douban.enabled", true).Bool() {
		return nil
	}
	data, err := s.loadMetadata(ctx,
		service.SysCache().KeyCatalogMetadataDouban(ctx, doubanId),
		"catalog.metadata.douban.cacheTtl",
		30*24*time.Hour,
		func() (*libmetadata.Item, error) {
			return s.douban.Get(ctx, libmetadata.GetInput{ExternalId: doubanId})
		},
	)
	if err != nil {
		g.Log().Warningf(ctx, "catalog metadata: resolve douban %s failed: %v", doubanId, err)
		return nil
	}
	return data
}

func (s *sCatalogMetadataUsecase) getBangumiMetadata(ctx context.Context, bangumiId string) *catalogout.TorrentMetadataItem {
	if !g.Cfg().MustGet(ctx, "catalog.metadata.bangumi.enabled", true).Bool() {
		return nil
	}
	data, err := s.loadMetadata(ctx,
		service.SysCache().KeyCatalogMetadataBangumi(ctx, bangumiId),
		"catalog.metadata.bangumi.cacheTtl",
		30*24*time.Hour,
		func() (*libmetadata.Item, error) {
			return s.bangumi.Get(ctx, libmetadata.GetInput{ExternalId: bangumiId})
		},
	)
	if err != nil {
		g.Log().Warningf(ctx, "catalog metadata: resolve bangumi %s failed: %v", bangumiId, err)
		return nil
	}
	return data
}

func (s *sCatalogMetadataUsecase) loadMetadata(
	ctx context.Context,
	cacheKey string,
	ttlConfigKey string,
	fallbackTtl time.Duration,
	load func() (*libmetadata.Item, error),
) (*catalogout.TorrentMetadataItem, error) {
	if cached := s.metadataFromCache(ctx, cacheKey); cached != nil {
		return cached, nil
	}
	providerData, err := load()
	if err != nil || providerData == nil {
		return nil, err
	}
	data := s.metadataItemOut(providerData)
	s.cacheMetadata(ctx, cacheKey, data, s.metadataCacheTtl(ctx, ttlConfigKey, fallbackTtl))
	return data, nil
}

func (s *sCatalogMetadataUsecase) metadataFromCache(ctx context.Context, cacheKey string) *catalogout.TorrentMetadataItem {
	value, err := g.Redis().Do(ctx, "GET", cacheKey)
	if err != nil || value.IsNil() {
		return nil
	}
	var cached catalogout.TorrentMetadataItem
	if json.Unmarshal(value.Bytes(), &cached) != nil {
		return nil
	}
	return &cached
}

func (s *sCatalogMetadataUsecase) cacheMetadata(ctx context.Context, cacheKey string, data *catalogout.TorrentMetadataItem, ttl time.Duration) {
	payload, err := json.Marshal(data)
	if err != nil {
		return
	}
	_, _ = g.Redis().Do(ctx, "SET", cacheKey, payload, "EX", int(ttl.Seconds()))
}

func (s *sCatalogMetadataUsecase) metadataCacheTtl(ctx context.Context, configKey string, fallback time.Duration) time.Duration {
	ttl := g.Cfg().MustGet(ctx, configKey, fallback.String()).Duration()
	if ttl <= 0 {
		return fallback
	}
	return ttl
}

func (s *sCatalogMetadataUsecase) metadataItemOut(data *libmetadata.Item) *catalogout.TorrentMetadataItem {
	if data == nil {
		return nil
	}
	return &catalogout.TorrentMetadataItem{
		Provider:      data.Provider,
		ProviderId:    data.ProviderId,
		TmdbType:      data.TmdbType,
		Title:         data.Title,
		OriginalTitle: data.OriginalTitle,
		Year:          data.Year,
		ReleaseDate:   data.ReleaseDate,
		Overview:      data.Overview,
		PosterUrl:     data.PosterUrl,
		BackdropUrl:   data.BackdropUrl,
		Rating:        data.Rating,
		Genres:        append([]string(nil), data.Genres...),
		ImdbId:        data.ImdbId,
	}
}

func (s *sCatalogMetadataUsecase) metadataLanguage(ctx context.Context) string {
	if language := strings.TrimSpace(gi18n.LanguageFromCtx(ctx)); language != "" {
		return language
	}
	return g.Cfg().MustGet(ctx, "i18n.default", "zh-CN").String()
}

func (s *sCatalogMetadataUsecase) bindingFromEntity(meta *entity.CatalogTorrentMeta) model.CatalogTorrentMetadataBinding {
	return model.CatalogTorrentMetadataBinding{
		ImdbId:    meta.ImdbId,
		DoubanId:  meta.DoubanId,
		BangumiId: meta.BangumiId,
		TmdbId:    meta.TmdbId,
		TmdbType:  meta.TmdbType,
	}.Normalized()
}

func (s *sCatalogMetadataUsecase) bindingOut(meta *entity.CatalogTorrentMeta, binding model.CatalogTorrentMetadataBinding) catalogout.TorrentMetadataBinding {
	return catalogout.TorrentMetadataBinding{
		ImdbId:        binding.ImdbId,
		ImdbRating:    meta.ImdbRating,
		DoubanId:      binding.DoubanId,
		DoubanRating:  meta.DoubanRating,
		BangumiId:     binding.BangumiId,
		BangumiRating: meta.BangumiRating,
		TmdbId:        binding.TmdbId,
		TmdbType:      binding.TmdbType,
		TmdbRating:    meta.TmdbRating,
	}
}

func (s *sCatalogMetadataUsecase) mergeMetadata(sources ...*catalogout.TorrentMetadataItem) *catalogout.TorrentMetadataItem {
	var merged *catalogout.TorrentMetadataItem
	for _, source := range sources {
		if source == nil {
			continue
		}
		if merged == nil {
			copy := *source
			copy.Genres = append([]string(nil), source.Genres...)
			merged = &copy
			continue
		}
		if merged.Title == "" {
			merged.Title = source.Title
		}
		if merged.OriginalTitle == "" {
			merged.OriginalTitle = source.OriginalTitle
		}
		if merged.Year == "" {
			merged.Year = source.Year
		}
		if merged.ReleaseDate == "" {
			merged.ReleaseDate = source.ReleaseDate
		}
		if merged.Overview == "" {
			merged.Overview = source.Overview
		}
		if merged.PosterUrl == "" {
			merged.PosterUrl = source.PosterUrl
		}
		if merged.BackdropUrl == "" {
			merged.BackdropUrl = source.BackdropUrl
		}
		if len(merged.Genres) == 0 {
			merged.Genres = append([]string(nil), source.Genres...)
		}
		if merged.ImdbId == "" {
			merged.ImdbId = source.ImdbId
		}
	}
	return merged
}
