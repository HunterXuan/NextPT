# NextPT

NextPT 是一个正在开发中的私有 Tracker 与社区站点项目。后端采用 GoFrame，前端采用 Nuxt 4；当前已覆盖账号与权限、Tracker、种子与论坛、站点运营、魔力与用户任务等核心链路。

## 项目状态

> **开发中，尚未发布稳定版本。** 数据库结构、站点配置、接口和权限模型仍可能随开发演进而调整。项目目前不应被视为可直接落地的成品建站方案。

如计划用于生产环境，请由了解 Go、MySQL、Redis、Docker、反向代理与 BT Tracker 运维的团队自行评估并承担运行责任。至少应在上线前完成：

- 阅读并审查与自身部署方式相关的代码、配置及数据库初始化脚本。
- 在隔离环境完成完整功能、升级和故障恢复演练。
- 通过密钥管理或部署环境注入所有真实凭据，避免将其写入版本库或镜像。
- 建立并定期验证 MySQL、Redis 与对象存储的备份和恢复流程。
- 配置 HTTPS、反向代理、日志、健康检查、监控和告警，并按站点规模评估限流与容量。

## 当前能力

- IAM：注册与邀请、邮箱验证、密码找回、TOTP 两步验证、Passkey、登录记录与设备会话管理。
- 资源：种子发布与审核、分类和标签、字幕、评论、收藏、赞赏、举报、优惠与外部媒体元数据。
- Tracker：announce、scrape、私有种子下载、RSS、客户端白名单、流量统计和做种魔力结算。
- 社区：论坛、站内通知、站务信箱、在线聊天室、公告、用户公开资料卡。
- 运营：角色自动升降级、用户任务、魔力商城、广告位、审计日志、定时任务、维护模式与健康检查。

能力范围与后续方向见 [产品路线图](doc/roadmap/pt_maturity.md)。业务用例、接口边界和架构约定见 [文档索引](doc/README.md)。

## 目录

```text
be/server/   GoFrame 后端、数据库初始化脚本和服务端资源
fe/ui/       Nuxt 4 前端
doc/         架构规范、业务用例和路线图
.github/     镜像构建工作流
```

## 本地开发

前置条件：Go 1.25、Node.js 24、MySQL、Redis。对象存储、Mailgun 与第三方元数据服务可按需配置。

初始化空数据库：

```bash
mysql -u <user> -p <database> < be/server/storage/migration/init.sql
mysql -u <user> -p <database> < be/server/storage/migration/bootstrap.sql
```

`bootstrap.sql` 只用于本地开发：它会创建可预测的管理员账号和 Tracker Passkey，**不得原样用于生产环境**。生产部署应通过受控初始化流程创建唯一的首个 Staff 账户，并立即设置独立密码与 Passkey。

先从模板创建本地配置，并填写数据库及按需启用的对象存储、邮件和元数据服务凭据：

```bash
cp be/server/manifest/config/config.example.yaml be/server/manifest/config/config.yaml
```

`config.yaml` 已被 Git 忽略，不应提交。随后启动后端：

```bash
cd be/server
go run main.go
```

另开终端启动前端：

```bash
cd fe/ui
npm ci
npm run dev
```

开发模式下 Nuxt 会将 `/api/**` 代理到 `http://localhost:8000/api/**`。生产环境需要由网关或反向代理将同一路径转发至后端服务。

## 容器镜像

推送至 `main` 时，[镜像构建工作流](.github/workflows/build.yml) 会将前端和后端镜像发布到 GHCR，并同时提供 `linux/amd64` 与 `linux/arm64`。工作流只负责构建和推送，不包含生产部署；部署编排、密钥和运行时配置需要由实际运维环境提供。

安全问题请遵循 [安全披露说明](SECURITY.md)，不要在公开 Issue 中提交漏洞细节或疑似凭据。
