package router

import (
	"context"

	"server/internal/controller/catalog"
	"server/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Catalog(_ context.Context, group *ghttp.RouterGroup) {
	group.Group("/catalog", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().CheckAuth, service.Middleware().RBAC)
		group.Bind(
			catalog.NewV1(),
		)
	})
}
