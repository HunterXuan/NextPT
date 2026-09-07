# Admin (后台聚合) 域设计

> **Admin 域定位**：Admin 不是一个独立的业务领域，它是面向管理员的 **Backend For Frontend (BFF)**。
> 它没有专属的基础设施表（少数特定的后台菜单配置表除外），它的底层操作全部是对其他 10 个业务域的跨域调用。

## Usecase 划分及 RESTful 接口设计

> **路由前缀约定**: `/api/admin`

### 1. AdminIamRoleUsecase, AdminIamPermissionUsecase 和 AdminIamSessionUsecase (权限与用户管理后台)
* **获取全站用户列表 (ListUsers)**
  * **Method/Path**: `GET /iam/users`
* **强制销毁用户所有会话 (DeleteSessions)**
  * **Method/Path**: `DELETE /iam/users/{userId}/sessions`
* **部分更新用户信息 (UpdateUser)**
  * **Method/Path**: `PATCH /iam/users/{id}`
* **用户处罚记录 (UserMods)**
  * **Method/Path**: `GET /mod/users/{id}/mods`, `POST /mod/users/{id}/mods`, `DELETE /mod/users/{id}/mods/{modId}`
  * **核心逻辑**: 处罚记录由 Mod 域维护，过期记录由定时任务清理。
* **增量更新用户统计数据 (IncrementUserStat)**
  * **Method/Path**: `POST /iam/users/{id}:incrementStat`
* **获取用户统计详情 (GetUserStat)**
  * **Method/Path**: `GET /iam/users/{id}/stat`
* **查看用户登录记录与资源权限**
  * **Method/Path**: `GET /iam/login-logs`, `GET /iam/users/{userId}/permissions`
* **授予或撤销用户资源权限**
  * **Method/Path**: `POST /iam/users/{userId}/permissions:grant`, `POST /iam/users/{userId}/permissions:revoke`
* **管理角色 (Roles)**
  * **Method/Path**: `GET /iam/roles`, `POST /iam/roles`, `PATCH /iam/roles/{id}`, `DELETE /iam/roles/{id}`
* **管理邀请码 (Invites)**
  * **Method/Path**: `GET /iam/invites`, `POST /iam/invites:grant`, `POST /iam/invites/{id}:recycle`

### 2. AdminForumCategoryUsecase, AdminForumNodeUsecase 和 AdminForumTopicUsecase (论坛管理后台)
* **管理版块分类 (NodeCategories)**
  * **Method/Path**: `GET /forum/categories`, `POST /forum/categories`, `PATCH /forum/categories/{id}`, `DELETE /forum/categories/{id}`
* **管理版块 (Nodes)**
  * **Method/Path**: `GET /forum/nodes`, `POST /forum/nodes`, `PATCH /forum/nodes/{id}`, `DELETE /forum/nodes/{id}`
* **强制锁定主题 (LockTopic)**
  * **Method/Path**: `POST /forum/topics/{id}:lock`
* **移动主题 (MoveTopic)**
  * **Method/Path**: `POST /forum/topics/{id}:move`
* **主题管理补充**
  * **Method/Path**: `POST /forum/topics/{id}:unlock`, `POST /forum/topics/{id}:pin`, `POST /forum/topics/{id}:unpin`, `DELETE /forum/topics/{id}`

### 3. AdminModReportUsecase, AdminModCheaterUsecase 和 AdminModUserUsecase (风控与审核后台)
* **查询和处理举报 (Reports)**
  * **Method/Path**: `GET /mod/reports`, `POST /mod/reports/{id}:resolve`
* **查询和处理作弊记录 (CheaterLogs)**
  * **Method/Path**: `GET /mod/cheaters`, `POST /mod/cheaters/{id}:resolve`
* **管理组信箱 (StaffMessages)**
  * **Method/Path**: `GET /mod/staff-messages`, `PATCH /mod/staff-messages/{id}`

### 4. AdminCatalogTorrentUsecase 和 AdminCatalogTagUsecase (资源管理后台)
* **查询审核队列 (ListTorrentReviews)**
  * **Method/Path**: `GET /catalog/torrents?status=0`
* **编辑种子 (UpdateTorrent)**
  * **Method/Path**: `PATCH /catalog/torrents/{id}`
  * Staff 使用 `admin:catalog/torrent:*`，与上传者的资源级 `update:catalog/torrent:{id}` 权限分离。
* **通过种子审核 (ApproveTorrent)**
  * **Method/Path**: `POST /catalog/torrents/{id}:approve`
* **拒绝种子审核 (RejectTorrent)**
  * **Method/Path**: `POST /catalog/torrents/{id}:reject`
  * 拒绝原因必填；通过和拒绝都使用条件更新防止重复处理，并写入审计日志、发送站内通知。
* **强制删除种子 (DeleteTorrent)**
  * **Method/Path**: `DELETE /catalog/torrents/{id}`
* **种子运营操作**
  * **Method/Path**: `POST /catalog/torrents/{id}:pin`, `:unpin`, `:feature`, `:unfeature`, `:promotion`, `:clearPromotion`
* **管理标签组 (TagGroups)**
  * **Method/Path**: `GET /catalog/tag-groups`, `POST /catalog/tag-groups`, `PATCH /catalog/tag-groups/{id}`, `DELETE /catalog/tag-groups/{id}`
* **管理标签 (Tags)**
  * **Method/Path**: `POST /catalog/tag-groups/{id}/tags`, `PATCH /catalog/tags/{id}`, `DELETE /catalog/tags/{id}`
  * 标签组仍有标签、或标签仍被种子使用时拒绝删除；所有操作使用 `admin:catalog/tag:*` 并写入审计日志。

### 5. AdminSysCronUsecase (系统基建后台)
* **查阅系统审计日志 (ListSiteAudits)**
  * **Method/Path**: `GET /site/audits`
* **查阅定时任务与日志 (ListCrons / ListCronLogs)**
  * **Method/Path**: `GET /sys/crons`, `GET /sys/crons/{name}/logs`

### 6. AdminSiteConfigUsecase 和 AdminSiteAuditUsecase (站点后台)
* **站点系统配置管理 (Configs)**
  * **Method/Path**: `GET /site/configs/{group}`, `PUT /site/configs/{group}/{key}`
