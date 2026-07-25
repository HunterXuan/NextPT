package catalog

import (
	"context"

	"server/internal/dao"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sCatalogMetadataDomain struct{}

func init() {
	service.RegisterCatalogMetadataDomain(NewCatalogMetadataDomain())
}

func NewCatalogMetadataDomain() *sCatalogMetadataDomain {
	return &sCatalogMetadataDomain{}
}

func (s *sCatalogMetadataDomain) GetTorrentMeta(ctx context.Context, torrentId uint64) (*entity.CatalogTorrentMeta, error) {
	var meta entity.CatalogTorrentMeta
	err := dao.CatalogTorrentMeta.Ctx(ctx).
		Where(dao.CatalogTorrentMeta.Columns().TorrentId, torrentId).
		Scan(&meta)
	if err != nil {
		return nil, err
	}
	if meta.Id == 0 {
		return nil, nil
	}
	return &meta, nil
}

func (s *sCatalogMetadataDomain) UpsertTorrentBinding(ctx context.Context, torrentId uint64, binding model.CatalogTorrentMetadataBinding) error {
	binding = binding.Normalized()
	meta, err := s.GetTorrentMeta(ctx, torrentId)
	if err != nil {
		return err
	}

	columns := dao.CatalogTorrentMeta.Columns()
	data := g.Map{
		columns.TorrentId: torrentId,
		columns.ImdbId:    binding.ImdbId,
		columns.DoubanId:  binding.DoubanId,
		columns.BangumiId: binding.BangumiId,
		columns.TmdbId:    binding.TmdbId,
		columns.TmdbType:  binding.TmdbType,
	}
	if meta != nil {
		if meta.ImdbId != binding.ImdbId {
			data[columns.ImdbRating] = nil
		}
		if meta.DoubanId != binding.DoubanId {
			data[columns.DoubanRating] = nil
		}
		if meta.BangumiId != binding.BangumiId {
			data[columns.BangumiRating] = nil
		}
		if meta.TmdbId != binding.TmdbId || meta.TmdbType != binding.TmdbType {
			data[columns.TmdbRating] = nil
		}
	}

	_, err = dao.CatalogTorrentMeta.Ctx(ctx).Data(data).Save()
	return err
}

func (s *sCatalogMetadataDomain) UpdateTmdbRating(ctx context.Context, tmdbId string, tmdbType string, rating float64) error {
	if tmdbId == "" || tmdbType == "" {
		return nil
	}
	columns := dao.CatalogTorrentMeta.Columns()
	_, err := dao.CatalogTorrentMeta.Ctx(ctx).
		Where(columns.TmdbId, tmdbId).
		Where(columns.TmdbType, tmdbType).
		Data(columns.TmdbRating, rating).
		Update()
	return err
}

func (s *sCatalogMetadataDomain) UpdateImdbRating(ctx context.Context, imdbId string, rating float64) error {
	if imdbId == "" {
		return nil
	}
	columns := dao.CatalogTorrentMeta.Columns()
	_, err := dao.CatalogTorrentMeta.Ctx(ctx).
		Where(columns.ImdbId, imdbId).
		Data(columns.ImdbRating, rating).
		Update()
	return err
}

func (s *sCatalogMetadataDomain) UpdateDoubanRating(ctx context.Context, doubanId string, rating float64) error {
	if doubanId == "" {
		return nil
	}
	columns := dao.CatalogTorrentMeta.Columns()
	_, err := dao.CatalogTorrentMeta.Ctx(ctx).
		Where(columns.DoubanId, doubanId).
		Data(columns.DoubanRating, rating).
		Update()
	return err
}

func (s *sCatalogMetadataDomain) UpdateBangumiRating(ctx context.Context, bangumiId string, rating float64) error {
	if bangumiId == "" {
		return nil
	}
	columns := dao.CatalogTorrentMeta.Columns()
	_, err := dao.CatalogTorrentMeta.Ctx(ctx).
		Where(columns.BangumiId, bangumiId).
		Data(columns.BangumiRating, rating).
		Update()
	return err
}
