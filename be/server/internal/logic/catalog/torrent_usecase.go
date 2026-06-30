package catalog

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/catalogin"
	"server/internal/model/in/modin"
	"server/internal/model/out/catalogout"
	"server/internal/service"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
)

type sCatalogTorrentUsecase struct{}

func init() {
	service.RegisterCatalogTorrentUsecase(NewCatalogTorrentUsecase())
}

func NewCatalogTorrentUsecase() *sCatalogTorrentUsecase {
	return &sCatalogTorrentUsecase{}
}

// List 获取种子分页列表
func (s *sCatalogTorrentUsecase) List(ctx context.Context, actor *model.Actor, in catalogin.TorrentListInp) (*catalogout.TorrentListOut, error) {
	entities, total, err := service.CatalogTorrentDomain().QueryTorrentsByConditions(ctx, actor, in.Keyword, in.CategoryIds, in.Page, in.Size)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.query_failed"))
	}

	list := s.formatTorrentListItems(ctx, actor, entities)

	return &catalogout.TorrentListOut{
		List:  list,
		Total: total,
	}, nil
}

// Download 获取用户专属的种子文件内容
func (s *sCatalogTorrentUsecase) Download(ctx context.Context, actor *model.Actor, in catalogin.TorrentDownloadInp) (*catalogout.TorrentDownloadOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	// 获取种子记录并校验可见性
	torrent, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return nil, err
	}

	if err := service.CatalogTorrentDomain().CheckTorrentDownloadPolicy(ctx, actor, torrent); err != nil {
		return nil, err
	}

	// 获取用户的 Passkey (用于注入 tracker)
	var user entity.IamUser
	userPtr, err := service.IamUserDomain().GetUserById(ctx, actor.Id)
	if err == nil && userPtr != nil {
		user = *userPtr
	}
	if err != nil || user.Passkey == "" {
		return nil, gerror.New(gi18n.T(ctx, "catalog.torrent.passkey_not_found"))
	}

	s3Key := fmt.Sprintf("torrents/%d.torrent", torrent.Id)

	// 获取本地存储/缓存路径
	targetPath := service.SysStorage().GetLocalPath(ctx, s3Key)

	// 如果本地不存在，则执行下载动作（对于 local 驱动，文件必定存在，天然跳过下载）
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		err := service.SysStorage().Download(ctx, s3Key, targetPath)
		if err != nil {
			g.Log().Errorf(ctx, "failed to download torrent from storage: %+v", err)
			return nil, gerror.New(gi18n.T(ctx, "catalog.torrent.load_failed"))
		}
	}

	// 读取文件并解析
	fileBytes, err := os.ReadFile(targetPath)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.read_cache_failed"))
	}

	mi, err := metainfo.Load(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.invalid_cache"))
	}

	// 注入 Announce (带 passkey)
	trackerUrl := s.getCatalogConfigCache(ctx, consts.SiteConfigTrackerUrl).String()
	mi.Announce = fmt.Sprintf("%s?passkey=%s", trackerUrl, user.Passkey)
	mi.AnnounceList = nil // 确保没有其他的 fallback trackers

	var outBuf bytes.Buffer
	if err := mi.Write(&outBuf); err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.generate_specific_failed"))
	}

	sourceStr := s.getCatalogConfigCache(ctx, consts.SiteConfigCatalogTorrentSource).String()
	fileName := fmt.Sprintf("[%s]%s.torrent", sourceStr, torrent.Name)
	return &catalogout.TorrentDownloadOut{
		Bytes:    outBuf.Bytes(),
		FileName: fileName,
	}, nil
}

