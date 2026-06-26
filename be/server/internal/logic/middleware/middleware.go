package middleware

import (
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/util/gconv"

	"server/internal/library/contexts"
	"server/internal/model"
	"server/internal/service"

	"github.com/anacrolix/torrent/bencode"
	"github.com/goflyfox/gtoken/v2/gtoken"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

type sMiddleware struct{}

func NewMiddleware() *sMiddleware {
	return &sMiddleware{}
}

func init() {
	service.RegisterMiddleware(NewMiddleware())
}

// =========================================================================
// 1. 全局与上下文中间件
// =========================================================================

// Ctx 初始化请求上下文
func (s *sMiddleware) Ctx(r *ghttp.Request) {
	contexts.Init(r, &model.Context{})

	if len(r.Cookie.GetSessionId()) == 0 {
		r.Cookie.SetSessionId(gctx.CtxId(r.Context()))
	}

	r.SetCtx(r.GetNeverDoneCtx())
	r.Middleware.Next()
}

// I18N 国际化中间件
func (s *sMiddleware) I18N(r *ghttp.Request) {
	lang := r.Get("lang").String()
	if lang == "" {
		// 提取 Accept-Language 的首个语言项 (e.g. "zh-CN,zh;q=0.9" -> "zh-CN")
		al := r.Header.Get("Accept-Language")
		if al != "" {
			lang = strings.Split(al, ",")[0]
			lang = strings.Split(lang, ";")[0]
			lang = strings.TrimSpace(lang)
		}
	}

	// 从配置中获取支持的语言列表和默认语言
	supportedLangsVar, _ := g.Cfg().Get(r.Context(), "i18n.lang")
	defaultLangVar, _ := g.Cfg().Get(r.Context(), "i18n.default")

	defaultLang := defaultLangVar.String()
	if defaultLang == "" {
		defaultLang = "en-US"
	}

	isSupported := false
	if !supportedLangsVar.IsEmpty() {
		// 检查当前请求的语言是否在配置的 key 中
		if _, ok := supportedLangsVar.Map()[lang]; ok {
			isSupported = true
		}
	}

	// 如果未匹配到支持的语言，兜底到默认语言
	if !isSupported {
		lang = defaultLang
	}

	ctx := gi18n.WithLanguage(r.Context(), lang)
	r.SetCtx(ctx)

	r.Middleware.Next()
}

// CORS 跨域处理
func (s *sMiddleware) CORS(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}

// ResponseHandler 统一 API 返回格式封装
func (s *sMiddleware) ResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()

	// 如果已经被处理过（比如 Tracker 直接写入了 Bencode 并且 ExitAll），直接返回
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		err = r.GetError()
		res = r.GetHandlerResponse()
	)

	// 统一 JSON 格式返回
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    1,
			Message: err.Error(),
			Data:    nil,
		})
	} else {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    0,
			Message: "success",
			Data:    res,
		})
	}
}

// =========================================================================
// 2. Web 鉴权中间件 (gtoken 负责核心拦截，这里只做权限补充)
// =========================================================================

// CheckAuth Web 端的普通用户鉴权拦截器 (委托给 auth 模块处理)
func (s *sMiddleware) CheckAuth(r *ghttp.Request) {
	handler := r.GetServeHandler()
	if handler != nil && handler.GetMetaTag("noAuth") == "true" {
		r.Middleware.Next()
		return
	}

	token, err := gtoken.GetRequestToken(r)
	if err != nil {
		service.IamSessionDomain().GetGFMiddleware().ResFun(r, err)
		return
	}

	userKey, err := service.IamSessionDomain().GetGFToken().Validate(r.Context(), token)
	if err != nil {
		service.IamSessionDomain().GetGFMiddleware().ResFun(r, err)
		return
	}
	r.SetCtxVar(gtoken.KeyUserKey, userKey)

	actor, err := service.IamUserUsecase().LoadActor(r.Context(), gconv.Uint64(userKey))
	if err != nil || actor == nil {
		_ = service.IamSessionDomain().RemoveToken(r.Context(), userKey)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    401,
			Message: gi18n.T(r.Context(), "iam.general.unauthorized"),
			Data:    nil,
		})
		r.ExitAll()
		return
	}
	contexts.SetActor(r.Context(), actor)

	r.Middleware.Next()
}

