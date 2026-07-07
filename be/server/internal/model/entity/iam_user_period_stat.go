// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserPeriodStat is the golang structure for table iam_user_period_stat.
type IamUserPeriodStat struct {
	Id            uint64      `json:"id"            orm:"id"             description:""`
	UserId        uint64      `json:"userId"        orm:"user_id"        description:""`
	PeriodType    uint        `json:"periodType"    orm:"period_type"    description:"1=每日 2=每月"`
	PeriodKey     string      `json:"periodKey"     orm:"period_key"     description:"YYYY-MM-DD 或 YYYY-MM"`
	Uploaded      uint64      `json:"uploaded"      orm:"uploaded"       description:"周期新增入账上传量 (bytes)"`
	Downloaded    uint64      `json:"downloaded"    orm:"downloaded"     description:"周期新增入账下载量 (bytes)"`
	RawUploaded   uint64      `json:"rawUploaded"   orm:"raw_uploaded"   description:"周期新增真实上传量 (bytes)"`
	RawDownloaded uint64      `json:"rawDownloaded" orm:"raw_downloaded" description:"周期新增真实下载量 (bytes)"`
	SeedTime      uint64      `json:"seedTime"      orm:"seed_time"      description:"周期新增做种时间 (秒)"`
	LeechTime     uint64      `json:"leechTime"     orm:"leech_time"     description:"周期新增下载时间 (秒)"`
	Bonus         float64     `json:"bonus"         orm:"bonus"          description:"周期获得魔力值"`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`
}
