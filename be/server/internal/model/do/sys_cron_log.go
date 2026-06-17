// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysCronLog is the golang structure of table sys_cron_log for DAO operations like Where/Data.
type SysCronLog struct {
	g.Meta       `orm:"table:sys_cron_log, do:true"`
	Id           any         //
	JobName      any         // 任务名称
	NodeIp       any         // 执行节点IP
	Status       any         // 状态 0: 执行中 1: 成功 2: 失败
	DurationMs   any         // 执行耗时(毫秒)
	ErrorMessage any         // 错误信息
	CreatedAt    *gtime.Time // 开始时间
	UpdatedAt    *gtime.Time // 更新时间
}
