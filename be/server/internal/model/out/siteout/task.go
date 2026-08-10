package siteout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
)

type TaskItem struct {
	Key             string                 `json:"key"`
	Enabled         bool                   `json:"enabled"`
	Cycle           string                 `json:"cycle"`
	NameI18N        map[string]string      `json:"nameI18n"`
	DescriptionI18N map[string]string      `json:"descriptionI18n"`
	Rule            model.SiteTaskRule     `json:"rule"`
	Rewards         []model.SiteTaskReward `json:"rewards"`
	UserTask        *UserTaskItem          `json:"userTask,omitempty"`
}

type UserTaskItem struct {
	Id             uint64      `json:"id"`
	TaskKey        string      `json:"taskKey"`
	CycleKey       string      `json:"cycleKey"`
	Status         int         `json:"status"`
	Progress       uint64      `json:"progress"`
	Target         uint64      `json:"target"`
	CycleStartedAt *gtime.Time `json:"cycleStartedAt"`
	CycleEndedAt   *gtime.Time `json:"cycleEndedAt"`
	ClaimedAt      *gtime.Time `json:"claimedAt"`
	CompletedAt    *gtime.Time `json:"completedAt"`
	RewardedAt     *gtime.Time `json:"rewardedAt"`
}

type TaskListOut struct {
	List []*TaskItem `json:"list"`
}

type TaskClaimOut struct {
	UserTask UserTaskItem `json:"userTask"`
}
