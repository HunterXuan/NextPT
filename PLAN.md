# NextPT 当前开发计划

这份文件记录当前阶段的工程关注点。数据库结构以 `be/server/storage/migration/init.sql` 为准，空库初始数据以 `be/server/storage/migration/bootstrap.sql` 为准；业务能力与产品方向见 `doc/roadmap/pt_maturity.md`。

## 当前状态

- IAM、RBAC、Tracker、Catalog、Forum、Economy、Site、Mod、Admin 和 Sys 主链路均已接入前后端。
- 已具备邮箱验证、密码找回、两步验证、Passkey、会话设备管理和新设备登录通知；默认注册角色由 `iam.default_register_role` 配置决定。
- 已具备种子审核、外部媒体元数据、RSS、求种与续种、用户任务、魔力商城、公告与通知、站务信箱、在线聊天室和广告位。
- 初始化脚本保持两层：`init.sql` 建表，`bootstrap.sql` 写入本地开发所需的最小角色、权限、分类、节点和站点配置。

## 工程约定

- API 统一挂在 `/api` 下，业务分组包括 `/iam`、`/catalog`、`/forum`、`/economy`、`/site`、`/mod` 和 `/admin`。
- 普通 JSON API 使用 `code`、`message`、`data` 响应包；下载和 Tracker 协议响应不使用该包装。
- `/api/iam/users/me` 是登录态主数据源，`/api/iam/users/me/permissions` 用于获取前端权限判断所需的通配权限。
- 后端 RBAC 是最终授权依据。前端只负责入口提示和交互限制，不能作为安全边界。
- Usecase / Domain 分层、事务边界和参数传递方式遵循 `doc/dev/architecture_guidelines.md`。

## 当前工程重点

- 补齐 Tracker、权限、站点配置、定时任务和数据库初始化脚本的关键回归测试。
- 收紧分页、筛选与配置输入校验，并持续检查跨域缓存失效和权限缓存一致性。
- 打磨生产部署方式：配置分离、凭据管理、备份恢复、可观测性、镜像升级与回滚。
- 根据实际运营数据再决定 RSS 自定义订阅、H&R、MediaInfo / 截图、Staff 风险画像和旧站迁移的优先级。

## 暂缓项

- 旧站数据迁移工具。
- RSS 自定义订阅管理、H&R / 保种考核、MediaInfo 与资源截图。
- 徽章市场、竞猜、复杂活动和私信系统。
