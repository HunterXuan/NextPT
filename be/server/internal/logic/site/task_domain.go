package site

import (
	"context"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sSiteTaskDomain struct{}

func init() {
	service.RegisterSiteTaskDomain(NewSiteTaskDomain())
}

func NewSiteTaskDomain() *sSiteTaskDomain {
	return &sSiteTaskDomain{}
}

func (s *sSiteTaskDomain) ListUserTasks(ctx context.Context, userId uint64) ([]entity.SiteUserTask, error) {
	columns := dao.SiteUserTask.Columns()
	var tasks []entity.SiteUserTask
	err := dao.SiteUserTask.Ctx(ctx).
		Where(columns.UserId, userId).
		OrderDesc(columns.UpdatedAt).
		Scan(&tasks)
	return tasks, err
}

func (s *sSiteTaskDomain) CreateUserTask(ctx context.Context, task entity.SiteUserTask) (uint64, error) {
	id, err := dao.SiteUserTask.Ctx(ctx).Data(task).InsertAndGetId()
	return uint64(id), err
}

func (s *sSiteTaskDomain) ListActiveUserTasks(ctx context.Context, afterId uint64, size int) ([]entity.SiteUserTask, error) {
	columns := dao.SiteUserTask.Columns()
	m := dao.SiteUserTask.Ctx(ctx).
		Where(columns.Status, consts.SiteTaskStatusActive).
		OrderAsc(columns.Id).
		Limit(size)
	if afterId > 0 {
		m = m.WhereGT(columns.Id, afterId)
	}
	var tasks []entity.SiteUserTask
	err := m.Scan(&tasks)
	return tasks, err
}

func (s *sSiteTaskDomain) UpdateProgress(ctx context.Context, id uint64, progress uint64) error {
	columns := dao.SiteUserTask.Columns()
	_, err := dao.SiteUserTask.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.Status, consts.SiteTaskStatusActive).
		Data(g.Map{columns.Progress: progress, columns.UpdatedAt: gtime.Now()}).
		Update()
	return err
}

func (s *sSiteTaskDomain) Complete(ctx context.Context, id uint64, progress uint64) (bool, error) {
	columns := dao.SiteUserTask.Columns()
	result, err := dao.SiteUserTask.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.Status, consts.SiteTaskStatusActive).
		Data(g.Map{
			columns.Progress:    progress,
			columns.Status:      consts.SiteTaskStatusCompleted,
			columns.CompletedAt: gtime.Now(),
			columns.UpdatedAt:   gtime.Now(),
		}).
		Update()
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	return rows == 1, err
}

func (s *sSiteTaskDomain) Expire(ctx context.Context, id uint64) (bool, error) {
	columns := dao.SiteUserTask.Columns()
	result, err := dao.SiteUserTask.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.Status, consts.SiteTaskStatusActive).
		Data(g.Map{columns.Status: consts.SiteTaskStatusExpired, columns.UpdatedAt: gtime.Now()}).
		Update()
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	return rows == 1, err
}

func (s *sSiteTaskDomain) GetUserTask(ctx context.Context, id uint64, userId uint64) (*entity.SiteUserTask, error) {
	columns := dao.SiteUserTask.Columns()
	var task entity.SiteUserTask
	err := dao.SiteUserTask.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.UserId, userId).
		Scan(&task)
	if err != nil || task.Id == 0 {
		return nil, err
	}
	return &task, nil
}

func (s *sSiteTaskDomain) Reward(ctx context.Context, id uint64) (bool, error) {
	columns := dao.SiteUserTask.Columns()
	result, err := dao.SiteUserTask.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.Status, consts.SiteTaskStatusCompleted).
		Data(g.Map{
			columns.Status:     consts.SiteTaskStatusRewarded,
			columns.RewardedAt: gtime.Now(),
			columns.UpdatedAt:  gtime.Now(),
		}).
		Update()
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	return rows == 1, err
}

func (s *sSiteTaskDomain) CleanupHistory(ctx context.Context, before *gtime.Time) (int, error) {
	columns := dao.SiteUserTask.Columns()
	result, err := dao.SiteUserTask.Ctx(ctx).
		WhereIn(columns.Status, []int{consts.SiteTaskStatusRewarded, consts.SiteTaskStatusExpired}).
		WhereNot(columns.CycleKey, consts.SiteTaskCycleOnce).
		WhereLT(columns.UpdatedAt, before).
		Delete()
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return int(rows), err
}
