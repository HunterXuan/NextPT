package contexts

import (
	"context"

	"server/internal/consts"
	"server/internal/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改
func Init(r *ghttp.Request, customCtx *model.Context) {
	r.SetCtxVar(consts.SysContextHTTPKey, customCtx)
}

// Get 获得上下文变量，如果没有设置，那么返回nil
func Get(ctx context.Context) *model.Context {
	value := ctx.Value(consts.SysContextHTTPKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.Context); ok {
		return localCtx
	}
	return nil
}

// SetActor 将上下文信息设置到上下文请求中，注意是完整覆盖
func SetActor(ctx context.Context, actor *model.Actor) {
	c := Get(ctx)
	if c == nil {
		g.Log().Warning(ctx, "contexts.SetActor, c == nil ")
		return
	}
	c.Actor = actor
}

func SetSessionId(ctx context.Context, sessionId string) {
	c := Get(ctx)
	if c == nil {
		g.Log().Warning(ctx, "contexts.SetSessionId, c == nil ")
		return
	}
	c.SessionId = sessionId
}

func SetResponseType(ctx context.Context, responseType model.ResponseType) {
	c := Get(ctx)
	if c == nil {
		g.Log().Warning(ctx, "contexts.SetResponseType, c == nil ")
		return
	}
	c.ResponseType = responseType
}

func GetResponseType(ctx context.Context) model.ResponseType {
	c := Get(ctx)
	if c == nil {
		return model.ResponseTypeJSON
	}
	if c.ResponseType == "" {
		return model.ResponseTypeJSON
	}
	return c.ResponseType
}

// GetActor 获取操作者信息
func GetActor(ctx context.Context) *model.Actor {
	c := Get(ctx)
	if c == nil {
		return nil
	}
	return c.Actor
}

func GetSessionId(ctx context.Context) string {
	c := Get(ctx)
	if c == nil {
		return ""
	}
	return c.SessionId
}

// GetUserId 获取用户ID
func GetUserId(ctx context.Context) uint64 {
	actor := GetActor(ctx)
	if actor == nil {
		return 0
	}
	return actor.Id
}

func GetUserRoleLevel(ctx context.Context) int {
	actor := GetActor(ctx)
	if actor == nil {
		return 0
	}
	return actor.RoleLevel
}
