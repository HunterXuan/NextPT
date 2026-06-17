package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

// SysCronLogListReq 获取定时任务执行日志列表
type SysCronLogListReq struct {
	g.Meta `path:"/sys/crons/{name}/logs" method:"get" tags:"AdminSys" summary:"获取系统定时任务日志" perm:"admin:sys/cron:*"`
	adminin.SysCronLogListInp
}

type SysCronLogListRes struct {
	adminout.SysCronLogListOut
}
