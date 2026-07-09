package router

import (
	"context"

	"server/internal/controller/site"
	"server/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Site(_ context.Context, group *ghttp.RouterGroup) {
	group.Group("/site", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().CheckAuth, service.Middleware().RBAC)
		group.Bind(
			site.NewV1(),
		)
	})
}
