package main

import (
	"server/internal/cmd"
	"server/internal/global"
	_ "server/internal/logic"
	_ "server/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	ctx := gctx.GetInitCtx()
	global.Init(ctx)
	cmd.Main.Run(ctx)
}
