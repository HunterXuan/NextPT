package adminout

import (
	"github.com/gogf/gf/v2/os/gtime"
)

type SysCronLogItem struct {
	Id           uint64      `json:"id"`
	JobName      string      `json:"jobName"`
	NodeIp       string      `json:"nodeIp"`
	Status       int         `json:"status"` // 0: 执行中 1: 成功 2: 失败
	DurationMs   int         `json:"durationMs"`
	ErrorMessage string      `json:"errorMessage"`
	CreatedAt    *gtime.Time `json:"createdAt"`
	UpdatedAt    *gtime.Time `json:"updatedAt"`
}

type SysCronLogListOut struct {
	Page  int               `json:"page"`
	Size  int               `json:"size"`
	Total int               `json:"total"`
	List  []*SysCronLogItem `json:"list"`
}
