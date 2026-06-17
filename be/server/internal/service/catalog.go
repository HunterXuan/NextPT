// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type (
	ICatalogCategoryDomain interface {
		ListCategories(ctx context.Context) ([]entity.CatalogCategory, error)
		ListTagGroups(ctx context.Context) ([]entity.CatalogTagGroup, []entity.CatalogTag, error)
	}
	ICatalogCategoryUsecase interface {
		ListCategories(ctx context.Context, actor *model.Actor, in catalogin.CategoryListInp) (*catalogout.CategoryListOut, error)
		ListTagGroups(ctx context.Context, actor *model.Actor, in catalogin.TagGroupListInp) (*catalogout.TagGroupListOut, error)
	}
	ICatalogCommentDomain interface {
		CreateComment(ctx context.Context, targetType string, targetId uint64, userId uint64, content string) (uint64, error)
		GetCommentById(ctx context.Context, id uint64) (*entity.CatalogComment, error)
		QueryCommentsByTarget(ctx context.Context, targetType string, targetId uint64, page int, size int) ([]entity.CatalogComment, int, error)
		ToggleLike(ctx context.Context, userId uint64, commentId uint64) (bool, error)
		GetCommentLikesByUser(ctx context.Context, userId uint64, commentIds []uint64) ([]entity.CatalogCommentLike, error)
		InsertCommentReward(ctx context.Context, userId uint64, commentId uint64, amount float64) error
		IncrementCommentRewardStats(ctx context.Context, commentId uint64) error
		DeleteCommentsByTarget(ctx context.Context, targetType string, targetId uint64) error
	}
	ICatalogCommentUsecase interface {
		Create(ctx context.Context, actor *model.Actor, in catalogin.CommentCreateInp) (*catalogout.CommentCreateOut, error)
		List(ctx context.Context, actor *model.Actor, in catalogin.CommentListInp) (*catalogout.CommentListOut, error)
		ToggleLike(ctx context.Context, actor *model.Actor, in catalogin.CommentToggleLikeInp) (*catalogout.CommentToggleLikeOut, error)
		Reward(ctx context.Context, actor *model.Actor, in catalogin.CommentRewardInp) error
		Report(ctx context.Context, actor *model.Actor, in catalogin.CommentReportInp) error
	}
	ICatalogSubtitleDomain interface {
		InsertSubtitle(ctx context.Context, torrentId uint64, userId uint64, fileName string, ext string, size int, language string) (uint64, error)
		UpdateSubtitleStoragePath(ctx context.Context, id uint64, path string) error
		GetSubtitleById(ctx context.Context, id uint64) (*entity.CatalogSubtitle, error)
		IncrementDownloadCount(ctx context.Context, id uint64) error
		QuerySubtitles(ctx context.Context, torrentId uint64, page int, size int) ([]entity.CatalogSubtitle, int, error)
		UpdateSubtitle(ctx context.Context, id uint64, language string, userId uint64, isAdmin bool) error
		GetSubtitlesByTorrentId(ctx context.Context, torrentId uint64) ([]*entity.CatalogSubtitle, error)
		DeleteSubtitlesByTorrentId(ctx context.Context, torrentId uint64) error
		DeleteSubtitle(ctx context.Context, id uint64) error
	}
	ICatalogSubtitleUsecase interface {
		List(ctx context.Context, actor *model.Actor, in catalogin.SubtitleListInp) (*catalogout.SubtitleListOut, error)
		ListByTorrent(ctx context.Context, actor *model.Actor, in catalogin.TorrentSubtitleListInp) (*catalogout.SubtitleListOut, error)
		Upload(ctx context.Context, actor *model.Actor, in catalogin.SubtitleUploadInp) (uint64, error)
		Download(ctx context.Context, actor *model.Actor, in catalogin.SubtitleDownloadInp) (*catalogout.SubtitleDownloadOut, error)
		Update(ctx context.Context, actor *model.Actor, in catalogin.SubtitleUpdateInp) error
		Report(ctx context.Context, actor *model.Actor, in catalogin.SubtitleReportInp) error
	}
	ICatalogTorrentDomain interface {
		GetTorrentById(ctx context.Context, id uint64) (*entity.CatalogTorrent, error)
		GetTorrentByInfoHash(ctx context.Context, infoHash string) (*entity.CatalogTorrent, error)
		CheckCategoryExists(ctx context.Context, categoryId uint) error
		CheckInfoHashExists(ctx context.Context, infoHashBytes []byte) error
		SaveTorrent(ctx context.Context, torrentInsert *entity.CatalogTorrent, filesToInsert []entity.CatalogTorrentFile) (uint64, error)
		Bookmark(ctx context.Context, userId uint64, torrentId uint64) error
		Unbookmark(ctx context.Context, userId uint64, torrentId uint64) error
		QueryBookmarkedTorrents(ctx context.Context, actor *model.Actor, page int, size int) ([]entity.CatalogTorrent, int, error)
		ToggleLike(ctx context.Context, torrentId uint64, userId uint64) (bool, error)
		QueryTorrentLikes(ctx context.Context, torrentId uint64, page int, size int) ([]entity.CatalogTorrentLike, int, error)
		UpdateTorrent(ctx context.Context, id uint64, name string, subTitle string, categoryId uint, description string, anonymous *bool) error
		GetTorrentFiles(ctx context.Context, torrentId uint64) ([]entity.CatalogTorrentFile, error)
		GetTorrentsByIds(ctx context.Context, ids []uint64) ([]*entity.CatalogTorrent, error)
		DeleteTorrent(ctx context.Context, id uint64) error
		QueryTorrents(ctx context.Context, actor *model.Actor, categoryId uint, torrentType *int, page int, size int) ([]entity.CatalogTorrent, int, error)
		InsertTorrentReward(ctx context.Context, reward *entity.CatalogTorrentReward) error
		UpdateTorrentRewardStats(ctx context.Context, torrentId uint64, amount float64) error
		CheckBookmarked(ctx context.Context, torrentId uint64, userId uint64) (bool, error)
		CheckLiked(ctx context.Context, torrentId uint64, userId uint64) (bool, error)
		BatchUpdateTorrents(ctx context.Context, updates []g.Map) error
		IncrementTorrentStats(ctx context.Context, torrentId uint64, field string, amount float64) error
		UpdateTorrentStatsFromSync(ctx context.Context, updates g.Map) error
		GetTorrentsByHashes(ctx context.Context, hashes []string) ([]entity.CatalogTorrent, error)
		QueryTorrentsByConditions(ctx context.Context, actor *model.Actor, categoryId uint, torrentType *int, page int, size int) ([]entity.CatalogTorrent, int, error)
		CheckTorrentBookmarked(ctx context.Context, torrentId uint64, userId uint64) (bool, error)
		CheckTorrentLiked(ctx context.Context, torrentId uint64, userId uint64) (bool, error)
		IncrementTorrentRewardStats(ctx context.Context, torrentId uint64, amount float64) error
		QueryActiveTorrentIds(ctx context.Context) ([]entity.CatalogTorrent, error)
		QueryTorrentRewards(ctx context.Context, torrentId uint64, page int, size int) ([]entity.CatalogTorrentReward, int, error)
		ApplyTorrentVisibleScope(m *gdb.Model, actor *model.Actor) *gdb.Model
		CheckTorrentVisiblePolicy(ctx context.Context, actor *model.Actor, torrent *entity.CatalogTorrent) error
		CheckTorrentDownloadPolicy(ctx context.Context, actor *model.Actor, torrent *entity.CatalogTorrent) error
		LoadVisibleTorrent(ctx context.Context, actor *model.Actor, id uint64) (*entity.CatalogTorrent, error)
	}
	ICatalogTorrentUsecase interface {
		// List 获取种子分页列表
		List(ctx context.Context, actor *model.Actor, in catalogin.TorrentListInp) (*catalogout.TorrentListOut, error)
		// Download 获取用户专属的种子文件内容
		Download(ctx context.Context, actor *model.Actor, in catalogin.TorrentDownloadInp) (*catalogout.TorrentDownloadOut, error)
		// Upload 解析用户上传的种子文件，处理为私有种子，并保存到数据库和 S3
		Upload(ctx context.Context, actor *model.Actor, in catalogin.TorrentUploadInp) (*catalogout.TorrentUploadOut, error)
		// HardDeleteTorrent 执行大统一的种子硬删除，包含所有关联业务数据和 S3 文件
		HardDeleteTorrent(ctx context.Context, torrentId uint64) error
		// Reward 赞赏种子
		Reward(ctx context.Context, actor *model.Actor, in catalogin.TorrentRewardInp) (*catalogout.TorrentRewardOut, error)
		// RewardList 获取种子赞赏列表
		RewardList(ctx context.Context, actor *model.Actor, in catalogin.TorrentRewardListInp) (*catalogout.TorrentRewardListOut, error)
		Bookmark(ctx context.Context, actor *model.Actor, in catalogin.TorrentBookmarkInp) error
		Unbookmark(ctx context.Context, actor *model.Actor, in catalogin.TorrentUnbookmarkInp) error
		ListBookmarkedTorrents(ctx context.Context, actor *model.Actor, in catalogin.TorrentBookmarkListInp) (*catalogout.TorrentBookmarkListOut, error)
		GetTorrent(ctx context.Context, actor *model.Actor, in catalogin.TorrentGetInp) (*catalogout.TorrentDetailOut, error)
		ToggleLike(ctx context.Context, actor *model.Actor, in catalogin.TorrentToggleLikeInp) (*catalogout.TorrentToggleLikeOut, error)
		ListLikes(ctx context.Context, actor *model.Actor, in catalogin.TorrentLikeListInp) (*catalogout.TorrentLikeListOut, error)
		Update(ctx context.Context, actor *model.Actor, in catalogin.TorrentUpdateInp) (*catalogout.TorrentUpdateOut, error)
		ListFiles(ctx context.Context, actor *model.Actor, in catalogin.TorrentFileListInp) (*catalogout.TorrentFileListOut, error)
		ListPeers(ctx context.Context, actor *model.Actor, in catalogin.TorrentPeerListInp) (*catalogout.TorrentPeerListOut, error)
		Report(ctx context.Context, actor *model.Actor, in catalogin.TorrentReportInp) (*catalogout.TorrentReportOut, error)
	}
)

