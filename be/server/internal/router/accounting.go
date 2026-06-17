package router

import (
	"context"

	"server/internal/controller/accounting"
	"server/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Accounting(_ context.Context, group *ghttp.RouterGroup) {
	group.Group("/accounting", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().CheckAuth, service.Middleware().RBAC)
		group.Bind(
			accounting.NewV1(),
		)
	})
}