// Upload 解析用户上传的种子文件，处理为私有种子，并保存到数据库和 S3
func (s *sCatalogTorrentUsecase) Upload(ctx context.Context, actor *model.Actor, in catalogin.TorrentUploadInp) (*catalogout.TorrentUploadOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	if in.File == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.torrent.file_req"))
	}

	category, err := service.CatalogCategoryDomain().GetCategoryById(ctx, in.CategoryId)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.category.invalid"))
	}

	mi, info, err := s.parseAndModifyTorrent(ctx, in)
	if err != nil {
		return nil, err
	}

	infoHashBytes, infoHashHex, err := s.calculateAndCheckInfoHash(ctx, mi)
	if err != nil {
		return nil, err
	}

	totalSize, fileCount := s.extractTorrentMetadata(info)
	releaseData, err := s.prepareUploadReleaseData(ctx, category, info.Name, in)
	if err != nil {
		return nil, err
	}

	var finalTorrentBuf bytes.Buffer
	if err := mi.Write(&finalTorrentBuf); err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.generate_final_failed"))
	}

	var torrentId uint64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		tid, err := s.saveTorrentToDB(ctx, actor, in, infoHashBytes, releaseData, totalSize, fileCount, info)
		if err != nil {
			return err
		}
		torrentId = tid
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 拿到数据库 ID 后，在事务外层执行对象存储上传，避免持有长事务锁
	s3Key := fmt.Sprintf("torrents/%d.torrent", torrentId)
	if err := service.SysStorage().Upload(ctx, s3Key, finalTorrentBuf.Bytes(), "application/x-bittorrent"); err != nil {
		_ = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			return service.CatalogTorrentDomain().DeleteTorrent(ctx, torrentId)
		})
		s.invalidateTorrentInfoHashCache(ctx, infoHashBytes)
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.upload_storage_failed"))
	}

	return &catalogout.TorrentUploadOut{
		TorrentId: torrentId,
		InfoHash:  infoHashHex,
	}, nil
}

// HardDeleteTorrent 执行大统一的种子硬删除，包含所有关联业务数据和 S3 文件
func (s *sCatalogTorrentUsecase) HardDeleteTorrent(ctx context.Context, torrentId uint64) error {
	torrent, err := service.CatalogTorrentDomain().GetTorrentById(ctx, torrentId)
	if err != nil {
		return err
	}

	// 1. 搜集需要删除的文件路径（字幕）
	subtitles, _ := service.CatalogSubtitleDomain().GetSubtitlesByTorrentId(ctx, torrentId)
	var subtitlePaths []string
	for _, sub := range subtitles {
		if sub.StoragePath != "" {
			subtitlePaths = append(subtitlePaths, sub.StoragePath)
		}
	}

	// 2. 执行跨域的 DB 大事务硬删除
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := service.CatalogTorrentDomain().DeleteTorrent(ctx, torrentId); err != nil {
			return err
		}
		if err := service.CatalogSubtitleDomain().DeleteSubtitlesByTorrentId(ctx, torrentId); err != nil {
			return err
		}
		if err := service.TrackerPeerDomain().DeletePeersByTorrentId(ctx, torrentId); err != nil {
			return err
		}
		if err := service.AccountingSnatchDomain().DeleteSnatchesByTorrentId(ctx, torrentId); err != nil {
			return err
		}
		if err := service.ModCheaterDomain().DeleteCheaterLogsByTorrentId(ctx, torrentId); err != nil {
			return err
		}
		if err := service.CatalogCommentDomain().DeleteCommentsByTarget(ctx, "torrent", torrentId); err != nil {
			return err
		}
		if err := service.ModReportDomain().DeleteReportsByTarget(ctx, "torrent", torrentId); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	s.invalidateTorrentInfoHashCache(ctx, torrent.InfoHash)

	// 3. 数据库删除成功后，清理 S3 和本地的文件
	// 删种子文件
	torrentS3Key := fmt.Sprintf("torrents/%d.torrent", torrentId)
	_ = service.SysStorage().Delete(ctx, torrentS3Key)

	// 删字幕文件
	for _, path := range subtitlePaths {
		_ = service.SysStorage().Delete(ctx, path)
	}

	return nil
}

func (s *sCatalogTorrentUsecase) invalidateTorrentInfoHashCache(ctx context.Context, infoHash []byte) {
	if len(infoHash) == 0 {
		return
	}

	cacheKey := service.SysCache().KeyCatalogTorrentInfoHash(ctx, hex.EncodeToString(infoHash))
	_, _ = gcache.Remove(ctx, cacheKey)
	_ = service.SysCache().PublishInvalidate(ctx, cacheKey)
}

func (s *sCatalogTorrentUsecase) validateCategory(ctx context.Context, categoryId uint) error {
	err := service.CatalogTorrentDomain().CheckCategoryExists(ctx, categoryId)
	if err != nil {
		return err
	}
	return nil
}

