package catalog

import (
	"context"
	"fmt"
	"strings"
	"time"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sCatalogTorrentDomain struct{}

func init() {
	service.RegisterCatalogTorrentDomain(NewCatalogTorrentDomain())
}

func NewCatalogTorrentDomain() *sCatalogTorrentDomain {
	return &sCatalogTorrentDomain{}
}

func (s *sCatalogTorrentDomain) GetTorrentById(ctx context.Context, id uint64) (*entity.CatalogTorrent, error) {
	var torrent entity.CatalogTorrent
	err := dao.CatalogTorrent.Ctx(ctx).Where(dao.CatalogTorrent.Columns().Id, id).Scan(&torrent)
	if err != nil || torrent.Id == 0 {
		return nil, gerror.New(gi18n.T(ctx, "catalog.torrent.not_found"))
	}
	return &torrent, nil
}

func (s *sCatalogTorrentDomain) GetTorrentByInfoHash(ctx context.Context, infoHash string) (*entity.CatalogTorrent, error) {
	var torrent *entity.CatalogTorrent
	err := dao.CatalogTorrent.Ctx(ctx).
		Where(dao.CatalogTorrent.Columns().InfoHash, infoHash).
		Scan(&torrent)

	if err != nil {
		return nil, err
	}
	return torrent, nil
}

func (s *sCatalogTorrentDomain) ResolveEffectiveTorrentPromotion(ctx context.Context, torrent *entity.CatalogTorrent, now *gtime.Time) model.CatalogTorrentPromotion {
	if torrent == nil {
		return model.CatalogTorrentPromotion{SpState: consts.ResourceTorrentSpNormal}
	}
	return model.ResolveCatalogTorrentPromotion(torrent.SpState, torrent.SpExpireAt, s.getGlobalPromotionConfig(ctx), now)
}

func (s *sCatalogTorrentDomain) PickNewTorrentPromotion(ctx context.Context, size uint64, now *gtime.Time) model.CatalogTorrentPromotion {
	return model.PickCatalogNewTorrentPromotion(s.getNewTorrentPromotionConfig(ctx), size, now)
}

func (s *sCatalogTorrentDomain) CalculatePromotedTorrentTraffic(ctx context.Context, torrent *entity.CatalogTorrent, rawUploaded, rawDownloaded int64, now *gtime.Time) (int64, int64) {
	return s.ResolveEffectiveTorrentPromotion(ctx, torrent, now).ApplyTraffic(rawUploaded, rawDownloaded)
}

func (s *sCatalogTorrentDomain) getGlobalPromotionConfig(ctx context.Context) model.CatalogTorrentGlobalPromotionConfig {
	var cfg model.CatalogTorrentGlobalPromotionConfig
	_ = gconv.Struct(s.getCatalogConfigCache(ctx, consts.SiteConfigCatalogGlobalPromotion).Val(), &cfg)
	return cfg
}

func (s *sCatalogTorrentDomain) getNewTorrentPromotionConfig(ctx context.Context) model.CatalogTorrentNewPromotionConfig {
	var cfg model.CatalogTorrentNewPromotionConfig
	_ = gconv.Struct(s.getCatalogConfigCache(ctx, consts.SiteConfigCatalogNewTorrentPromotion).Val(), &cfg)
	return cfg
}

func (s *sCatalogTorrentDomain) getCatalogConfigCache(ctx context.Context, key string) *gvar.Var {
	cacheKey := service.SysCache().KeySiteConfigFullPath(ctx, key)
	val, err := gcache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		return service.SiteConfigDomain().GetByPath(ctx, key).Val(), nil
	}, 5*time.Minute)
	if err != nil || val.IsNil() {
		return service.SiteConfigDomain().GetByPath(ctx, key)
	}
	return gvar.New(val.Val())
}

func (s *sCatalogTorrentDomain) CheckCategoryExists(ctx context.Context, categoryId uint) error {
	count, err := dao.CatalogCategory.Ctx(ctx).Where(dao.CatalogCategory.Columns().Id, categoryId).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New(gi18n.T(ctx, "catalog.category.invalid"))
	}
	return nil
}

