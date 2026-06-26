// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IMiddleware interface {
		// Ctx 初始化请求上下文
		Ctx(r *ghttp.Request)
		// I18N 国际化中间件
		I18N(r *ghttp.Request)
		// CORS 跨域处理
		CORS(r *ghttp.Request)
		// ResponseHandler 统一 API 返回格式封装
		ResponseHandler(r *ghttp.Request)
		// CheckAuth Web 端的普通用户鉴权拦截器 (委托给 auth 模块处理)
		CheckAuth(r *ghttp.Request)
		RBAC(r *ghttp.Request)
		// RequireStaff 管理员后台拦截器，强制要求必须是 Staff
		RequireStaff(r *ghttp.Request)
		// SetTrackerResponseType 标记当前请求需要返回 Tracker/Bencode 响应
		SetTrackerResponseType(r *ghttp.Request)
		// CheckTrackerAuth Tracker 的 Passkey 鉴权
		CheckTrackerAuth(r *ghttp.Request)
	}
)

var (
	localMiddleware IMiddleware
)

func Middleware() IMiddleware {
	if localMiddleware == nil {
		panic("implement not found for interface IMiddleware, forgot register?")
	}
	return localMiddleware
}

func RegisterMiddleware(i IMiddleware) {
	localMiddleware = i
}
