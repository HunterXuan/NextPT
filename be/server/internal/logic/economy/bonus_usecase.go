package economy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/economyin"
	"server/internal/model/out/economyout"
	"server/internal/service"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/glog"
)

type sEconomyBonusUsecase struct{}

func init() {
	service.RegisterEconomyBonusUsecase(NewEconomyBonusUsecase())
}

func NewEconomyBonusUsecase() *sEconomyBonusUsecase {
	return &sEconomyBonusUsecase{}
}

func (s *sEconomyBonusUsecase) ListMyBonusLogs(ctx context.Context, actor *model.Actor, in economyin.BonusLogsInp) (*economyout.BonusLogsOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	page := in.Page
	size := in.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	logs, total, err := service.EconomyBonusDomain().QueryBonusLogs(ctx, actor.Id, in.Action, page, size)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "economy.bonus.get_log_list_failed"))
	}

	var list []economyout.BonusLogItem
	for _, l := range logs {
		list = append(list, economyout.BonusLogItem{
			Id:           l.Id,
			UserId:       l.UserId,
			Amount:       l.Amount,
			BalanceAfter: l.BalanceAfter,
			Action:       l.Action,
			TargetType:   l.TargetType,
			TargetId:     l.TargetId,
			Remark:       l.Remark,
			CreatedAt:    l.CreatedAt,
		})
	}

	return &economyout.BonusLogsOut{
		Page:  page,
		Size:  size,
		Total: total,
		List:  list,
	}, nil
}

func (s *sEconomyBonusUsecase) TransferBonus(ctx context.Context, fromUserId, toUserId uint64, amount float64, targetType string, targetId uint64, remarkFrom, remarkTo string) error {
	if amount <= 0 {
		return gerror.New(gi18n.T(ctx, "economy.bonus.invalid_amount"))
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// Debit sender
		err := service.EconomyBonusDomain().DebitBonusIfEnough(ctx, fromUserId, amount)
		if err != nil {
			return err
		}

		// Credit receiver
		err = service.EconomyBonusDomain().CreditBonus(ctx, toUserId, amount)
		if err != nil {
			return err
		}

		// Read balances
		newSenderBonus, err := service.EconomyBonusDomain().GetUserBonus(ctx, fromUserId)
		if err != nil {
			return err
		}
		newReceiverBonus, err := service.EconomyBonusDomain().GetUserBonus(ctx, toUserId)
		if err != nil {
			return err
		}

		// Write logs
		logs := []entity.EconomyBonusLog{
			{
				UserId:       fromUserId,
				Amount:       -amount,
				BalanceAfter: newSenderBonus,
				Action:       consts.EconomyBonusActionTransferSent,
				TargetType:   targetType,
				TargetId:     targetId,
				Remark:       remarkFrom,
			},
			{
				UserId:       toUserId,
				Amount:       amount,
				BalanceAfter: newReceiverBonus,
				Action:       consts.EconomyBonusActionTransferReceived,
				TargetType:   targetType,
				TargetId:     targetId,
				Remark:       remarkTo,
			},
		}
		return service.EconomyBonusDomain().InsertBonusLogs(ctx, logs)
	})
}

func (s *sEconomyBonusUsecase) AddBonus(ctx context.Context, userId uint64, amount float64, action string, targetType string, targetId uint64, remark string, period string) error {
	if amount == 0 {
		return nil
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if amount > 0 {
			err := service.EconomyBonusDomain().CreditBonus(ctx, userId, amount)
			if err != nil {
				return err
			}
		} else {
			err := service.EconomyBonusDomain().DebitBonusIfEnough(ctx, userId, -amount)
			if err != nil {
				return err
			}
		}

		newBonus, err := service.EconomyBonusDomain().GetUserBonus(ctx, userId)
		if err != nil {
			return err
		}

		log := entity.EconomyBonusLog{
			UserId:       userId,
			Amount:       amount,
			BalanceAfter: newBonus,
			Action:       action,
			TargetType:   targetType,
			TargetId:     targetId,
			Remark:       remark,
			Period:       period,
		}
		return service.EconomyBonusDomain().InsertBonusLog(ctx, log)
	})
}

