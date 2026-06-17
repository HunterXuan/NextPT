// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserMonthlyStat is the golang structure of table iam_user_monthly_stat for DAO operations like Where/Data.
type IamUserMonthlyStat struct {
	g.Meta     `orm:"table:iam_user_monthly_stat, do:true"`
	Id         any         //
	UserId     any         //
	YearMonth  any         // YYYY-MM
	Uploaded   any         // 当月新增上传量 (bytes)
	Downloaded any         // 当月新增下载量 (bytes)
	SeedTime   any         // 当月新增做种时间 (秒)
	LeechTime  any         // 当月新增下载时间 (秒)
	Bonus      any         // 当月获得魔力值
	CreatedAt  *gtime.Time //
}
