# Site (站点展示与配置) 域设计

> **Site 域定位**：Site 是管理全站级状态、显示静态聚合内容的核心模块。它不仅维护面向前端展示的百科、常见问题、规则，更是掌控 NextPT 全站系统级配置命脉的地方。

## 核心实体 (Domain Entities)
- `SiteConfig` (全局动态配置项)
- `SiteAudit` (管理员及高危系统操作审计)
- `SiteAnnouncement` (全站公告)
- `SiteMessage` (站内消息 / 通知)

Forum 域中的公告节点只用于社区讨论和长期沉淀，不承担全站公告系统职责。全站公告先保持轻量，只负责官方信息发布和用户已读状态；站内消息用于用户个人事件触达，也不与论坛私信或工单混用。

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

### 2. SiteAnnouncement (全站公告)

> **定位**：全站公告是站点向用户发布重要信息的官方通道，包括维护通知、规则变动、活动说明和高优先级提醒。

* **SiteAnnouncementDomain**
  * `ListPublished(ctx, userId, page, size)`：查询已发布公告，并带出当前用户是否已读。
  * `MarkRead(ctx, userId, announcementId)`：标记当前用户已读。
  * `AdminList(ctx, page, size, status)`：后台分页查询公告。
  * `AdminCreate(ctx, data)` / `AdminUpdate(ctx, id, data)` / `AdminDelete(ctx, id)`：后台维护公告。
* **SiteAnnouncementUsecase**
  * 面向普通用户提供公告列表和已读操作。
  * 后台写入动作由 Admin usecase 转调 domain，并写入 `SiteAudit`。

#### API

- `GET /api/site/announcements`
  - 权限：`read:site/announcement:*`
  - 返回已发布公告列表，默认按发布时间倒序。
- `POST /api/site/announcements/{id}:read`
  - 权限：`read:site/announcement:*`
  - 标记公告已读。
- `GET /api/admin/site/announcements`
  - 权限：`admin:site/announcement:*`
  - 后台公告列表。
- `POST /api/admin/site/announcements`
  - 权限：`admin:site/announcement:*`
  - 创建公告。
- `PATCH /api/admin/site/announcements/{id}`
  - 权限：`admin:site/announcement:*`
  - 更新公告。
- `DELETE /api/admin/site/announcements/{id}`
  - 权限：`admin:site/announcement:*`
  - 删除公告。

#### MVP 字段

- `title`：公告标题。
- `content`：公告正文，默认按 Markdown 渲染。
- `status`：`draft` / `published` / `archived`。
- `published_at`：发布时间。
- `created_by` / `updated_by`：后台维护人。

### 3. SiteMessage (站内消息 / 通知中心)

> **定位**：站内消息是用户个人收件箱，用于系统通知和站内事件提醒。它不是论坛私信，也不是管理工单。

* **SiteMessageDomain**
  * `ListByReceiver(ctx, userId, page, size, isRead)`：查询用户消息，可通过 `isRead=false` 获得未读消息总数。
  * `MarkRead(ctx, userId, id)` / `MarkAllRead(ctx, userId)`：标记已读。
  * `Create(ctx, data)` / `BatchCreate(ctx, items)`：创建系统通知。
  * `AdminList(ctx, page, size)`：后台查看通知发送记录。
* **SiteMessageUsecase**
  * 面向用户提供通知中心和未读数。
  * 后续其它域只需要调用 `SiteMessageUsecase.Notify(...)` 或 domain 创建消息，不直接拼前端展示文案。

#### API

- `GET /api/site/messages`
  - 权限：`read:site/message:*`
  - 查询当前用户站内消息；未读数通过 `isRead=false&page=1&size=1` 返回的 `total` 获得。
- `POST /api/site/messages/{id}:read`
  - 权限：`read:site/message:*`
  - 标记单条消息已读。
- `POST /api/site/messages:readAll`
  - 权限：`read:site/message:*`
  - 标记全部已读。
- `GET /api/admin/site/messages`
  - 权限：`admin:site/message:*`
  - 后台查看通知发送记录。
- `POST /api/admin/site/messages`
  - 权限：`admin:site/message:*`
  - 后台发送系统通知。

#### MVP 字段

- `sender_id`：0 表示系统通知，管理员发送时记录管理员 ID。
- `receiver_id`：接收用户 ID。
- `title`：简短标题。
- `content`：通知正文。
- `target_type` / `target_id`：关联资源信息，由前端根据资源类型自行拼接跳转地址。
- `is_read` / `read_at`：已读状态。
