package cmd

import (
	"context"

	"server/internal/router"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			if err := s.SetConfigWithMap(g.Map{
				"ClientMaxBodySize": "100mb",
			}); err != nil {
				return err
			}
			s.AddStaticPath("/api/public", "resource/public")

			s.Group("/api", func(group *ghttp.RouterGroup) {
				group.Middleware(
					service.Middleware().Ctx,
					service.Middleware().I18N,
					service.Middleware().ResponseHandler,
				)

				router.Accounting(ctx, group)
				router.Admin(ctx, group)
				router.Catalog(ctx, group)
				router.Economy(ctx, group)
				router.Forum(ctx, group)
				router.Tracker(ctx, group)
				router.Iam(ctx, group)
			})
			s.Run()
			return nil
		},
	}
)
