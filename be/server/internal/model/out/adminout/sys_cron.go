package adminout

import "github.com/gogf/gf/v2/os/gtime"

type SysCronItem struct {
	Name         string      `json:"name"`
	Status       int         `json:"status"` // 0: runnable, 1: running, 2: stopped, -1: closed
	RegisterTime *gtime.Time `json:"registerTime"`
}

type SysCronListOut struct {
	List []*SysCronItem `json:"list"`
}
