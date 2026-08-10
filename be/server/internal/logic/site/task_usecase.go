package site

import (
	"context"
	"time"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/sitein"
	"server/internal/model/out/siteout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	siteTaskConfigCacheDuration = 10 * time.Minute
	siteTaskGiB                 = uint64(1024 * 1024 * 1024)
)

type sSiteTaskUsecase struct{}

func init() {
	service.RegisterSiteTaskUsecase(NewSiteTaskUsecase())
}

func NewSiteTaskUsecase() *sSiteTaskUsecase {
	return &sSiteTaskUsecase{}
}

func (s *sSiteTaskUsecase) List(ctx context.Context, actor *model.Actor) (*siteout.TaskListOut, error) {
	if actor == nil {
		return nil, gerror.New("unauthorized")
	}
	tasks, err := s.loadTasksCache(ctx)
	if err != nil {
		return nil, err
	}
	instances, err := service.SiteTaskDomain().ListUserTasks(ctx, actor.Id)
	if err != nil {
		return nil, err
	}
	instanceByCycle := make(map[string]entity.SiteUserTask, len(instances))
	for _, instance := range instances {
		key := instance.TaskKey + ":" + instance.CycleKey
		if _, exists := instanceByCycle[key]; !exists {
			instanceByCycle[key] = instance
		}
	}
	items := make([]*siteout.TaskItem, 0, len(tasks))
	for _, task := range tasks {
		item := &siteout.TaskItem{
			Key:             task.Key,
			Enabled:         task.Enabled,
			Cycle:           task.Cycle,
			NameI18N:        task.NameI18N,
			DescriptionI18N: task.DescriptionI18N,
			Rule:            task.Rule,
			Rewards:         task.Rewards,
		}
		cycleKey, _, _ := task.CyclePeriod(time.Now())
		if instance, exists := instanceByCycle[task.Key+":"+cycleKey]; exists {
			userTask := s.userTaskItem(instance)
			item.UserTask = &userTask
		}
		items = append(items, item)
	}
	return &siteout.TaskListOut{List: items}, nil
}

func (s *sSiteTaskUsecase) Claim(ctx context.Context, actor *model.Actor, in sitein.TaskClaimInp) (*siteout.TaskClaimOut, error) {
	if actor == nil {
		return nil, gerror.New("unauthorized")
	}
	tasks, err := s.loadTasksCache(ctx)
	if err != nil {
		return nil, err
	}
	task := tasks.Find(in.Key)
	if task == nil || !task.Enabled {
		return nil, gerror.New("task is unavailable")
	}
	if task.Rule.Type == consts.SiteTaskRuleTypeRoleLevelReached && actor.IsStaff {
		return nil, gerror.New("staff users cannot claim role tasks")
	}

	now := time.Now()
	cycleKey, cycleStartedAt, cycleEndedAt := task.CyclePeriod(now)
	existing, err := service.SiteTaskDomain().ListUserTasks(ctx, actor.Id)
	if err != nil {
		return nil, err
	}
	for _, item := range existing {
		if item.TaskKey == task.Key && item.CycleKey == cycleKey {
			return nil, gerror.New("task has already been claimed for this cycle")
		}
	}

	claimedAt := gtime.NewFromTime(now)
	instance := entity.SiteUserTask{
		UserId:         actor.Id,
		TaskKey:        task.Key,
		CycleKey:       cycleKey,
		CycleStartedAt: cycleStartedAt,
		CycleEndedAt:   cycleEndedAt,
		Status:         consts.SiteTaskStatusActive,
		Target:         s.taskTarget(*task),
		TaskSnapshot:   gjson.New(task),
		ClaimedAt:      claimedAt,
		CreatedAt:      claimedAt,
		UpdatedAt:      claimedAt,
	}
	id, err := service.SiteTaskDomain().CreateUserTask(ctx, instance)
	if err != nil {
		return nil, err
	}
	instance.Id = id
	item := s.userTaskItem(instance)
	return &siteout.TaskClaimOut{UserTask: item}, nil
}