func (s *sCatalogTorrentDomain) CountTorrentsByCategory(ctx context.Context, categoryId uint) (int, error) {
	return dao.CatalogTorrent.Ctx(ctx).Where(dao.CatalogTorrent.Columns().CategoryId, categoryId).Count()
}

func (s *sCatalogTorrentDomain) CheckInfoHashExists(ctx context.Context, infoHashBytes []byte) error {
	count, err := dao.CatalogTorrent.Ctx(ctx).Where(dao.CatalogTorrent.Columns().InfoHash, infoHashBytes).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New(gi18n.T(ctx, "catalog.torrent.duplicate"))
	}
	return nil
}

func (s *sCatalogTorrentDomain) SaveTorrent(ctx context.Context, torrentInsert *entity.CatalogTorrent, filesToInsert []entity.CatalogTorrentFile) (uint64, error) {
	insertRes, err := dao.CatalogTorrent.Ctx(ctx).Data(torrentInsert).Insert()
	if err != nil {
		return 0, err
	}
	tid, _ := insertRes.LastInsertId()
	torrentId := uint64(tid)

	for i := range filesToInsert {
		filesToInsert[i].TorrentId = torrentId
	}

	if len(filesToInsert) > 0 {
		_, err = dao.CatalogTorrentFile.Ctx(ctx).Data(filesToInsert).Batch(500).Insert()
		if err != nil {
			return 0, err
		}
	}

	_, err = dao.CatalogTorrentMeta.Ctx(ctx).Data(&entity.CatalogTorrentMeta{
		TorrentId: torrentId,
	}).Insert()

	return torrentId, err
}

func (s *sCatalogTorrentDomain) Bookmark(ctx context.Context, userId uint64, torrentId uint64) error {
	count, err := dao.CatalogTorrentBookmark.Ctx(ctx).Where(g.Map{
		dao.CatalogTorrentBookmark.Columns().UserId:    userId,
		dao.CatalogTorrentBookmark.Columns().TorrentId: torrentId,
	}).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Already bookmarked
	}
	_, err = dao.CatalogTorrentBookmark.Ctx(ctx).Data(entity.CatalogTorrentBookmark{
		UserId:    userId,
		TorrentId: torrentId,
	}).Insert()
	return err
}

func (s *sCatalogTorrentDomain) Unbookmark(ctx context.Context, userId uint64, torrentId uint64) error {
	_, err := dao.CatalogTorrentBookmark.Ctx(ctx).Where(g.Map{
		dao.CatalogTorrentBookmark.Columns().UserId:    userId,
		dao.CatalogTorrentBookmark.Columns().TorrentId: torrentId,
	}).Delete()
	return err
}