func (s *sCatalogTorrentUsecase) parseAndModifyTorrent(ctx context.Context, in catalogin.TorrentUploadInp) (*metainfo.MetaInfo, *metainfo.Info, error) {
	file, err := in.File.Open()
	if err != nil {
		return nil, nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.open_upload_failed"))
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.read_upload_failed"))
	}

	mi, err := metainfo.Load(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.invalid_format"))
	}

	info, err := mi.UnmarshalInfo()
	if err != nil {
		return nil, nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.unmarshal_info_failed"))
	}

	privateFlag := true
	info.Private = &privateFlag

	sourceFlag := s.getCatalogConfigCache(ctx, consts.SiteConfigCatalogTorrentSource).String()
	if sourceFlag == "" {
		sourceFlag, _ = consts.SiteConfigDefaults[consts.SiteConfigCatalogTorrentSource].(string)
	}
	if sourceFlag == "" {
		sourceFlag = "NextPT"
	}
	info.Source = sourceFlag

	mi.Announce = ""
	mi.AnnounceList = nil

	infoBytes, err := bencode.Marshal(info)
	if err != nil {
		return nil, nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.reencode_info_failed"))
	}
	mi.InfoBytes = infoBytes

	return mi, &info, nil
}

func (s *sCatalogTorrentUsecase) calculateAndCheckInfoHash(ctx context.Context, mi *metainfo.MetaInfo) ([]byte, string, error) {
	newInfoHash := mi.HashInfoBytes()
	infoHashBytes := newInfoHash.Bytes()
	infoHashHex := newInfoHash.HexString()

	err := service.CatalogTorrentDomain().CheckInfoHashExists(ctx, infoHashBytes)
	if err != nil {
		return nil, "", err
	}
	return infoHashBytes, infoHashHex, nil
}

func (s *sCatalogTorrentUsecase) extractTorrentMetadata(info *metainfo.Info) (uint64, uint) {
	var totalSize uint64
	var fileCount uint = 0

	if len(info.Files) > 0 {
		for _, f := range info.Files {
			totalSize += uint64(f.Length)
			fileCount++
		}
	} else {
		totalSize = uint64(info.Length)
		fileCount = 1
	}
	return totalSize, fileCount
}

func (s *sCatalogTorrentUsecase) saveTorrentToDB(ctx context.Context, actor *model.Actor, in catalogin.TorrentUploadInp, infoHashBytes []byte, releaseData *uploadReleaseData, totalSize uint64, fileCount uint, info *metainfo.Info) (uint64, error) {
	if releaseData == nil {
		releaseData = &uploadReleaseData{Name: s.resolveUploadName(strings.TrimSpace(in.Name), info.Name)}
	}

	torrentInsert := &entity.CatalogTorrent{
		InfoHash:    infoHashBytes,
		Name:        releaseData.Name,
		SubTitle:    in.SubTitle,
		CategoryId:  in.CategoryId,
		Description: in.Description,
		FileName:    in.File.Filename,
		Size:        totalSize,
		FileCount:   fileCount,
		OwnerId:     actor.Id,
		Anonymous:   in.Anonymous,
		Visible:     true,
	}
	if releaseData != nil && len(releaseData.Fields) > 0 {
		torrentInsert.ReleaseFields = gjson.New(releaseData.Fields)
	}

	var filesToInsert []entity.CatalogTorrentFile
	if len(info.Files) > 0 {
		for _, f := range info.Files {
			filesToInsert = append(filesToInsert, entity.CatalogTorrentFile{
				FilePath: filepath.Join(f.Path...),
				Size:     uint64(f.Length),
			})
		}
	} else {
		filesToInsert = append(filesToInsert, entity.CatalogTorrentFile{
			FilePath: info.Name,
			Size:     uint64(info.Length),
		})
	}

	return service.CatalogTorrentDomain().SaveTorrent(ctx, torrentInsert, filesToInsert)
}

// Reward 赞赏种子
func (s *sCatalogTorrentUsecase) Reward(ctx context.Context, actor *model.Actor, in catalogin.TorrentRewardInp) (*catalogout.TorrentRewardOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	userId := actor.Id

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 获取种子信息并校验可见性
		torrent, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
		if err != nil {
			return err
		}

		if torrent.OwnerId == userId {
			return gerror.New(gi18n.T(ctx, "catalog.reward.self_denied"))
		}

		// 调用 Economy Bonus 域的 TransferBonus 方法
		remarkFrom := fmt.Sprintf("Reward torrent #%d", in.Id)
		remarkTo := fmt.Sprintf("Received reward for torrent #%d", in.Id)
		err = service.EconomyBonusUsecase().TransferBonus(ctx, userId, torrent.OwnerId, in.Amount, consts.EconomyBonusTargetTypeTorrent, in.Id, remarkFrom, remarkTo)
		if err != nil {
			return err
		}

		// 写入赞赏记录
		err = service.CatalogTorrentDomain().InsertTorrentReward(ctx, &entity.CatalogTorrentReward{TorrentId: in.Id, UserId: userId, Amount: in.Amount})
		if err != nil {
			return gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.write_reward_failed"))
		}

		// 更新种子赞赏统计
		err = service.CatalogTorrentDomain().IncrementTorrentRewardStats(ctx, in.Id, in.Amount)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &catalogout.TorrentRewardOut{}, nil
}

