// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserDailyStat is the golang structure for table iam_user_daily_stat.
type IamUserDailyStat struct {
	Id         uint64      `json:"id"         orm:"id"         description:""`
	UserId     uint64      `json:"userId"     orm:"user_id"    description:""`
	Date       *gtime.Time `json:"date"       orm:"date"       description:"统计日期 YYYY-MM-DD"`
	Uploaded   uint64      `json:"uploaded"   orm:"uploaded"   description:"当日新增上传量 (bytes)"`
	Downloaded uint64      `json:"downloaded" orm:"downloaded" description:"当日新增下载量 (bytes)"`
	SeedTime   uint64      `json:"seedTime"   orm:"seed_time"  description:"当日新增做种时间 (秒)"`
	LeechTime  uint64      `json:"leechTime"  orm:"leech_time" description:"当日新增下载时间 (秒)"`
	Bonus      float64     `json:"bonus"      orm:"bonus"      description:"当日获得魔力值"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at" description:""`
}