func (s *sCatalogTorrentDomain) QueryBookmarkedTorrents(ctx context.Context, actor *model.Actor, page, size int) ([]entity.CatalogTorrent, int, error) {
	userId := actor.Id

	// 1. 获取用户收藏的所有种子 ID，按收藏时间倒序
	var bookmarks []entity.CatalogTorrentBookmark
	err := dao.CatalogTorrentBookmark.Ctx(ctx).
		Where(dao.CatalogTorrentBookmark.Columns().UserId, userId).
		OrderDesc(dao.CatalogTorrentBookmark.Columns().CreatedAt).
		Scan(&bookmarks)
	if err != nil || len(bookmarks) == 0 {
		return nil, 0, err
	}

	bookmarkIds := make([]uint64, 0, len(bookmarks))
	for _, bookmark := range bookmarks {
		bookmarkIds = append(bookmarkIds, bookmark.TorrentId)
	}

	// 2. 批量查询这些种子中，当前用户可见的种子 ID
	var visibleTorrents []entity.CatalogTorrent
	m := dao.CatalogTorrent.Ctx(ctx).
		Fields(dao.CatalogTorrent.Columns().Id).
		WhereIn(dao.CatalogTorrent.Columns().Id, bookmarkIds)
	m = s.ApplyTorrentVisibleScope(m, actor)
	err = m.Scan(&visibleTorrents)
	if err != nil || len(visibleTorrents) == 0 {
		return nil, 0, err
	}

	// 将可见种子 ID 转成 Map，方便快速检索
	visibleMap := make(map[uint64]bool)
	for _, torrent := range visibleTorrents {
		visibleMap[torrent.Id] = true
	}

	// 3. 过滤出既被收藏又对当前用户可见的种子 ID（保持原收藏倒序）
	var finalIds []uint64
	for _, id := range bookmarkIds {
		if visibleMap[id] {
			finalIds = append(finalIds, id)
		}
	}

	total := len(finalIds)
	if total == 0 {
		return nil, 0, nil
	}

	// 4. 执行内存分页切片
	start := (page - 1) * size
	if start >= total {
		return nil, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	pageIds := finalIds[start:end]

	// 5. 批量查询分页后种子的详细数据
	var unorderedTorrents []entity.CatalogTorrent
	err = dao.CatalogTorrent.Ctx(ctx).
		WhereIn(dao.CatalogTorrent.Columns().Id, pageIds).
		Scan(&unorderedTorrents)
	if err != nil {
		return nil, 0, err
	}

	// 6. 将查询出来的详细数据按 pageIds 顺序重排
	torrentMap := make(map[uint64]entity.CatalogTorrent)
	for _, t := range unorderedTorrents {
		torrentMap[t.Id] = t
	}

	var torrents []entity.CatalogTorrent
	for _, id := range pageIds {
		if t, ok := torrentMap[id]; ok {
			torrents = append(torrents, t)
		}
	}

	return torrents, total, nil
}

func (s *sCatalogTorrentDomain) ToggleLike(ctx context.Context, torrentId uint64, userId uint64) (bool, error) {
	likeColumns := dao.CatalogTorrentLike.Columns()
	torrentColumns := dao.CatalogTorrent.Columns()

	isLiked := false
	err := dao.CatalogTorrentLike.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		m := dao.CatalogTorrentLike.Ctx(ctx).
			Where(likeColumns.TorrentId, torrentId).
			Where(likeColumns.UserId, userId)

		count, err := m.Count()
		if err != nil {
			return err
		}

		if count > 0 {
			if _, err = m.Delete(); err != nil {
				return err
			}
			if _, err = dao.CatalogTorrent.Ctx(ctx).Where(torrentColumns.Id, torrentId).Decrement(torrentColumns.LikeCount, 1); err != nil {
				return err
			}
			isLiked = false
			return nil
		}

		if _, err = dao.CatalogTorrentLike.Ctx(ctx).Data(g.Map{
			likeColumns.TorrentId: torrentId,
			likeColumns.UserId:    userId,
			likeColumns.CreatedAt: gtime.Now(),
		}).Insert(); err != nil {
			return err
		}
		if _, err = dao.CatalogTorrent.Ctx(ctx).Where(torrentColumns.Id, torrentId).Increment(torrentColumns.LikeCount, 1); err != nil {
			return err
		}
		isLiked = true
		return nil
	})
	return isLiked, err
}

func (s *sCatalogTorrentDomain) QueryTorrentLikes(ctx context.Context, torrentId uint64, page, size int) ([]entity.CatalogTorrentLike, int, error) {
	m := dao.CatalogTorrentLike.Ctx(ctx).Where(dao.CatalogTorrentLike.Columns().TorrentId, torrentId)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []entity.CatalogTorrentLike
	err = m.Page(page, size).OrderDesc(dao.CatalogTorrentLike.Columns().CreatedAt).Scan(&list)
	return list, total, err
}

