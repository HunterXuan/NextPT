package global

import (
	"context"

	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func Init(ctx context.Context) {
	if err := gtime.SetTimeZone("Asia/Shanghai"); err != nil {
		g.Log().Fatalf(ctx, "SetTimeZone fail: %+v", err)
	}

	// 启动所有系统定时任务
	service.SysCron().Start(ctx)
}
