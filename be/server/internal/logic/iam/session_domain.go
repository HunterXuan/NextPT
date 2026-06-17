package iam

import (
	"context"

	"server/internal/service"

	"github.com/goflyfox/gtoken/v2/gtoken"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type sIamSessionDomain struct {
	gfToken    gtoken.Token
	middleware gtoken.Middleware
}

func init() {
	service.RegisterIamSessionDomain(NewIamSessionDomain())
}

func NewIamSessionDomain() *sIamSessionDomain {
	ctx := context.TODO()

	cacheMode := g.Cfg().MustGet(ctx, "iam.cacheMode", gtoken.CacheModeCache).Int8()
	cachePreKey := g.Cfg().MustGet(ctx, "iam.cachePreKey", "NextPT_Auth:").String()
	timeoutDur := g.Cfg().MustGet(ctx, "iam.timeout", "10d").Duration()
	maxRefreshDur := g.Cfg().MustGet(ctx, "iam.maxRefresh", "5d").Duration()
	multiLogin := g.Cfg().MustGet(ctx, "iam.multiLogin", true).Bool()
	tokenDelimiter := g.Cfg().MustGet(ctx, "iam.tokenDelimiter", "_").String()

	tm := gtoken.NewDefaultToken(gtoken.Options{
		CacheMode:      cacheMode,
		CachePreKey:    cachePreKey,
		Timeout:        timeoutDur.Milliseconds(),
		MaxRefresh:     maxRefreshDur.Milliseconds(),
		MultiLogin:     multiLogin,
		TokenDelimiter: tokenDelimiter,
	})

	mw := gtoken.NewDefaultMiddleware(tm)
	mw.ResFun = func(r *ghttp.Request, err error) {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    401,
			Message: "Unauthorized: " + err.Error(),
			Data:    nil,
		})
		r.ExitAll()
	}

	return &sIamSessionDomain{
		gfToken:    tm,
		middleware: mw,
	}
}

func (s *sIamSessionDomain) GetGFToken() gtoken.Token {
	return s.gfToken
}

func (s *sIamSessionDomain) GetGFMiddleware() gtoken.Middleware {
	return s.middleware
}

func (s *sIamSessionDomain) GenerateToken(ctx context.Context, userKey string, data any) (string, error) {
	return s.gfToken.Generate(ctx, userKey, data)
}

func (s *sIamSessionDomain) RemoveToken(ctx context.Context, userKey string) error {
	return s.gfToken.Destroy(ctx, userKey)
}
