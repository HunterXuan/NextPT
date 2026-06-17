package router

import (
	"context"

	"server/internal/controller/iam"
	"server/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Iam(_ context.Context, group *ghttp.RouterGroup) {
	group.Group("/iam", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().CheckAuth, service.Middleware().RBAC)
		group.Bind(
			iam.NewV1(),
		)
	})
}