var (
	localCatalogCategoryDomain  ICatalogCategoryDomain
	localCatalogCategoryUsecase ICatalogCategoryUsecase
	localCatalogCommentDomain   ICatalogCommentDomain
	localCatalogCommentUsecase  ICatalogCommentUsecase
	localCatalogSubtitleDomain  ICatalogSubtitleDomain
	localCatalogSubtitleUsecase ICatalogSubtitleUsecase
	localCatalogTorrentDomain   ICatalogTorrentDomain
	localCatalogTorrentUsecase  ICatalogTorrentUsecase
)

func CatalogCategoryDomain() ICatalogCategoryDomain {
	if localCatalogCategoryDomain == nil {
		panic("implement not found for interface ICatalogCategoryDomain, forgot register?")
	}
	return localCatalogCategoryDomain
}

func RegisterCatalogCategoryDomain(i ICatalogCategoryDomain) {
	localCatalogCategoryDomain = i
}

func CatalogCategoryUsecase() ICatalogCategoryUsecase {
	if localCatalogCategoryUsecase == nil {
		panic("implement not found for interface ICatalogCategoryUsecase, forgot register?")
	}
	return localCatalogCategoryUsecase
}

func RegisterCatalogCategoryUsecase(i ICatalogCategoryUsecase) {
	localCatalogCategoryUsecase = i
}

