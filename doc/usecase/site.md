# Site (站点展示与配置) 域设计

> **Site 域定位**：Site 是管理全站级状态、显示静态聚合内容的核心模块。它不仅维护面向前端展示的百科、常见问题、规则，更是掌控 NextPT 全站系统级配置命脉的地方。

## 核心实体 (Domain Entities)
- `SiteAudit` (全站管理员及高危系统操作审计)
- `SystemConfig` (全局动态配置项)

> *注意：原先规划的 `News`, `Rule`, `Faq`, `Wiki` 等静态展示内容，由于可以通过直接复用 `Forum` 域的基础能力（例如建立“官方公告”、“规则说明”等只读权限板块，发布置顶帖）来实现，故不再在 Site 域冗余建表与写接口。*

## Usecase 与 Domain 划分及 RESTful 接口设计

> **路由前缀约定**: `/api/v1/site`
> 
> *注意：此域直接暴露的大部分是供普通用户或未登录用户访问的公共接口。后台核心的配置写入接口，全权交由 Admin 域管理。*

### 1. SiteConfig (站点配置服务)
> **架构分离**：配置服务按照 Clean Architecture 进行了 Domain 和 Usecase 的彻底切割。
* **SiteConfigDomain (底层数据读写)**
  * 提供 `GetConfigByGroupAndKey(ctx, group, key)` 等原子的 DB 检索。
  * **禁止缓存**，纯净的数据源对接。
  * **内部使用：Get(ctx, group, key)**：读取指定配置项（无缓存，直接查库）。
  * **内部使用：GetByPath(ctx, path)**：模仿 `gcfg` 风格读取（如 `site.config.group_name.key_name`），支持代码内置兜底默认值（无缓存，直接查库）。调用方如果有高频读取需求，需自己在其 Usecase 层加缓存。
* **SiteConfigUsecase (废弃/预留)**
  * 由于站点名称、Logo 等属于前端环境变量直传，上传体积、功能开关等属于用户登录后才下发的私有态信息，因此**全站不再设立全局对外的动态 Config API**。`SiteConfig` 完全成为后端内部的基础设施模块。
