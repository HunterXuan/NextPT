// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserStat is the golang structure for table iam_user_stat.
type IamUserStat struct {
	Id           uint64      `json:"id"           orm:"id"            description:""`
	UserId       uint64      `json:"userId"       orm:"user_id"       description:""`
	Uploaded     uint64      `json:"uploaded"     orm:"uploaded"      description:"总上传量 (bytes)"`
	Downloaded   uint64      `json:"downloaded"   orm:"downloaded"    description:"总下载量 (bytes)"`
	SeedTime     uint64      `json:"seedTime"     orm:"seed_time"     description:"总做种时间 (秒)"`
	LeechTime    uint64      `json:"leechTime"    orm:"leech_time"    description:"总下载时间 (秒)"`
	Bonus        float64     `json:"bonus"        orm:"bonus"         description:"魔力值"`
	BonusCharity float64     `json:"bonusCharity" orm:"bonus_charity" description:"捐赠魔力值"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:""`
}