// RewardList 获取种子赞赏列表
func (s *sCatalogTorrentUsecase) RewardList(ctx context.Context, actor *model.Actor, in catalogin.TorrentRewardListInp) (*catalogout.TorrentRewardListOut, error) {
	_, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return nil, err
	}

	summaries, total, err := service.CatalogTorrentDomain().QueryTorrentRewards(ctx, in.Id, in.Page, in.Size)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.torrent.get_reward_list_failed"))
	}

	userIds := make([]uint64, 0, len(summaries))
	for _, e := range summaries {
		userIds = append(userIds, e.UserId)
	}
	usernameMap := s.loadUsernameMap(ctx, userIds)

	var list []catalogout.TorrentRewardItem
	for _, e := range summaries {
		list = append(list, catalogout.TorrentRewardItem{
			UserId:       e.UserId,
			Username:     usernameMap[e.UserId],
			Amount:       e.Amount,
			RewardCount:  e.RewardCount,
			LastRewardAt: s.formatTime(e.LastRewardAt),
		})
	}

	return &catalogout.TorrentRewardListOut{
		List:  list,
		Total: total,
	}, nil
}

func (s *sCatalogTorrentUsecase) Bookmark(ctx context.Context, actor *model.Actor, in catalogin.TorrentBookmarkInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	userId := actor.Id
	_, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return err
	}
	return service.CatalogTorrentDomain().Bookmark(ctx, userId, in.Id)
}

func (s *sCatalogTorrentUsecase) Unbookmark(ctx context.Context, actor *model.Actor, in catalogin.TorrentUnbookmarkInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	userId := actor.Id
	return service.CatalogTorrentDomain().Unbookmark(ctx, userId, in.Id)
}

func (s *sCatalogTorrentUsecase) ListBookmarkedTorrents(ctx context.Context, actor *model.Actor, in catalogin.TorrentBookmarkListInp) (*catalogout.TorrentBookmarkListOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}

	torrents, total, err := service.CatalogTorrentDomain().QueryBookmarkedTorrents(ctx, actor, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	list := s.formatTorrentListItems(ctx, actor, torrents)

	return &catalogout.TorrentBookmarkListOut{
		List:  list,
		Total: total,
	}, nil
}

func (s *sCatalogTorrentUsecase) formatTorrentListItems(ctx context.Context, actor *model.Actor, entities []entity.CatalogTorrent) []catalogout.TorrentListItem {
	var ownerIds []uint64
	for _, e := range entities {
		if !s.shouldHideTorrentOwner(actor, e) && e.OwnerId > 0 {
			ownerIds = append(ownerIds, e.OwnerId)
		}
	}

	userMap := make(map[uint64]string)
	if len(ownerIds) > 0 {
		users, _ := service.IamUserDomain().GetUsersByIds(ctx, ownerIds)
		for _, u := range users {
			userMap[u.Id] = u.Username
		}
	}

	var list []catalogout.TorrentListItem
	for _, e := range entities {
		ownerId := e.OwnerId
		ownerName := userMap[e.OwnerId]
		if s.shouldHideTorrentOwner(actor, e) {
			ownerId = 0
			ownerName = ""
		}
		list = append(list, catalogout.TorrentListItem{
			Id:         e.Id,
			Name:       e.Name,
			SubTitle:   e.SubTitle,
			CategoryId: e.CategoryId,
			Size:       e.Size,
			FileCount:  e.FileCount,
			SpState:    e.SpState,
			SpExpireAt: s.formatTime(e.SpExpireAt),
			IsFeatured: e.IsFeatured,
			IsPinned:   e.IsPinned,
			Seeders:    e.Seeders,
			Leechers:   e.Leechers,
			Snatched:   e.TimesCompleted,
			LikeCount:  e.LikeCount,
			OwnerId:    ownerId,
			OwnerName:  ownerName,
			Anonymous:  e.Anonymous,
			CreatedAt:  e.CreatedAt.String(),
		})
	}
	return list
}

