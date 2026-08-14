package router

import (
	"context"

	"server/internal/controller/mod"
	"server/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Mod(_ context.Context, group *ghttp.RouterGroup) {
	group.Group("/mod", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().CheckAuth, service.Middleware().RBAC)
		group.Bind(mod.NewV1())
	})
}
