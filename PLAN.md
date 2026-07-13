# NextPT 当前开发计划

这份文件只记录当前阶段的工作重点。早期的全量 DDL 设计稿已经过时，数据库结构以 `be/server/storage/migration/init.sql` 为准；空库初始化数据以 `be/server/storage/migration/bootstrap.sql` 为准。

## 当前状态

- 后端核心模块已进入可支撑前端开发的阶段：IAM、RBAC、Tracker、Catalog、Forum、Mod、Admin 基础链路已经接入。
- RBAC 已支持 Web 和 Tracker 场景共用权限中间件，用户处罚会落到用户权限表并通过现有缓存权限校验生效。
- Tracker 已处理完成事件幂等、完成水位保护、peer 候选数量限制和 announce interval 调整。
- 用户个人中心基础信息已可通过 `/users/me` 获取，个人资料更新和修改密码接口已补齐；修改密码后会销毁当前用户 token，需要重新登录。
- 初始化脚本拆为两层：`init.sql` 建表，`bootstrap.sql` 写入本地开发最小可用数据。

## 已完成验收

- 已用临时空库执行 `init.sql + bootstrap.sql`，确认建表和初始化数据可执行。
- 已用 bootstrap 默认管理员走通登录、`/users/me` 等价用例、资料更新、改密码、旧 token 失效、新密码登录。
- 已验证普通用户注册后能拿到默认角色 1。
- 已验证后台基础数据可见：站点配置、Tracker 客户端白名单、种子分类、论坛节点。
- 已补 IAM 相关轻量单测：资料更新不额外加载用户、改密码旧密码错误、新旧密码相同、成功后更新 hash 并销毁 token。

## 前端 API 约定

- API 统一挂在 `/api` 下，业务分组包括 `/iam`、`/catalog`、`/forum`、`/admin` 等。
- 普通 JSON API 使用统一响应包：`code`、`message`、`data`。下载类接口会直接写文件响应，不走普通 JSON data。
- 登录接口返回 token；前端后续请求放到 `Authorization` 头。
- `/api/iam/users/me` 是前端登录态主数据源，包含用户基础信息、流量统计、`role`、`roleLevel`、`isStaff`。
- 列表接口原则上使用 `page`、`size` 输入，输出使用 `list`、`total`；前端第一阶段按这个契约接。
- 权限判断以后端 RBAC 为准；前端可用 `isStaff` 控制后台入口显示，但最终仍以接口 403 为准。

## 后端后续优化

- 补更多核心单测：tracker completed/水位/幂等、站点配置 JSON 取值、bootstrap SQL 权限约束。
- 优化热门种子 peer 列表：继续评估 Redis active set 采样、随机窗口或 Lua 层筛选。
- 完善站点配置管理：明确哪些配置允许后台编辑，哪些只保留代码默认值。
- 梳理上传资源的 owner 权限授予，确保 `update:*:{id}` 只授予具体资源，不在普通角色里配置 `update:*:*`。
- 继续收紧分页输入校验：少量接口还只有默认值，缺少统一的 `min:1|max:100`。
- 成熟 PT 站能力差距与后续产品路线见 `doc/roadmap/pt_maturity.md`；站点公告、通知中心、种子 RSS、求种与续种请求 MVP 已完成。H&R 暂缓，近期可从魔力商城、RSS 自定义订阅或高级种子搜索中选择下一项。

## 暂不进入 MVP 的内容

- 旧站数据迁移脚本。
- 徽章市场、竞猜、复杂活动系统。
- 私信系统和复杂通知聚合。
- 两步验证、邮箱验证和完整风控策略。