func (s *sCatalogTorrentDomain) UpdateTorrent(ctx context.Context, id uint64, data model.CatalogTorrentUpdate) error {
	updateMap := g.Map{
		dao.CatalogTorrent.Columns().Name:          data.Name,
		dao.CatalogTorrent.Columns().SubTitle:      data.SubTitle,
		dao.CatalogTorrent.Columns().CategoryId:    data.CategoryId,
		dao.CatalogTorrent.Columns().Description:   data.Description,
		dao.CatalogTorrent.Columns().ReleaseFields: data.ReleaseFields,
	}
	if data.Anonymous != nil {
		if *data.Anonymous {
			updateMap[dao.CatalogTorrent.Columns().Anonymous] = 1
		} else {
			updateMap[dao.CatalogTorrent.Columns().Anonymous] = 0
		}
	}

	_, err := dao.CatalogTorrent.Ctx(ctx).Where(dao.CatalogTorrent.Columns().Id, id).Data(updateMap).Update()
	return err
}

func (s *sCatalogTorrentDomain) AdminSetTorrentPinned(ctx context.Context, id uint64, pinned bool, pinWeight int) error {
	if !pinned {
		pinWeight = 0
	}
	_, err := dao.CatalogTorrent.Ctx(ctx).Where(dao.CatalogTorrent.Columns().Id, id).Data(g.Map{
		dao.CatalogTorrent.Columns().IsPinned:  pinned,
		dao.CatalogTorrent.Columns().PinWeight: pinWeight,
		dao.CatalogTorrent.Columns().UpdatedAt: gtime.Now(),
	}).Update()
	return err
}

func (s *sCatalogTorrentDomain) AdminSetTorrentFeatured(ctx context.Context, id uint64, featured bool) error {
	_, err := dao.CatalogTorrent.Ctx(ctx).Where(dao.CatalogTorrent.Columns().Id, id).Data(g.Map{
		dao.CatalogTorrent.Columns().IsFeatured: featured,
		dao.CatalogTorrent.Columns().UpdatedAt:  gtime.Now(),
	}).Update()
	return err
}

func (s *sCatalogTorrentDomain) AdminSetTorrentPromotion(ctx context.Context, id uint64, spState int, spExpireAt *gtime.Time) error {
	if spState == consts.ResourceTorrentSpNormal {
		spExpireAt = nil
	}
	_, err := dao.CatalogTorrent.Ctx(ctx).Where(dao.CatalogTorrent.Columns().Id, id).Data(g.Map{
		dao.CatalogTorrent.Columns().SpState:    spState,
		dao.CatalogTorrent.Columns().SpExpireAt: spExpireAt,
		dao.CatalogTorrent.Columns().UpdatedAt:  gtime.Now(),
	}).Update()
	return err
}

func (s *sCatalogTorrentDomain) GetTorrentFiles(ctx context.Context, torrentId uint64) ([]entity.CatalogTorrentFile, error) {
	var files []entity.CatalogTorrentFile
	err := dao.CatalogTorrentFile.Ctx(ctx).Where(dao.CatalogTorrentFile.Columns().TorrentId, torrentId).Scan(&files)
	return files, err
}

func (s *sCatalogTorrentDomain) GetTorrentsByIds(ctx context.Context, ids []uint64) ([]*entity.CatalogTorrent, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []*entity.CatalogTorrent
	err := dao.CatalogTorrent.Ctx(ctx).WhereIn("id", ids).Scan(&list)
	return list, err
}

func (s *sCatalogTorrentDomain) DeleteTorrent(ctx context.Context, id uint64) error {
	// 删除种子扩展元数据
	if _, err := dao.CatalogTorrentMeta.Ctx(ctx).Where(dao.CatalogTorrentMeta.Columns().TorrentId, id).Delete(); err != nil {
		return err
	}
	// 删除种子文件列表
	if _, err := dao.CatalogTorrentFile.Ctx(ctx).Where(dao.CatalogTorrentFile.Columns().TorrentId, id).Delete(); err != nil {
		return err
	}
	// 删除种子标签关联
	if _, err := dao.CatalogTorrentTag.Ctx(ctx).Where(dao.CatalogTorrentTag.Columns().TorrentId, id).Delete(); err != nil {
		return err
	}
	// 删除种子收藏
	if _, err := dao.CatalogTorrentBookmark.Ctx(ctx).Where(dao.CatalogTorrentBookmark.Columns().TorrentId, id).Delete(); err != nil {
		return err
	}
	// 删除种子点赞
	if _, err := dao.CatalogTorrentLike.Ctx(ctx).Where(dao.CatalogTorrentLike.Columns().TorrentId, id).Delete(); err != nil {
		return err
	}
	// 最后删除主表
	_, err := dao.CatalogTorrent.Ctx(ctx).WherePri(id).Delete()
	return err
}

