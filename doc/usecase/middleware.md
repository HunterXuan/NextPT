# Middleware (技术中间件) 域设计

## 定位
`Middleware` 并不是业务域（Business Domain），而是一个**横切关注点 (Cross-cutting Concern)**。
它负责将与业务核心无关的技术逻辑从 Controller / Usecase 中剥离出来。

## 核心中间件划分

### 1. AuthMiddleware (认证鉴权)
- **CheckAuth**: 校验 Web Session Token 与 `X-Device-Id` 摘要的一致性，将 `model.Actor` 写入请求上下文；Controller 再显式向 Usecase 传递 Actor。
- **CheckTrackerAuth**: 专为 Tracker 域准备，从请求参数提取 Passkey，校验账户状态后建立 Actor。它覆盖 announce、scrape、RSS 与种子文件下载等 Tracker 路由。
- **RBAC**: 同时处理 Web 和 Tracker 请求的通配权限校验；资源级权限仍由业务 Usecase 处理 owner 与状态等约束。

### 2. SecurityMiddleware (安全风控)
- **Maintenance**: 维护模式下阻断普通用户 API 请求，同时保留 Staff、Tracker、`/healthz` 与 `/readyz`。
- 通用限流与 CORS 当前由部署入口或后续网关配置承担，不在本服务中声明独立中间件。

### 3. ToolingMiddleware (工具链)
- **I18N**: 根据请求语言设置多语言上下文。
- **ResponseHandler**: 将普通 Web API 的成功和错误响应包装为 `{code, message, data}`；Tracker 协议响应保持 Bencode 格式。
