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

func (s *sCatalogTorrentDomain) GetTorrentByIdForUpdate(ctx context.Context, id uint64) (*entity.CatalogTorrent, error) {
	var torrent entity.CatalogTorrent
	err := dao.CatalogTorrent.Ctx(ctx).LockUpdate().Where(dao.CatalogTorrent.Columns().Id, id).Scan(&torrent)
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
	m = s.ApplyTorrentPublishedScope(m)
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

func (s *sCatalogTorrentDomain) AdminQueryReviewTorrents(ctx context.Context, options model.CatalogTorrentReviewListOptions) ([]entity.CatalogTorrent, int, error) {
	columns := dao.CatalogTorrent.Columns()
	m := dao.CatalogTorrent.Ctx(ctx)
	if options.Status >= 0 {
		m = m.Where(columns.Status, options.Status)
	}
	if options.CategoryId > 0 {
		m = m.Where(columns.CategoryId, options.CategoryId)
	}
	if keyword := strings.TrimSpace(options.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		m = m.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", columns.Name, columns.SubTitle), like, like)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var torrents []entity.CatalogTorrent
	err = m.OrderDesc(columns.SubmittedAt).
		OrderDesc(columns.Id).
		Page(options.Page, options.Size).
		Scan(&torrents)
	return torrents, total, err
}

func (s *sCatalogTorrentDomain) QueryUserTorrents(ctx context.Context, userId uint64, options model.CatalogUserTorrentListOptions) ([]entity.CatalogTorrent, int, error) {
	columns := dao.CatalogTorrent.Columns()
	m := dao.CatalogTorrent.Ctx(ctx).Where(columns.OwnerId, userId)
	if options.Status >= 0 {
		m = m.Where(columns.Status, options.Status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var torrents []entity.CatalogTorrent
	err = m.OrderDesc(columns.UpdatedAt).
		OrderDesc(columns.Id).
		Page(options.Page, options.Size).
		Scan(&torrents)
	return torrents, total, err
}

func (s *sCatalogTorrentDomain) AdminApproveTorrent(ctx context.Context, id uint64, reviewedBy uint64, comment string, publishedAt *gtime.Time, spState int, spExpireAt *gtime.Time) (bool, error) {
	columns := dao.CatalogTorrent.Columns()
	result, err := dao.CatalogTorrent.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.Status, consts.CatalogTorrentStatusPending).
		Data(g.Map{
			columns.Status:        consts.CatalogTorrentStatusPublished,
			columns.PublishedAt:   publishedAt,
			columns.ReviewedBy:    reviewedBy,
			columns.ReviewedAt:    publishedAt,
			columns.ReviewComment: strings.TrimSpace(comment),
			columns.SpState:       spState,
			columns.SpExpireAt:    spExpireAt,
			columns.UpdatedAt:     publishedAt,
		}).Update()
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	return affected > 0, err
}

func (s *sCatalogTorrentDomain) AdminRejectTorrent(ctx context.Context, id uint64, reviewedBy uint64, comment string, reviewedAt *gtime.Time) (bool, error) {
	columns := dao.CatalogTorrent.Columns()
	result, err := dao.CatalogTorrent.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.Status, consts.CatalogTorrentStatusPending).
		Data(g.Map{
			columns.Status:        consts.CatalogTorrentStatusRejected,
			columns.ReviewedBy:    reviewedBy,
			columns.ReviewedAt:    reviewedAt,
			columns.ReviewComment: strings.TrimSpace(comment),
			columns.UpdatedAt:     reviewedAt,
		}).Update()
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	return affected > 0, err
}

func (s *sCatalogTorrentDomain) ResubmitTorrent(ctx context.Context, id uint64, status int, submittedAt, publishedAt *gtime.Time, spState int, spExpireAt *gtime.Time) (bool, error) {
	columns := dao.CatalogTorrent.Columns()
	result, err := dao.CatalogTorrent.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.Status, consts.CatalogTorrentStatusRejected).
		Data(g.Map{
			columns.Status:        status,
			columns.SubmittedAt:   submittedAt,
			columns.PublishedAt:   publishedAt,
			columns.ReviewedBy:    0,
			columns.ReviewedAt:    nil,
			columns.ReviewComment: "",
			columns.SpState:       spState,
			columns.SpExpireAt:    spExpireAt,
			columns.UpdatedAt:     gtime.Now(),
		}).Update()
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	return affected > 0, err
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
		Fields(dao.CatalogTorrent.Columns().Id, dao.CatalogTorrent.Columns().Status, dao.CatalogTorrent.Columns().Banned).
		Where(dao.CatalogTorrent.Columns().InfoHash, hashes).
		Scan(&torrents)
	return torrents, err
}

func (s *sCatalogTorrentDomain) QueryTorrents(ctx context.Context, actor *model.Actor, options model.CatalogTorrentListOptions) ([]entity.CatalogTorrent, int, error) {
	options = options.Normalized()
	m := s.buildTorrentQuery(ctx, actor, options)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var entities []entity.CatalogTorrent
	err = s.applyTorrentOrder(m, options.Sort, true).
		Page(options.Page, options.Size).
		Scan(&entities)
	return entities, total, err
}

func (s *sCatalogTorrentDomain) QueryRssTorrents(ctx context.Context, actor *model.Actor, options model.CatalogTorrentListOptions) ([]entity.CatalogTorrent, error) {
	options = options.Normalized()
	columns := dao.CatalogTorrent.Columns()
	var entities []entity.CatalogTorrent
	err := s.buildTorrentQuery(ctx, actor, options).
		OrderDesc(columns.PublishedAt).
		OrderDesc(columns.Id).
		Limit(options.Size).
		Scan(&entities)
	return entities, err
}

func (s *sCatalogTorrentDomain) buildTorrentQuery(ctx context.Context, actor *model.Actor, options model.CatalogTorrentListOptions) *gdb.Model {
	options = options.Normalized()
	m := s.ApplyTorrentPublishedScope(dao.CatalogTorrent.Ctx(ctx))
	columns := dao.CatalogTorrent.Columns()
	if len(options.CategoryIds) > 0 {
		m = m.WhereIn(columns.CategoryId, options.CategoryIds)
	}
	m = s.applyTorrentTagFilter(m, options.TagGroups)
	m = s.applyTorrentMetadataFilter(m, options.Metadata)

	if options.Keyword != "" {
		like := "%" + options.Keyword + "%"
		m = m.Where(fmt.Sprintf("(%s LIKE ? OR %s LIKE ?)", columns.Name, columns.SubTitle), like, like)
	}
	if options.FeaturedOnly {
		m = m.Where(columns.IsFeatured, true)
	}
	if options.MinSize > 0 {
		m = m.WhereGTE(columns.Size, options.MinSize)
	}
	if options.MaxSize > 0 {
		m = m.WhereLTE(columns.Size, options.MaxSize)
	}
	if options.PublishedWithinDays > 0 {
		m = m.WhereGTE(columns.PublishedAt, gtime.Now().Add(-time.Duration(options.PublishedWithinDays)*24*time.Hour))
	}

	switch options.SeedStatus {
	case consts.CatalogTorrentSeedStatusSeeded:
		m = m.WhereGT(columns.Seeders, 0)
	case consts.CatalogTorrentSeedStatusUnseeded:
		m = m.Where(columns.Seeders, 0)
	}

	return s.applyTorrentPromotionFilter(ctx, m, options.Promotion)
}

func (s *sCatalogTorrentDomain) applyTorrentTagFilter(m *gdb.Model, groups [][]uint) *gdb.Model {
	torrentTable := dao.CatalogTorrent.Table()
	torrentColumns := dao.CatalogTorrent.Columns()
	relationTable := dao.CatalogTorrentTag.Table()
	relationColumns := dao.CatalogTorrentTag.Columns()
	for _, tagIds := range groups {
		if len(tagIds) == 0 {
			continue
		}
		query := fmt.Sprintf(
			"EXISTS (SELECT 1 FROM %s WHERE %s.%s = %s.%s AND %s.%s IN (?))",
			relationTable,
			relationTable, relationColumns.TorrentId,
			torrentTable, torrentColumns.Id,
			relationTable, relationColumns.TagId,
		)
		m = m.Where(query, tagIds)
	}
	return m
}

func (s *sCatalogTorrentDomain) applyTorrentMetadataFilter(m *gdb.Model, filter model.CatalogTorrentMetadataFilter) *gdb.Model {
	filter = filter.Normalized()
	if filter.IsEmpty() {
		return m
	}

	torrentTable := dao.CatalogTorrent.Table()
	metaTable := dao.CatalogTorrentMeta.Table()
	torrentColumns := dao.CatalogTorrent.Columns()
	metaColumns := dao.CatalogTorrentMeta.Columns()
	conditions := []string{fmt.Sprintf("%s.%s = %s.%s", metaTable, metaColumns.TorrentId, torrentTable, torrentColumns.Id)}
	args := make([]any, 0, 5)
	if filter.ImdbId != "" {
		conditions = append(conditions, fmt.Sprintf("%s.%s = ?", metaTable, metaColumns.ImdbId))
		args = append(args, filter.ImdbId)
	}
	if filter.DoubanId != "" {
		conditions = append(conditions, fmt.Sprintf("%s.%s = ?", metaTable, metaColumns.DoubanId))
		args = append(args, filter.DoubanId)
	}
	if filter.BangumiId != "" {
		conditions = append(conditions, fmt.Sprintf("%s.%s = ?", metaTable, metaColumns.BangumiId))
		args = append(args, filter.BangumiId)
	}
	if filter.TmdbId != "" {
		conditions = append(conditions, fmt.Sprintf("%s.%s = ?", metaTable, metaColumns.TmdbId))
		args = append(args, filter.TmdbId)
		if filter.TmdbType != "" {
			conditions = append(conditions, fmt.Sprintf("%s.%s = ?", metaTable, metaColumns.TmdbType))
			args = append(args, filter.TmdbType)
		}
	}

	query := fmt.Sprintf("EXISTS (SELECT 1 FROM %s WHERE %s)", metaTable, strings.Join(conditions, " AND "))
	return m.Where(query, args...)
}

func (s *sCatalogTorrentDomain) applyTorrentPromotionFilter(ctx context.Context, m *gdb.Model, filter string) *gdb.Model {
	if filter == consts.CatalogTorrentPromotionFilterAll {
		return m
	}

	globalPromotion := model.ResolveCatalogTorrentPromotion(
		consts.ResourceTorrentSpNormal,
		nil,
		s.getGlobalPromotionConfig(ctx),
		nil,
	)
	if globalPromotion.SpState != consts.ResourceTorrentSpNormal {
		if filter == consts.CatalogTorrentPromotionFilterPromoted || model.CatalogTorrentPromotionStateToSp(filter) == globalPromotion.SpState {
			return m
		}
		return m.Where("1 = 0")
	}

	columns := dao.CatalogTorrent.Columns()
	now := gtime.Now()
	activeCondition := fmt.Sprintf("(%s IS NULL OR %s > ?)", columns.SpExpireAt, columns.SpExpireAt)
	expiredCondition := fmt.Sprintf("(%s IS NOT NULL AND %s <= ?)", columns.SpExpireAt, columns.SpExpireAt)
	switch filter {
	case consts.CatalogTorrentPromotionFilterPromoted:
		return m.Where(fmt.Sprintf("%s != ? AND %s", columns.SpState, activeCondition), consts.ResourceTorrentSpNormal, now)
	case consts.ResourceTorrentPromotionStateNormal:
		return m.Where(fmt.Sprintf("(%s = ? OR %s)", columns.SpState, expiredCondition), consts.ResourceTorrentSpNormal, now)
	default:
		spState := model.CatalogTorrentPromotionStateToSp(filter)
		if spState == consts.ResourceTorrentSpNormal {
			return m.Where("1 = 0")
		}
		return m.Where(fmt.Sprintf("%s = ? AND %s", columns.SpState, activeCondition), spState, now)
	}
}

func (s *sCatalogTorrentDomain) applyTorrentOrder(m *gdb.Model, sort string, pinnedFirst bool) *gdb.Model {
	columns := dao.CatalogTorrent.Columns()
	if pinnedFirst {
		m = m.OrderDesc(columns.IsPinned).OrderDesc(columns.PinWeight)
	}

	switch sort {
	case consts.CatalogTorrentSortOldest:
		return m.OrderAsc(columns.PublishedAt).OrderAsc(columns.Id)
	case consts.CatalogTorrentSortSeeders:
		return m.OrderDesc(columns.Seeders).OrderDesc(columns.PublishedAt).OrderDesc(columns.Id)
	case consts.CatalogTorrentSortLeechers:
		return m.OrderDesc(columns.Leechers).OrderDesc(columns.PublishedAt).OrderDesc(columns.Id)
	case consts.CatalogTorrentSortComplete:
		return m.OrderDesc(columns.TimesCompleted).OrderDesc(columns.PublishedAt).OrderDesc(columns.Id)
	case consts.CatalogTorrentSortSizeAsc:
		return m.OrderAsc(columns.Size).OrderDesc(columns.PublishedAt).OrderDesc(columns.Id)
	case consts.CatalogTorrentSortSizeDesc:
		return m.OrderDesc(columns.Size).OrderDesc(columns.PublishedAt).OrderDesc(columns.Id)
	default:
		return m.OrderDesc(columns.PublishedAt).OrderDesc(columns.Id)
	}
}

func (s *sCatalogTorrentDomain) QueryHotVisibleTorrents(ctx context.Context, size int) ([]entity.CatalogTorrent, error) {
	if size <= 0 {
		size = 5
	}
	if size > 50 {
		size = 50
	}

	columns := dao.CatalogTorrent.Columns()
	order := fmt.Sprintf(
		"(%s * 3 + %s * 2 + %s + %s * 2) DESC, %s DESC",
		columns.Seeders,
		columns.Leechers,
		columns.TimesCompleted,
		columns.LikeCount,
		columns.PublishedAt,
	)

	var entities []entity.CatalogTorrent
	err := dao.CatalogTorrent.Ctx(ctx).
		Where(g.Map{
			columns.Status: consts.CatalogTorrentStatusPublished,
			columns.Banned: false,
		}).
		Order(order).
		Limit(size).
		Scan(&entities)
	return entities, err
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

func (s *sCatalogTorrentDomain) ApplyTorrentPublishedScope(m *gdb.Model) *gdb.Model {
	columns := dao.CatalogTorrent.Columns()
	return m.Where(columns.Status, consts.CatalogTorrentStatusPublished).Where(columns.Banned, false)
}

func (s *sCatalogTorrentDomain) CheckTorrentViewPolicy(ctx context.Context, actor *model.Actor, torrent *entity.CatalogTorrent) error {
	if actor != nil && actor.IsStaff {
		return nil
	}
	if torrent.Banned {
		return gerror.New(gi18n.T(ctx, "catalog.torrent.banned"))
	}
	if torrent.Status == consts.CatalogTorrentStatusPublished {
		return nil
	}
	if actor != nil && torrent.OwnerId == actor.Id {
		return nil
	}
	return gerror.New(gi18n.T(ctx, "catalog.torrent.not_visible"))
}

func (s *sCatalogTorrentDomain) CheckTorrentDownloadPolicy(ctx context.Context, actor *model.Actor, torrent *entity.CatalogTorrent) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	if torrent.Banned {
		return gerror.New(gi18n.T(ctx, "catalog.torrent.banned"))
	}
	if torrent.Status == consts.CatalogTorrentStatusPublished || actor.IsStaff || torrent.OwnerId == actor.Id {
		return nil
	}
	return gerror.New(gi18n.T(ctx, "catalog.torrent.not_visible"))
}

func (s *sCatalogTorrentDomain) CheckTorrentAnnouncePolicy(ctx context.Context, actor *model.Actor, torrent *entity.CatalogTorrent) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	if torrent.Banned || torrent.Status == consts.CatalogTorrentStatusRejected {
		return gerror.New(gi18n.T(ctx, "catalog.torrent.banned"))
	}
	if torrent.Status == consts.CatalogTorrentStatusPublished || actor.IsStaff || torrent.OwnerId == actor.Id {
		return nil
	}
	return gerror.New(gi18n.T(ctx, "catalog.torrent.not_visible"))
}

func (s *sCatalogTorrentDomain) LoadViewableTorrent(ctx context.Context, actor *model.Actor, id uint64) (*entity.CatalogTorrent, error) {
	torrent, err := s.GetTorrentById(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.CheckTorrentViewPolicy(ctx, actor, torrent); err != nil {
		return nil, err
	}
	return torrent, nil
}