func (s *sCatalogTorrentDomain) UpdateTorrentRewardStats(ctx context.Context, torrentId uint64, amount float64) error {
	_, err := dao.CatalogTorrent.Ctx(ctx).
		Where(dao.CatalogTorrent.Columns().Id, torrentId).
		Increment(dao.CatalogTorrent.Columns().RewardsCount, 1)
	if err != nil {
		return err
	}
	_, err = dao.CatalogTorrent.Ctx(ctx).
		Where(dao.CatalogTorrent.Columns().Id, torrentId).
		Increment(dao.CatalogTorrent.Columns().RewardsAmount, amount)
	return err
}

func (s *sCatalogTorrentDomain) CheckBookmarked(ctx context.Context, torrentId, userId uint64) (bool, error) {
	count, err := dao.CatalogTorrentBookmark.Ctx(ctx).Where(dao.CatalogTorrentBookmark.Columns().TorrentId, torrentId).Where(dao.CatalogTorrentBookmark.Columns().UserId, userId).Count()
	return count > 0, err
}

func (s *sCatalogTorrentDomain) CheckLiked(ctx context.Context, torrentId, userId uint64) (bool, error) {
	count, err := dao.CatalogTorrentLike.Ctx(ctx).Where(dao.CatalogTorrentLike.Columns().TorrentId, torrentId).Where(dao.CatalogTorrentLike.Columns().UserId, userId).Count()
	return count > 0, err
}

func (s *sCatalogTorrentDomain) UpdateTorrentPeerStats(ctx context.Context, torrentId uint64, seeders int, leechers int) error {
	if torrentId == 0 {
		return nil
	}
	columns := dao.CatalogTorrent.Columns()
	_, err := dao.CatalogTorrent.Ctx(ctx).
		Where(columns.Id, torrentId).
		Data(g.Map{
			columns.Seeders:  seeders,
			columns.Leechers: leechers,
		}).
		Update()
	return err
}

func (s *sCatalogTorrentDomain) IncrementTorrentStats(ctx context.Context, torrentId uint64, field string, amount float64) error {
	_, err := dao.CatalogTorrent.Ctx(ctx).Where(dao.CatalogTorrent.Columns().Id, torrentId).Increment(field, amount)
	return err
}

func (s *sCatalogTorrentDomain) GetTorrentsByHashes(ctx context.Context, hashes []string) ([]entity.CatalogTorrent, error) {
	var torrents []entity.CatalogTorrent
	err := dao.CatalogTorrent.Ctx(ctx).
		Fields(dao.CatalogTorrent.Columns().Id, dao.CatalogTorrent.Columns().Visible).
		Where(dao.CatalogTorrent.Columns().InfoHash, hashes).
		Scan(&torrents)
	return torrents, err
}

func (s *sCatalogTorrentDomain) QueryTorrentsByConditions(ctx context.Context, actor *model.Actor, keyword string, categoryIds []uint, page, size int) ([]entity.CatalogTorrent, int, error) {
	m := dao.CatalogTorrent.Ctx(ctx)
	m = s.ApplyTorrentVisibleScope(m, actor)

	columns := dao.CatalogTorrent.Columns()
	categoryIds = s.normalizeTorrentCategoryIds(categoryIds)
	if len(categoryIds) > 0 {
		m = m.WhereIn(columns.CategoryId, categoryIds)
	}

	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		m = m.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", columns.Name, columns.SubTitle), like, like)
	}

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var entities []entity.CatalogTorrent
	err = m.Page(page, size).
		OrderDesc(columns.IsPinned).
		OrderDesc(columns.PinWeight).
		OrderDesc(columns.CreatedAt).
		Scan(&entities)
	return entities, total, err
}

