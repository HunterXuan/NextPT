// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteUserTask is the golang structure for table site_user_task.
type SiteUserTask struct {
	Id             uint64      `json:"id"             orm:"id"               description:""`
	UserId         uint64      `json:"userId"         orm:"user_id"          description:""`
	TaskKey        string      `json:"taskKey"        orm:"task_key"         description:""`
	CycleKey       string      `json:"cycleKey"       orm:"cycle_key"        description:"once、YYYY-Www 或 YYYY-MM"`
	CycleStartedAt *gtime.Time `json:"cycleStartedAt" orm:"cycle_started_at" description:""`
	CycleEndedAt   *gtime.Time `json:"cycleEndedAt"   orm:"cycle_ended_at"   description:""`
	Status         int         `json:"status"         orm:"status"           description:"0=active 1=completed 2=rewarded 3=expired"`
	Progress       uint64      `json:"progress"       orm:"progress"         description:""`
	Target         uint64      `json:"target"         orm:"target"           description:""`
	TaskSnapshot   *gjson.Json `json:"taskSnapshot"   orm:"task_snapshot"    description:""`
	ClaimedAt      *gtime.Time `json:"claimedAt"      orm:"claimed_at"       description:""`
	CompletedAt    *gtime.Time `json:"completedAt"    orm:"completed_at"     description:""`
	RewardedAt     *gtime.Time `json:"rewardedAt"     orm:"rewarded_at"      description:""`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:""`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       description:""`
}
