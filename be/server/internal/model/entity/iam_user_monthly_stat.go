// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserMonthlyStat is the golang structure for table iam_user_monthly_stat.
type IamUserMonthlyStat struct {
	Id         uint64      `json:"id"         orm:"id"         description:""`
	UserId     uint64      `json:"userId"     orm:"user_id"    description:""`
	YearMonth  string      `json:"yearMonth"  orm:"year_month" description:"YYYY-MM"`
	Uploaded   uint64      `json:"uploaded"   orm:"uploaded"   description:"当月新增上传量 (bytes)"`
	Downloaded uint64      `json:"downloaded" orm:"downloaded" description:"当月新增下载量 (bytes)"`
	SeedTime   uint64      `json:"seedTime"   orm:"seed_time"  description:"当月新增做种时间 (秒)"`
	LeechTime  uint64      `json:"leechTime"  orm:"leech_time" description:"当月新增下载时间 (秒)"`
	Bonus      float64     `json:"bonus"      orm:"bonus"      description:"当月获得魔力值"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at" description:""`
}