func (s *sSiteTaskUsecase) ClaimReward(ctx context.Context, actor *model.Actor, in sitein.UserTaskRewardClaimInp) error {
	if actor == nil {
		return gerror.New("unauthorized")
	}
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		instance, err := service.SiteTaskDomain().GetUserTask(ctx, in.Id, actor.Id)
		if err != nil {
			return err
		}
		if instance == nil {
			return gerror.New("user task not found")
		}
		if instance.Status != consts.SiteTaskStatusCompleted {
			return gerror.New("task reward is not available")
		}
		var task model.SiteTaskDefinition
		if instance.TaskSnapshot == nil || instance.TaskSnapshot.Scan(&task) != nil {
			return gerror.New("invalid task snapshot")
		}
		rewarded, err := service.SiteTaskDomain().Reward(ctx, instance.Id)
		if err != nil {
			return err
		}
		if !rewarded {
			return gerror.New("task reward is no longer available")
		}
		for _, reward := range task.Rewards {
			if err := s.fulfillReward(ctx, instance.Id, instance.UserId, reward); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	service.IamUserUsecase().InvalidateUserCache(ctx, actor.Id)
	return nil
}

func (s *sSiteTaskUsecase) Settle(ctx context.Context, size int) (int, error) {
	if size <= 0 {
		size = 500
	}
	updated := 0
	now := time.Now()
	var afterId uint64
	for {
		instances, err := service.SiteTaskDomain().ListActiveUserTasks(ctx, afterId, size)
		if err != nil {
			return updated, err
		}
		if len(instances) == 0 {
			return updated, nil
		}
		for _, instance := range instances {
			changed, err := s.settleUserTask(ctx, instance, now)
			if err != nil {
				return updated, err
			}
			if changed {
				updated++
			}
			afterId = instance.Id
		}
		if len(instances) < size {
			return updated, nil
		}
	}
}

func (s *sSiteTaskUsecase) CleanupHistory(ctx context.Context) (int, error) {
	before := gtime.Now().AddDate(0, 0, -consts.SiteTaskHistoryRetentionDays)
	return service.SiteTaskDomain().CleanupHistory(ctx, before)
}

func (s *sSiteTaskUsecase) settleUserTask(ctx context.Context, instance entity.SiteUserTask, now time.Time) (bool, error) {
	var task model.SiteTaskDefinition
	if instance.TaskSnapshot == nil || instance.TaskSnapshot.Scan(&task) != nil {
		return false, gerror.New("invalid task snapshot")
	}
	if !task.SupportsRule() {
		return service.SiteTaskDomain().Expire(ctx, instance.Id)
	}
	progress, err := s.taskProgress(ctx, instance, task, now)
	if err != nil {
		return false, err
	}
	if progress >= instance.Target {
		return service.SiteTaskDomain().Complete(ctx, instance.Id, progress)
	}
	changed := false
	if progress != instance.Progress {
		if err := service.SiteTaskDomain().UpdateProgress(ctx, instance.Id, progress); err != nil {
			return false, err
		}
		changed = true
	}
	if instance.CycleEndedAt != nil && !now.Before(instance.CycleEndedAt.Time) {
		return service.SiteTaskDomain().Expire(ctx, instance.Id)
	}
	return changed, nil
}

func (s *sSiteTaskUsecase) taskProgress(ctx context.Context, instance entity.SiteUserTask, task model.SiteTaskDefinition, now time.Time) (uint64, error) {
	switch task.Rule.Type {
	case consts.SiteTaskRuleTypeRoleLevelReached:
		user, err := service.IamUserDomain().GetUserById(ctx, instance.UserId)
		if err != nil || user == nil {
			return 0, err
		}
		role, err := service.IamRoleDomain().GetRoleById(ctx, user.Role)
		if err != nil || role == nil {
			return 0, err
		}
		if !role.IsStaff && uint64(role.Level) >= instance.Target {
			return instance.Target, nil
		}
		return 0, nil
	case consts.SiteTaskRuleTypeTorrentPublished:
		return s.publishedTorrentCount(ctx, instance, now)
	case consts.SiteTaskRuleTypeSeedDuration, consts.SiteTaskRuleTypeUploaded:
		return s.trackerProgress(ctx, instance, task, now)
	default:
		return 0, gerror.New("unsupported task rule type")
	}
}

func (s *sSiteTaskUsecase) publishedTorrentCount(ctx context.Context, instance entity.SiteUserTask, now time.Time) (uint64, error) {
	columns := dao.CatalogTorrent.Columns()
	m := dao.CatalogTorrent.Ctx(ctx).
		Where(columns.OwnerId, instance.UserId).
		Where(columns.Status, consts.CatalogTorrentStatusPublished).
		WhereGTE(columns.PublishedAt, instance.CycleStartedAt)
	if instance.CycleEndedAt != nil && !now.Before(instance.CycleEndedAt.Time) {
		m = m.WhereLT(columns.PublishedAt, instance.CycleEndedAt)
	}
	count, err := m.Count()
	return uint64(count), err
}

func (s *sSiteTaskUsecase) trackerProgress(ctx context.Context, instance entity.SiteUserTask, task model.SiteTaskDefinition, now time.Time) (uint64, error) {
	end := now
	if instance.CycleEndedAt != nil && !now.Before(instance.CycleEndedAt.Time) {
		end = instance.CycleEndedAt.Time.Add(-time.Second)
	}
	stats, err := service.AccountingTrafficDomain().QueryPeriodStats(
		ctx,
		instance.UserId,
		consts.AccountingStatPeriodDaily,
		instance.CycleStartedAt,
		gtime.NewFromTime(end),
	)
	if err != nil {
		return 0, err
	}
	var progress uint64
	for _, stat := range stats {
		if task.Rule.Type == consts.SiteTaskRuleTypeSeedDuration {
			progress += stat.SeedTime
		} else {
			progress += stat.RawUploaded
		}
	}
	return progress, nil
}

func (s *sSiteTaskUsecase) fulfillReward(ctx context.Context, taskId uint64, userId uint64, reward model.SiteTaskReward) error {
	switch reward.Type {
	case consts.SiteTaskRewardTypeBonus:
		return service.EconomyBonusUsecase().AddBonus(
			ctx,
			userId,
			reward.Amount,
			consts.EconomyBonusActionTaskReward,
			consts.EconomyBonusTargetTypeSiteUserTask,
			taskId,
			"",
			"",
		)
	case consts.SiteTaskRewardTypeVip:
		return service.IamUserDomain().ExtendUserVip(ctx, userId, int(reward.Amount), consts.SiteTaskVipRemark)
	case consts.SiteTaskRewardTypeInvite:
		for index := 0; index < int(reward.Amount); index++ {
			if _, err := service.IamInviteDomain().CreateInvite(ctx, userId, false, nil); err != nil {
				return err
			}
		}
		return nil
	default:
		return gerror.New("unsupported task reward type")
	}
}

func (s *sSiteTaskUsecase) loadTasksCache(ctx context.Context) (model.SiteTasks, error) {
	cacheKey := service.SysCache().KeySiteConfigFullPath(ctx, consts.SiteConfigSiteTasks)
	value, err := gcache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		return s.loadTasks(ctx)
	}, siteTaskConfigCacheDuration)
	if err != nil {
		return nil, err
	}
	if tasks, ok := value.Val().(model.SiteTasks); ok {
		return tasks, nil
	}
	return s.loadTasks(ctx)
}

