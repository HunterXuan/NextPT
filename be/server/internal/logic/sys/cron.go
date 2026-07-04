package sys

import (
	"context"
	"os"
	"time"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
)

type sSysCron struct{}

func init() {
	cron := NewCron()
	service.RegisterSysCron(cron)
}

func NewCron() *sSysCron {
	return &sSysCron{}
}

// Start 启动所有系统定时任务，在全局 Init 中调用
func (s *sSysCron) Start(ctx context.Context) {
	gcron.AddSingleton(ctx, "0 */5 * * * *", s.runWrapper("tracker_ghost_peer_cleanup", 180, func(ctx context.Context) error {
		return service.TrackerSyncUsecase().CleanupGhostPeers(ctx)
	}), "tracker_ghost_peer_cleanup")

	gcron.AddSingleton(ctx, "0 * * * * *", s.runWrapper("tracker_data_alignment", 60, func(ctx context.Context) error {
		return service.TrackerSyncUsecase().SyncTorrentData(ctx)
	}), "tracker_data_alignment")

	gcron.AddSingleton(ctx, "30 * * * * *", s.runWrapper("mod_user_expiry_cleanup", 300, func(ctx context.Context) error {
		rows, err := service.ModUserUsecase().CleanupExpired(ctx)
		if err != nil {
			return err
		}
		if rows > 0 {
			glog.Infof(ctx, "[Cron] Cleanup expired user mods completed. Restored %d mod records.", rows)
		}
		return nil
	}), "mod_user_expiry_cleanup")

	gcron.AddSingleton(ctx, "45 * * * * *", s.runWrapper("iam_invite_expiry_cleanup", 300, func(ctx context.Context) error {
		rows, err := service.IamInviteUsecase().CleanupExpired(ctx)
		if err != nil {
			return err
		}
		if rows > 0 {
			glog.Infof(ctx, "[Cron] Cleanup expired invites completed. Expired %d invite records.", rows)
		}
		return nil
	}), "iam_invite_expiry_cleanup")

	gcron.AddSingleton(ctx, "@hourly", s.runWrapper("tracker_bonus_points", 1800, func(ctx context.Context) error {
		return service.EconomyBonusUsecase().DistributeBonusPoints(ctx)
	}), "tracker_bonus_points")

	// 每天清理历史定时任务日志，保留最近 7 天
	gcron.AddSingleton(ctx, "@daily", s.runWrapper("sys_cron_log_cleanup", 3600, func(ctx context.Context) error {
		return s.cleanupCronLogs(ctx)
	}), "sys_cron_log_cleanup")

	// 每天清理 Tracker 幂等表，保留最近 2 天即可 (防重放通常只需较短时间)
	gcron.AddSingleton(ctx, "@daily", s.runWrapper("tracker_idempotency_cleanup", 3600, func(ctx context.Context) error {
		return service.TrackerSyncUsecase().CleanupIdempotency(ctx)
	}), "tracker_idempotency_cleanup")
}

// cleanupCronLogs 清理历史的定时任务日志，保留最近 7 天
func (s *sSysCron) cleanupCronLogs(ctx context.Context) error {
	retentionDays := 7
	beforeTime := gtime.Now().Add(-time.Duration(retentionDays*24) * time.Hour)
	result, err := dao.SysCronLog.Ctx(ctx).Where("created_at < ?", beforeTime).Delete()
	if err != nil {
		glog.Error(ctx, "[Cron] Failed to cleanup cron logs:", err)
		return err
	}
	rows, _ := result.RowsAffected()
	if rows > 0 {
		glog.Infof(ctx, "[Cron] Cleanup cron logs completed. Removed %d expired log entries before %s.", rows, beforeTime)
	}
	return nil
}

func (s *sSysCron) tryAcquireCronLock(ctx context.Context, lockName string, expireSeconds int) bool {
	lockKey := service.SysCache().KeySysCronLock(ctx, lockName)
	v, err := g.Redis().Do(ctx, "SET", lockKey, gtime.Now().Unix(), "NX", "EX", expireSeconds)
	if err != nil {
		glog.Error(ctx, "[Cron] Failed to acquire distributed lock for "+lockName+":", err)
		return false
	}
	if v.String() != "OK" {
		return false
	}
	return true
}

func (s *sSysCron) runWrapper(jobName string, lockExpire int, jobFunc func(ctx context.Context) error) func(ctx context.Context) {
	return func(ctx context.Context) {
		// 分布式锁，抢不到直接返回
		if !s.tryAcquireCronLock(ctx, jobName, lockExpire) {
			return
		}

		nodeIp, _ := os.Hostname()
		if nodeIp == "" {
			nodeIp = gctx.CtxId(ctx) // 降级使用 trace id
		}
		startTime := gtime.Now()

		// 写入一条进行中的日志
		logId, err := dao.SysCronLog.Ctx(ctx).Data(entity.SysCronLog{
			JobName:   jobName,
			NodeIp:    nodeIp,
			Status:    consts.SysCronStatusRunning,
			CreatedAt: startTime,
			UpdatedAt: startTime,
		}).InsertAndGetId()

		if err != nil {
			glog.Error(ctx, "[Cron] Failed to insert cron log for "+jobName+":", err)
		}

		// 捕获 panic 并执行业务
		var jobErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					jobErr = gerror.Newf("panic: %v", r)
				}
			}()
			jobErr = jobFunc(ctx)
		}()

		endTime := gtime.Now()
		durationMs := int(endTime.TimestampMilli() - startTime.TimestampMilli())

		status := consts.SysCronStatusSuccess
		errorMsg := ""
		if jobErr != nil {
			status = consts.SysCronStatusFailed
			errorMsg = jobErr.Error()
		}

		// 更新日志状态
		if logId > 0 {
			dbCtx := gctx.NeverDone(ctx)
			_, err = dao.SysCronLog.Ctx(dbCtx).Data(g.Map{
				"status":        status,
				"duration_ms":   durationMs,
				"error_message": errorMsg,
				"updated_at":    endTime,
			}).Where("id", logId).Update()
			if err != nil {
				glog.Error(ctx, "[Cron] Failed to update cron log for "+jobName+":", err)
			}
		}
	}
}

func (s *sSysCron) AdminListCronLogs(ctx context.Context, jobName string, page, size int) ([]*entity.SysCronLog, int, error) {
	m := dao.SysCronLog.Ctx(ctx).Where("job_name", jobName)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.SysCronLog
	err = m.Page(page, size).Order("id desc").Scan(&list)
	return list, total, err
}