func (s *sMiddleware) RBAC(r *ghttp.Request) {
	ctx := r.GetCtx()
	handler := r.GetServeHandler()
	if handler == nil {
		r.Middleware.Next()
		return
	}
	if handler.GetMetaTag("noAuth") == "true" {
		r.Middleware.Next()
		return
	}

	user := contexts.GetActor(ctx)
	if user == nil {
		s.writeRBACError(r, 401, gi18n.T(ctx, "iam.general.unauthorized"), "Invalid passkey or account disabled")
		return
	}

	permTpl := handler.GetMetaTag("perm")
	if permTpl == "" {
		if handler.GetMetaTag("noPerm") != "true" && strings.HasPrefix(r.Router.Uri, "/api/admin") {
			s.writeRBACError(r, 403, gi18n.T(ctx, "iam.general.no_permission"), "Permission denied")
			return
		}
		r.Middleware.Next()
		return
	}

	// Dynamic parse {xxx} to parameters
	re := regexp.MustCompile(`\{([^}]+)\}`)
	permKey := re.ReplaceAllStringFunc(permTpl, func(match string) string {
		paramName := match[1 : len(match)-1]
		return r.Get(paramName).String()
	})

	hasPerm, err := service.IamUserUsecase().CheckPermission(ctx, user, permKey)
	if err != nil || !hasPerm {
		s.writeRBACError(r, 403, gi18n.T(ctx, "iam.general.no_permission"), s.trackerPermissionMessage(permKey))
		return
	}

	r.Middleware.Next()
}

func (s *sMiddleware) writeRBACError(r *ghttp.Request, code int, webMessage, trackerMessage string) {
	if s.isTrackerRequest(r) {
		s.writeBencodeError(r, trackerMessage)
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    code,
		Message: webMessage,
		Data:    nil,
	})
	r.ExitAll()
}

func (s *sMiddleware) trackerPermissionMessage(permKey string) string {
	if strings.HasPrefix(permKey, "download:") {
		return "user is not allowed to download"
	}
	return "Permission denied"
}

func (s *sMiddleware) isTrackerRequest(r *ghttp.Request) bool {
	return contexts.GetResponseType(r.Context()) == model.ResponseTypeTracker
}

// RequireStaff 管理员后台拦截器，强制要求必须是 Staff
func (s *sMiddleware) RequireStaff(r *ghttp.Request) {
	ctx := r.GetCtx()
	user := contexts.GetActor(ctx)
	if user == nil || !user.IsStaff {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    403,
			Message: gi18n.T(ctx, "iam.general.no_permission"),
			Data:    nil,
		})
		r.ExitAll()
		return
	}
	r.Middleware.Next()
}

// =========================================================================
// 3. Tracker 专用中间件 (必须返回 Bencode)
// =========================================================================

// SetTrackerResponseType 标记当前请求需要返回 Tracker/Bencode 响应
func (s *sMiddleware) SetTrackerResponseType(r *ghttp.Request) {
	contexts.SetResponseType(r.Context(), model.ResponseTypeTracker)
	r.Middleware.Next()
}

// CheckTrackerAuth Tracker 的 Passkey 鉴权
func (s *sMiddleware) CheckTrackerAuth(r *ghttp.Request) {
	passkey := r.GetQuery("passkey").String()
	if passkey == "" || len(passkey) != 32 {
		s.writeBencodeError(r, "Invalid passkey")
		return
	}

	actor, err := service.IamSessionUsecase().VerifyPasskey(r.Context(), passkey)
	if err != nil || actor == nil {
		s.writeBencodeError(r, "Invalid passkey or account disabled")
		return
	}

	contexts.SetActor(r.Context(), actor)

	r.Middleware.Next()
}

// writeBencodeError 写入标准 BitTorrent 错误格式并中断请求
func (s *sMiddleware) writeBencodeError(r *ghttp.Request, msg string) {
	r.Response.ClearBuffer()
	r.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")

	m := map[string]string{"failure reason": msg}
	b, _ := bencode.Marshal(m)
	r.Response.Write(b)
	r.ExitAll()
}
