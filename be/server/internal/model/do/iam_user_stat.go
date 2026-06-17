// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserStat is the golang structure of table iam_user_stat for DAO operations like Where/Data.
type IamUserStat struct {
	g.Meta       `orm:"table:iam_user_stat, do:true"`
	Id           any         //
	UserId       any         //
	Uploaded     any         // 总上传量 (bytes)
	Downloaded   any         // 总下载量 (bytes)
	SeedTime     any         // 总做种时间 (秒)
	LeechTime    any         // 总下载时间 (秒)
	Bonus        any         // 魔力值
	BonusCharity any         // 捐赠魔力值
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
