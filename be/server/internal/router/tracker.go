package router

import (
	"context"

	"server/internal/controller/tracker"
	"server/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Tracker(_ context.Context, group *ghttp.RouterGroup) {
	group.Group("/tracker", func(group *ghttp.RouterGroup) {
		// BT 客户端使用 Passkey 鉴权，不要使用 Web JWT (CheckAuth)
		group.Middleware(
			service.Middleware().SetTrackerResponseType,
			service.Middleware().CheckTrackerAuth,
			service.Middleware().RBAC,
		)
		group.Bind(tracker.NewV1())
	})
}
