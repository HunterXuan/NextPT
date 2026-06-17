// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysCronLog is the golang structure for table sys_cron_log.
type SysCronLog struct {
	Id           uint64      `json:"id"           orm:"id"            description:""`
	JobName      string      `json:"jobName"      orm:"job_name"      description:"任务名称"`
	NodeIp       string      `json:"nodeIp"       orm:"node_ip"       description:"执行节点IP"`
	Status       int         `json:"status"       orm:"status"        description:"状态 0: 执行中 1: 成功 2: 失败"`
	DurationMs   int         `json:"durationMs"   orm:"duration_ms"   description:"执行耗时(毫秒)"`
	ErrorMessage string      `json:"errorMessage" orm:"error_message" description:"错误信息"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"开始时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"更新时间"`
}
