# Middleware (技术中间件) 域设计

## 定位
`Middleware` 并不是业务域（Business Domain），而是一个**横切关注点 (Cross-cutting Concern)**。
它负责将与业务核心无关的技术逻辑从 Controller / Usecase 中剥离出来。

## 核心中间件划分

### 1. AuthMiddleware (认证鉴权)
- **JwtParser**: 解析请求头中的 Token，并将解析出的 `user_id`、`role` 注入到当前请求的 `context.Context` 中。
- **PasskeyParser**: 专为 Tracker 域准备，拦截所有 `/announce` 和 `/scrape` 请求，从 URL Params 中提取 Passkey，校验后放行。

### 2. SecurityMiddleware (安全风控)
- **RateLimiter**: 接口限流器，防止恶意爬虫和高频访问。
- **Cors**: 跨域资源共享配置。

### 3. ToolingMiddleware (工具链)
- **I18nInjector**: 根据请求头（如 `Accept-Language`）自动设置多语言环境上下文。
- **ResponseWrapper**: 统一拦截成功的 Response 或抛出的 Error，包装为 `{code, message, data}` 标准 JSON 格式返回。
