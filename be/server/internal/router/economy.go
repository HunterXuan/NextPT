package router

import (
	"context"

	"server/internal/controller/economy"
	"server/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Economy(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/economy", func(group *ghttp.RouterGroup) {
		// 中间件绑定：需要登录
		group.Middleware(service.Middleware().CheckAuth, service.Middleware().RBAC)
		// 绑定 Economy API
		group.Bind(
			economy.NewV1(),
		)
	})
}
