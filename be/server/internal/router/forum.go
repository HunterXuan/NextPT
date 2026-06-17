package router

import (
	"context"

	"server/internal/controller/forum"
	"server/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Forum(_ context.Context, group *ghttp.RouterGroup) {
	group.Group("/forum", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().CheckAuth, service.Middleware().RBAC)
		group.Bind(
			forum.NewV1(),
		)
	})
}