func (s *sCatalogTorrentUsecase) formatTime(value *gtime.Time) string {
	if value == nil {
		return ""
	}
	return value.String()
}

func (s *sCatalogTorrentUsecase) shouldHideTorrentOwner(actor *model.Actor, torrent entity.CatalogTorrent) bool {
	if !torrent.Anonymous {
		return false
	}
	return actor == nil || (!actor.IsStaff && actor.Id != torrent.OwnerId)
}

func (s *sCatalogTorrentUsecase) loadUsernameMap(ctx context.Context, userIds []uint64) map[uint64]string {
	userMap := make(map[uint64]string)
	if len(userIds) == 0 {
		return userMap
	}

	uniqueIds := make([]uint64, 0, len(userIds))
	seen := make(map[uint64]struct{}, len(userIds))
	for _, id := range userIds {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIds = append(uniqueIds, id)
	}
	if len(uniqueIds) == 0 {
		return userMap
	}

	users, err := service.IamUserDomain().GetUsersByIds(ctx, uniqueIds)
	if err != nil {
		return userMap
	}
	for _, user := range users {
		userMap[user.Id] = user.Username
	}
	return userMap
}

func (s *sCatalogTorrentUsecase) GetTorrent(ctx context.Context, actor *model.Actor, in catalogin.TorrentGetInp) (*catalogout.TorrentDetailOut, error) {
	torrent, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return nil, err
	}

	isBookmarked := false
	isLiked := false
	if actor != nil && actor.Id > 0 {
		isBookmarked, _ = service.CatalogTorrentDomain().CheckTorrentBookmarked(ctx, in.Id, actor.Id)
		isLiked, _ = service.CatalogTorrentDomain().CheckTorrentLiked(ctx, in.Id, actor.Id)
	}

	listItems := s.formatTorrentListItems(ctx, actor, []entity.CatalogTorrent{*torrent})
	var item catalogout.TorrentListItem
	if len(listItems) > 0 {
		item = listItems[0]
	}

	var releaseFields map[string]any
	if torrent.ReleaseFields != nil {
		_ = torrent.ReleaseFields.Scan(&releaseFields)
	}

	return &catalogout.TorrentDetailOut{
		TorrentListItem: item,
		Description:     torrent.Description,
		ReleaseFields:   releaseFields,
		IsBookmarked:    isBookmarked,
		IsLiked:         isLiked,
	}, nil
}

func (s *sCatalogTorrentUsecase) ToggleLike(ctx context.Context, actor *model.Actor, in catalogin.TorrentToggleLikeInp) (*catalogout.TorrentToggleLikeOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}
	_, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return nil, err
	}

	isLiked, err := service.CatalogTorrentDomain().ToggleLike(ctx, in.Id, actor.Id)
	if err != nil {
		return nil, err
	}
	return &catalogout.TorrentToggleLikeOut{IsLiked: isLiked}, nil
}

func (s *sCatalogTorrentUsecase) ListLikes(ctx context.Context, actor *model.Actor, in catalogin.TorrentLikeListInp) (*catalogout.TorrentLikeListOut, error) {
	_, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return nil, err
	}

	likes, total, err := service.CatalogTorrentDomain().QueryTorrentLikes(ctx, in.Id, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	var userIds []uint64
	for _, e := range likes {
		userIds = append(userIds, e.UserId)
	}

	userMap := make(map[uint64]string)
	if len(userIds) > 0 {
		users, _ := service.IamUserDomain().GetUsersByIds(ctx, userIds)
		for _, u := range users {
			userMap[u.Id] = u.Username
		}
	}

	var list []catalogout.TorrentLikeItem
	for _, e := range likes {
		list = append(list, catalogout.TorrentLikeItem{
			Id:        e.Id,
			UserId:    e.UserId,
			Username:  userMap[e.UserId],
			CreatedAt: e.CreatedAt.String(),
		})
	}

	return &catalogout.TorrentLikeListOut{
		List:  list,
		Total: total,
	}, nil
}

