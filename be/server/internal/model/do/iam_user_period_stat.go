// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserPeriodStat is the golang structure of table iam_user_period_stat for DAO operations like Where/Data.
type IamUserPeriodStat struct {
	g.Meta        `orm:"table:iam_user_period_stat, do:true"`
	Id            any         //
	UserId        any         //
	PeriodType    any         // 1=每日 2=每月
	PeriodKey     any         // YYYY-MM-DD 或 YYYY-MM
	Uploaded      any         // 周期新增入账上传量 (bytes)
	Downloaded    any         // 周期新增入账下载量 (bytes)
	RawUploaded   any         // 周期新增真实上传量 (bytes)
	RawDownloaded any         // 周期新增真实下载量 (bytes)
	SeedTime      any         // 周期新增做种时间 (秒)
	LeechTime     any         // 周期新增下载时间 (秒)
	Bonus         any         // 周期获得魔力值
	CreatedAt     *gtime.Time //
}
