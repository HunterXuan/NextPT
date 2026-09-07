# NextPT 文档索引

- `usecase/`：按 IAM、Catalog、Tracker、Forum、Economy、Site、Mod、Admin、Accounting、Sys 和 Middleware 记录当前核心用例、边界与接口设计。
- `dev/architecture_guidelines.md`：后端 Usecase / Domain 分层、事务边界、参数传递和 API 声明规范。
- `roadmap/pt_maturity.md`：已完成能力、暂缓项和后续产品方向。

接口实际挂载于 `/api`，以 `be/server/api` 中的 `g.Meta` 声明和 `be/server/internal/router` 的路由组装为准。数据库表与初始数据分别以 `be/server/storage/migration/init.sql`、`be/server/storage/migration/bootstrap.sql` 为准。
