// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SiteUserTask is the golang structure of table site_user_task for DAO operations like Where/Data.
type SiteUserTask struct {
	g.Meta         `orm:"table:site_user_task, do:true"`
	Id             any         //
	UserId         any         //
	TaskKey        any         //
	CycleKey       any         // once、YYYY-Www 或 YYYY-MM
	CycleStartedAt *gtime.Time //
	CycleEndedAt   *gtime.Time //
	Status         any         // 0=active 1=completed 2=rewarded 3=expired
	Progress       any         //
	Target         any         //
	TaskSnapshot   *gjson.Json //
	ClaimedAt      *gtime.Time //
	CompletedAt    *gtime.Time //
	RewardedAt     *gtime.Time //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
