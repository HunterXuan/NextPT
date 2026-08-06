package router

import (
	"context"
	"net/http"

	"server/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Health registers operational probes outside the application API middleware stack.
func Health(_ context.Context, server *ghttp.Server) {
	server.BindHandler("GET:/healthz", func(r *ghttp.Request) {
		r.Response.WriteHeader(http.StatusNoContent)
	})

	server.BindHandler("GET:/readyz", func(r *ghttp.Request) {
		readiness := service.SysHealth().Readiness(r.Context())
		status := http.StatusOK
		if !readiness.IsReady() {
			status = http.StatusServiceUnavailable
		}
		r.Response.Header().Set("Content-Type", "application/json")
		r.Response.WriteHeader(status)
		r.Response.WriteJson(readiness)
	})
}
