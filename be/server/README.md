# NextPT Server

NextPT 后端服务，基于 GoFrame。它提供 Web API、Tracker API、定时任务、健康检查与静态公共资源。

## 运行依赖

- Go 1.25
- MySQL
- Redis
- 可选：S3 兼容对象存储、Mailgun、TMDB / IMDb / 豆瓣 / Bangumi 元数据服务

## 本地启动

先使用 [初始化脚本](storage/migration/) 创建空数据库并写入最小 Bootstrap 数据，再按本地环境调整 `manifest/config/config.yaml`：

```bash
go run main.go
```

默认 HTTP 服务监听 `:8000`。接口统一位于 `/api`，探针位于根路径：

- `GET /healthz`：仅确认 HTTP 服务存活。
- `GET /readyz`：检查 MySQL 和 Redis 是否就绪。

## 配置与部署

`manifest/config/config.yaml` 是本地开发配置入口。生产环境应通过部署系统维护独立配置和凭据，不能复用开发凭据或将真实密钥提交到版本库。前端生产环境需要由反向代理把 `/api/**` 转发至本服务。

Dockerfile 位于 `manifest/docker/Dockerfile`，支持通过 Buildx 构建 AMD64 与 ARM64 镜像。镜像只包含服务运行所需的二进制、`manifest/i18n`、`resource` 与 `storage`；数据库、Redis、对象存储和入口网关需要由部署环境提供。
