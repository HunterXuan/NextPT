package catalog

import (
	"context"

	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

type sCatalogSubtitleDomain struct{}

func NewCatalogSubtitleDomain() *sCatalogSubtitleDomain {
	return &sCatalogSubtitleDomain{}
}

func init() {
	service.RegisterCatalogSubtitleDomain(NewCatalogSubtitleDomain())
}

func (s *sCatalogSubtitleDomain) InsertSubtitle(ctx context.Context, torrentId uint64, userId uint64, fileName, ext string, size int, language string, anonymous bool) (uint64, error) {
	insertData := g.Map{
		dao.CatalogSubtitle.Columns().TorrentId:     torrentId,
		dao.CatalogSubtitle.Columns().UserId:        userId,
		dao.CatalogSubtitle.Columns().Title:         fileName,
		dao.CatalogSubtitle.Columns().FileName:      fileName,
		dao.CatalogSubtitle.Columns().FileExt:       ext,
		dao.CatalogSubtitle.Columns().FileSize:      size,
		dao.CatalogSubtitle.Columns().StoragePath:   "",
		dao.CatalogSubtitle.Columns().Language:      language,
		dao.CatalogSubtitle.Columns().DownloadCount: 0,
		dao.CatalogSubtitle.Columns().Anonymous:     anonymous,
		dao.CatalogSubtitle.Columns().CreatedAt:     gtime.Now(),
	}

	id, err := dao.CatalogSubtitle.Ctx(ctx).Data(insertData).InsertAndGetId()
	return uint64(id), err
}

func (s *sCatalogSubtitleDomain) UpdateSubtitleStoragePath(ctx context.Context, id uint64, path string) error {
	_, err := dao.CatalogSubtitle.Ctx(ctx).Where(dao.CatalogSubtitle.Columns().Id, id).Data(g.Map{
		dao.CatalogSubtitle.Columns().StoragePath: path,
	}).Update()
	return err
}

func (s *sCatalogSubtitleDomain) GetSubtitleById(ctx context.Context, id uint64) (*entity.CatalogSubtitle, error) {
	var sub entity.CatalogSubtitle
	err := dao.CatalogSubtitle.Ctx(ctx).Where(dao.CatalogSubtitle.Columns().Id, id).Scan(&sub)
	if err != nil || sub.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "catalog.subtitle.not_found"))
	}
	return &sub, nil
}

func (s *sCatalogSubtitleDomain) IncrementDownloadCount(ctx context.Context, id uint64) error {
	_, err := dao.CatalogSubtitle.Ctx(ctx).Where(dao.CatalogSubtitle.Columns().Id, id).Increment(dao.CatalogSubtitle.Columns().DownloadCount, 1)
	return err
}

func (s *sCatalogSubtitleDomain) QuerySubtitles(ctx context.Context, torrentId uint64, page, size int) ([]entity.CatalogSubtitle, int, error) {
	m := dao.CatalogSubtitle.Ctx(ctx)
	if torrentId > 0 {
		m = m.Where(dao.CatalogSubtitle.Columns().TorrentId, torrentId)
	}

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var list []entity.CatalogSubtitle
	err = m.Page(page, size).OrderDesc(dao.CatalogSubtitle.Columns().CreatedAt).Scan(&list)
	return list, total, err
}

func (s *sCatalogSubtitleDomain) UpdateSubtitle(ctx context.Context, id uint64, language string, userId uint64, isAdmin bool) error {
	sub, err := s.GetSubtitleById(ctx, id)
	if err != nil {
		return err
	}
	if sub.UserId != userId && !isAdmin {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}

	_, err = dao.CatalogSubtitle.Ctx(ctx).Where(dao.CatalogSubtitle.Columns().Id, id).Data(g.Map{
		dao.CatalogSubtitle.Columns().Language: language,
	}).Update()
	return err
}

func (s *sCatalogSubtitleDomain) GetSubtitlesByTorrentId(ctx context.Context, torrentId uint64) ([]*entity.CatalogSubtitle, error) {
	var list []*entity.CatalogSubtitle
	err := dao.CatalogSubtitle.Ctx(ctx).Where(dao.CatalogSubtitle.Columns().TorrentId, torrentId).Scan(&list)
	return list, err
}

func (s *sCatalogSubtitleDomain) DeleteSubtitlesByTorrentId(ctx context.Context, torrentId uint64) error {
	_, err := dao.CatalogSubtitle.Ctx(ctx).Where(dao.CatalogSubtitle.Columns().TorrentId, torrentId).Delete()
	return err
}

func (s *sCatalogSubtitleDomain) DeleteSubtitle(ctx context.Context, id uint64) error {
	_, err := dao.CatalogSubtitle.Ctx(ctx).WherePri(id).Delete()
	return err
}
