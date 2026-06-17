// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserDailyStat is the golang structure of table iam_user_daily_stat for DAO operations like Where/Data.
type IamUserDailyStat struct {
	g.Meta     `orm:"table:iam_user_daily_stat, do:true"`
	Id         any         //
	UserId     any         //
	Date       *gtime.Time // 统计日期 YYYY-MM-DD
	Uploaded   any         // 当日新增上传量 (bytes)
	Downloaded any         // 当日新增下载量 (bytes)
	SeedTime   any         // 当日新增做种时间 (秒)
	LeechTime  any         // 当日新增下载时间 (秒)
	Bonus      any         // 当日获得魔力值
	CreatedAt  *gtime.Time //
}