func (s *sCatalogTorrentUsecase) Update(ctx context.Context, actor *model.Actor, in catalogin.TorrentUpdateInp) (*catalogout.TorrentUpdateOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}

	torrent, err := service.CatalogTorrentDomain().GetTorrentById(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if torrent == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.torrent.not_found"))
	}

	if !s.canEditTorrent(actor, torrent) {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.forbidden"))
	}

	categoryId := torrent.CategoryId
	if in.CategoryId > 0 {
		categoryId = in.CategoryId
	}
	category, err := service.CatalogCategoryDomain().GetCategoryById(ctx, categoryId)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.category.invalid"))
	}

	releaseData, err := s.prepareUpdateReleaseData(ctx, category, torrent, in)
	if err != nil {
		return nil, err
	}

	name := releaseData.Name
	subTitle := strings.TrimSpace(in.SubTitle)
	description := strings.TrimSpace(in.Description)
	updateData := model.CatalogTorrentUpdate{
		Name:        name,
		SubTitle:    subTitle,
		CategoryId:  categoryId,
		Description: description,
		Anonymous:   in.Anonymous,
	}
	if releaseData != nil && len(releaseData.Fields) > 0 {
		updateData.ReleaseFields = gjson.New(releaseData.Fields)
	}

	err = service.CatalogTorrentDomain().UpdateTorrent(ctx, in.Id, updateData)
	if err != nil {
		return nil, err
	}

	return &catalogout.TorrentUpdateOut{Success: true}, nil
}

func (s *sCatalogTorrentUsecase) canEditTorrent(actor *model.Actor, torrent *entity.CatalogTorrent) bool {
	return actor != nil && torrent != nil && (actor.IsStaff || torrent.OwnerId == actor.Id)
}

func (s *sCatalogTorrentUsecase) ListFiles(ctx context.Context, actor *model.Actor, in catalogin.TorrentFileListInp) (*catalogout.TorrentFileListOut, error) {
	_, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return nil, err
	}

	files, err := service.CatalogTorrentDomain().GetTorrentFiles(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	var list []catalogout.TorrentFileItem
	for _, f := range files {
		list = append(list, catalogout.TorrentFileItem{
			Path: f.FilePath,
			Size: f.Size,
		})
	}

	return &catalogout.TorrentFileListOut{List: list}, nil
}

func (s *sCatalogTorrentUsecase) ListPeers(ctx context.Context, actor *model.Actor, in catalogin.TorrentPeerListInp) (*catalogout.TorrentPeerListOut, error) {
	_, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return nil, err
	}

	peers, err := service.TrackerPeerDomain().GetActivePeers(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	var userIds []uint64
	for _, p := range peers {
		userIds = append(userIds, p.UserId)
	}

	userMap := make(map[uint64]string)
	if len(userIds) > 0 {
		users, _ := service.IamUserDomain().GetUsersByIds(ctx, userIds)
		for _, u := range users {
			userMap[u.Id] = u.Username
		}
	}

	var list []catalogout.TorrentPeerItem
	for _, p := range peers {
		startedAt := ""
		if p.StartedAt != nil {
			startedAt = p.StartedAt.String()
		}
		list = append(list, catalogout.TorrentPeerItem{
			UserId:     p.UserId,
			Username:   userMap[p.UserId],
			IsSeeder:   p.IsSeeder,
			Uploaded:   p.Uploaded,
			Downloaded: p.Downloaded,
			StartedAt:  startedAt,
		})
	}

	return &catalogout.TorrentPeerListOut{List: list}, nil
}

func (s *sCatalogTorrentUsecase) Report(ctx context.Context, actor *model.Actor, in catalogin.TorrentReportInp) (*catalogout.TorrentReportOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}

	_, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return nil, err
	}

	err = service.ModReportUsecase().Create(ctx, actor, modin.CreateReportInp{
		TargetType: "torrent",
		TargetId:   in.Id,
		Reason:     in.Reason,
	})
	if err != nil {
		return nil, err
	}

	return &catalogout.TorrentReportOut{Success: true}, nil
}

func (s *sCatalogTorrentUsecase) getCatalogConfigCache(ctx context.Context, key string) *gvar.Var {
	cacheKey := service.SysCache().KeySiteConfigFullPath(ctx, key)
	val, err := gcache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		return service.SiteConfigDomain().GetByPath(ctx, key).Val(), nil
	}, 5*time.Minute)
	if err != nil || val.IsNil() {
		return service.SiteConfigDomain().GetByPath(ctx, key)
	}
	return gvar.New(val.Val())
}
