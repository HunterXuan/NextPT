package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type SysCronListReq struct {
	g.Meta `path:"/sys/crons" method:"get" tags:"AdminSys" summary:"获取系统定时任务列表" perm:"admin:sys/cron:*"`
	adminin.SysCronListInp
}

type SysCronListRes struct {
	adminout.SysCronListOut
}
