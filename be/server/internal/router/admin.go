package router

import (
	"context"

	"server/internal/controller/admin"
	"server/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Admin(_ context.Context, group *ghttp.RouterGroup) {
	group.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().CheckAuth, service.Middleware().RequireStaff, service.Middleware().RBAC)
		group.Bind(
			admin.NewV1(),
		)
	})
}