func (s *sSiteTaskUsecase) loadTasks(ctx context.Context) (model.SiteTasks, error) {
	var tasks model.SiteTasks
	if err := service.SiteConfigDomain().GetByPath(ctx, consts.SiteConfigSiteTasks).Scan(&tasks); err != nil {
		return nil, gerror.New("invalid site tasks config")
	}
	tasks = tasks.Normalized()
	tasks = tasks.Supported()
	if err := tasks.Validate(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *sSiteTaskUsecase) userTaskItem(instance entity.SiteUserTask) siteout.UserTaskItem {
	return siteout.UserTaskItem{
		Id:             instance.Id,
		TaskKey:        instance.TaskKey,
		CycleKey:       instance.CycleKey,
		Status:         instance.Status,
		Progress:       instance.Progress,
		Target:         instance.Target,
		CycleStartedAt: instance.CycleStartedAt,
		CycleEndedAt:   instance.CycleEndedAt,
		ClaimedAt:      instance.ClaimedAt,
		CompletedAt:    instance.CompletedAt,
		RewardedAt:     instance.RewardedAt,
	}
}

func (s *sSiteTaskUsecase) taskTarget(task model.SiteTaskDefinition) uint64 {
	if task.Rule.Type == consts.SiteTaskRuleTypeSeedDuration {
		return task.Rule.Target * 3600
	}
	if task.Rule.Type == consts.SiteTaskRuleTypeUploaded {
		return task.Rule.Target * siteTaskGiB
	}
	return task.Rule.Target
}