func (s *sEconomyBonusUsecase) GetMyHourlyBonus(ctx context.Context, actor *model.Actor, in economyin.HourlyBonusInp) (*economyout.HourlyBonusOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	hourly, err := s.CalculateHourlyBonus(ctx, actor.Id)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "economy.bonus.calc_hourly_failed"))
	}
	return &economyout.HourlyBonusOut{HourlyBonus: hourly}, nil
}

// DistributeBonusPoints 魔力值自动发放 (NexusPHP Formula)
func (s *sEconomyBonusUsecase) DistributeBonusPoints(ctx context.Context) error {
	// 1. 获取所有当前有做种记录的用户 ID
	activeUserIds, err := service.TrackerPeerDomain().GetSeedingUsers(ctx)
	if err != nil {
		glog.Error(ctx, "[Cron] Bonus calculation fetch seeding users failed:", err)
		return err
	}
	if len(activeUserIds) == 0 {
		return nil
	}

	// 从数据库动态获取公式参数，支持管理员随时调整，并使用默认值兜底
	T0 := s.getEconomyConfigCache(ctx, consts.SiteConfigTrackerBonusT0).Float64()
	N0 := s.getEconomyConfigCache(ctx, consts.SiteConfigTrackerBonusN0).Float64()
	B0 := s.getEconomyConfigCache(ctx, consts.SiteConfigTrackerBonusB0).Float64()
	L := s.getEconomyConfigCache(ctx, consts.SiteConfigTrackerBonusL).Float64()
	basePoints := s.getEconomyConfigCache(ctx, consts.SiteConfigTrackerBonusBase).Float64()

	now := time.Now()
	period := fmt.Sprintf("hourly:%s", now.Format("2006010215"))
	totalRewardedUsers := 0

	// 2. 按用户分批处理，避免内存爆炸
	chunkSize := 500
	for i := 0; i < len(activeUserIds); i += chunkSize {
		end := i + chunkSize
		if end > len(activeUserIds) {
			end = len(activeUserIds)
		}
		chunk := activeUserIds[i:end]
		var seeders []entity.TrackerPeer

		for _, uId := range chunk {
			peers, err := service.TrackerPeerDomain().GetUserSeedingPeers(ctx, uId)
			if err == nil && len(peers) > 0 {
				seeders = append(seeders, peers...)
			}
		}

		// 收集这批记录涉及到的唯一种子 ID
		torrentIdsSet := make(map[uint64]bool)
		for _, p := range seeders {
			torrentIdsSet[p.TorrentId] = true
		}
		var torrentIds []any
		for id := range torrentIdsSet {
			torrentIds = append(torrentIds, id)
		}

		var torrentIdsUint64 []uint64
		for _, id := range torrentIds {
			torrentIdsUint64 = append(torrentIdsUint64, id.(uint64))
		}
		torrents, _ := service.CatalogTorrentDomain().GetTorrentsByIds(ctx, torrentIdsUint64)
		torrentMap := s.buildBonusTorrentSnapshots(torrents)

		config := economyin.BonusFormulaConfig{T0: T0, N0: N0, B0: B0, L: L, BasePoints: basePoints}
		userPeersMap := make(map[uint64][]economyin.BonusPeerSnapshot)
		for _, p := range seeders {
			userPeersMap[p.UserId] = append(userPeersMap[p.UserId], economyin.BonusPeerSnapshot{
				UserId:    p.UserId,
				TorrentId: p.TorrentId,
			})
		}

		// 执行单批次的单条无事务积分更新 -> 改为有事务的积分更新
		for userId, userPeers := range userPeersMap {
			totalBonus := service.EconomyBonusDomain().CalculateBonusForPeers(userPeers, torrentMap, config, now)

			err := s.AddBonus(ctx, userId, totalBonus, consts.EconomyBonusActionSeedBonus, "", 0, "System hourly seeding bonus distribution", period)
			if err != nil {
				if strings.Contains(err.Error(), "1062") {
					glog.Warningf(ctx, "[Cron] Bonus already distributed for user %d at period %s, skipping", userId, period)
					continue
				}
				glog.Error(ctx, "[Cron] Bonus calculation failed to update user credit:", err)
				return err
			}
			totalRewardedUsers++
		}
	}

	glog.Infof(ctx, "[Cron] Bonus calculation completed. Rewarded %d users in total.", totalRewardedUsers)
	return nil
}