func CatalogCommentDomain() ICatalogCommentDomain {
	if localCatalogCommentDomain == nil {
		panic("implement not found for interface ICatalogCommentDomain, forgot register?")
	}
	return localCatalogCommentDomain
}

func RegisterCatalogCommentDomain(i ICatalogCommentDomain) {
	localCatalogCommentDomain = i
}

func CatalogCommentUsecase() ICatalogCommentUsecase {
	if localCatalogCommentUsecase == nil {
		panic("implement not found for interface ICatalogCommentUsecase, forgot register?")
	}
	return localCatalogCommentUsecase
}

func RegisterCatalogCommentUsecase(i ICatalogCommentUsecase) {
	localCatalogCommentUsecase = i
}

func CatalogSubtitleDomain() ICatalogSubtitleDomain {
	if localCatalogSubtitleDomain == nil {
		panic("implement not found for interface ICatalogSubtitleDomain, forgot register?")
	}
	return localCatalogSubtitleDomain
}

func RegisterCatalogSubtitleDomain(i ICatalogSubtitleDomain) {
	localCatalogSubtitleDomain = i
}

func CatalogSubtitleUsecase() ICatalogSubtitleUsecase {
	if localCatalogSubtitleUsecase == nil {
		panic("implement not found for interface ICatalogSubtitleUsecase, forgot register?")
	}
	return localCatalogSubtitleUsecase
}

func RegisterCatalogSubtitleUsecase(i ICatalogSubtitleUsecase) {
	localCatalogSubtitleUsecase = i
}

func CatalogTorrentDomain() ICatalogTorrentDomain {
	if localCatalogTorrentDomain == nil {
		panic("implement not found for interface ICatalogTorrentDomain, forgot register?")
	}
	return localCatalogTorrentDomain
}

func RegisterCatalogTorrentDomain(i ICatalogTorrentDomain) {
	localCatalogTorrentDomain = i
}

func CatalogTorrentUsecase() ICatalogTorrentUsecase {
	if localCatalogTorrentUsecase == nil {
		panic("implement not found for interface ICatalogTorrentUsecase, forgot register?")
	}
	return localCatalogTorrentUsecase
}

func RegisterCatalogTorrentUsecase(i ICatalogTorrentUsecase) {
	localCatalogTorrentUsecase = i
}