func (s *sCatalogTorrentDomain) normalizeTorrentCategoryIds(categoryIds []uint) []uint {
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

func (s *sCatalogTorrentDomain) CheckTorrentBookmarked(ctx context.Context, torrentId, userId uint64) (bool, error) {
	count, err := dao.CatalogTorrentBookmark.Ctx(ctx).Where(dao.CatalogTorrentBookmark.Columns().TorrentId, torrentId).Where(dao.CatalogTorrentBookmark.Columns().UserId, userId).Count()
	return count > 0, err
}

func (s *sCatalogTorrentDomain) CheckTorrentLiked(ctx context.Context, torrentId, userId uint64) (bool, error) {
	count, err := dao.CatalogTorrentLike.Ctx(ctx).Where(dao.CatalogTorrentLike.Columns().TorrentId, torrentId).Where(dao.CatalogTorrentLike.Columns().UserId, userId).Count()
	return count > 0, err
}

func (s *sCatalogTorrentDomain) IncrementTorrentRewardStats(ctx context.Context, torrentId uint64, amount float64) error {
	_, err := dao.CatalogTorrent.Ctx(ctx).
		Where(dao.CatalogTorrent.Columns().Id, torrentId).
		Increment(dao.CatalogTorrent.Columns().RewardsCount, 1)
	if err != nil {
		return err
	}
	_, err = dao.CatalogTorrent.Ctx(ctx).
		Where(dao.CatalogTorrent.Columns().Id, torrentId).
		Increment(dao.CatalogTorrent.Columns().RewardsAmount, amount)
	return err
}

func (s *sCatalogTorrentDomain) QueryActiveTorrentIds(ctx context.Context) ([]entity.CatalogTorrent, error) {
	var dbTorrents []entity.CatalogTorrent
	err := dao.CatalogTorrent.Ctx(ctx).
		Where("seeders > 0 OR leechers > 0").
		Fields("id").
		Scan(&dbTorrents)
	return dbTorrents, err
}

func (s *sCatalogTorrentDomain) ApplyTorrentVisibleScope(m *gdb.Model, actor *model.Actor) *gdb.Model {
	if actor != nil && actor.IsStaff {
		return m
	}
	if actor != nil {
		return m.Where("banned = 0 AND (visible = 1 OR owner_id = ?)", actor.Id)
	}
	return m.Where("visible = 1 AND banned = 0")
}

func (s *sCatalogTorrentDomain) CheckTorrentVisiblePolicy(ctx context.Context, actor *model.Actor, torrent *entity.CatalogTorrent) error {
	if actor != nil && actor.IsStaff {
		return nil
	}
	if torrent.Banned {
		return gerror.New(gi18n.T(ctx, "catalog.torrent.banned"))
	}
	if actor != nil && torrent.OwnerId == actor.Id {
		return nil
	}
	if !torrent.Visible {
		return gerror.New(gi18n.T(ctx, "catalog.torrent.not_visible"))
	}
	return nil
}

func (s *sCatalogTorrentDomain) CheckTorrentDownloadPolicy(ctx context.Context, actor *model.Actor, torrent *entity.CatalogTorrent) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	return s.CheckTorrentVisiblePolicy(ctx, actor, torrent)
}

func (s *sCatalogTorrentDomain) LoadVisibleTorrent(ctx context.Context, actor *model.Actor, id uint64) (*entity.CatalogTorrent, error) {
	torrent, err := s.GetTorrentById(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.CheckTorrentVisiblePolicy(ctx, actor, torrent); err != nil {
		return nil, err
	}
	return torrent, nil
}