// CalculateHourlyBonus 计算指定用户当前每小时可获得魔力值
func (s *sEconomyBonusUsecase) CalculateHourlyBonus(ctx context.Context, userId uint64) (float64, error) {
	// 获取做种记录
	peers, err := service.TrackerPeerDomain().GetUserSeedingPeers(ctx, userId)
	if err != nil {
		return 0, err
	}
	if len(peers) == 0 {
		return 0, nil
	}

	// 提取 torrent_id
	var torrentIdsUint64 []uint64
	for _, p := range peers {
		torrentIdsUint64 = append(torrentIdsUint64, p.TorrentId)
	}

	// 批量获取种子详情
	torrents, err := service.CatalogTorrentDomain().GetTorrentsByIds(ctx, torrentIdsUint64)
	if err != nil {
		return 0, err
	}

	torrentMap := s.buildBonusTorrentSnapshots(torrents)
	peerSnapshots := s.buildBonusPeerSnapshots(peers)

	config := economyin.BonusFormulaConfig{
		T0:         s.getEconomyConfigCache(ctx, consts.SiteConfigTrackerBonusT0).Float64(),
		N0:         s.getEconomyConfigCache(ctx, consts.SiteConfigTrackerBonusN0).Float64(),
		B0:         s.getEconomyConfigCache(ctx, consts.SiteConfigTrackerBonusB0).Float64(),
		L:          s.getEconomyConfigCache(ctx, consts.SiteConfigTrackerBonusL).Float64(),
		BasePoints: s.getEconomyConfigCache(ctx, consts.SiteConfigTrackerBonusBase).Float64(),
	}

	totalBonus := service.EconomyBonusDomain().CalculateBonusForPeers(peerSnapshots, torrentMap, config, time.Now())
	return totalBonus, nil
}

func (s *sEconomyBonusUsecase) buildBonusPeerSnapshots(peers []entity.TrackerPeer) []economyin.BonusPeerSnapshot {
	snapshots := make([]economyin.BonusPeerSnapshot, 0, len(peers))
	for _, p := range peers {
		snapshots = append(snapshots, economyin.BonusPeerSnapshot{
			UserId:    p.UserId,
			TorrentId: p.TorrentId,
		})
	}
	return snapshots
}

func (s *sEconomyBonusUsecase) buildBonusTorrentSnapshots(torrents []*entity.CatalogTorrent) map[uint64]economyin.BonusTorrentSnapshot {
	torrentMap := make(map[uint64]economyin.BonusTorrentSnapshot, len(torrents))
	for _, t := range torrents {
		torrentMap[t.Id] = economyin.BonusTorrentSnapshot{
			Id:        t.Id,
			Size:      t.Size,
			Seeders:   t.Seeders,
			CreatedAt: t.CreatedAt,
		}
	}
	return torrentMap
}

func (s *sEconomyBonusUsecase) getEconomyConfigCache(ctx context.Context, key string) *gvar.Var {
	cacheKey := service.SysCache().KeySiteConfigFullPath(ctx, key)
	val, err := gcache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		return service.SiteConfigDomain().GetByPath(ctx, key).Val(), nil
	}, 10*time.Minute)
	if err != nil || val.IsNil() {
		return service.SiteConfigDomain().GetByPath(ctx, key)
	}
	return gvar.New(val.Val())
}
