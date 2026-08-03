# Admin (后台聚合) 域设计

> **Admin 域定位**：Admin 不是一个独立的业务领域，它是面向管理员的 **Backend For Frontend (BFF)**。
> 它没有专属的基础设施表（少数特定的后台菜单配置表除外），它的底层操作全部是对其他 10 个业务域的跨域调用。

## Usecase 划分及 RESTful 接口设计

> **路由前缀约定**: `/api/v1/admin`

### 1. AdminIamRoleUsecase, AdminIamPermissionUsecase 和 AdminIamSessionUsecase (权限与用户管理后台)
* **获取全站用户列表 (ListUsers)**
  * **Method/Path**: `GET /iam/users`
* **强制销毁用户所有会话 (DeleteSessions)**
  * **Method/Path**: `DELETE /iam/users/{id}/sessions`
* **部分更新用户信息 (UpdateUser)**
  * **Method/Path**: `PATCH /iam/users/{id}`
* **直接封禁用户 (BanUser)**
  * **Method/Path**: `POST /iam/users/{id}:ban`
  * **核心逻辑**: 调用 Mod 记录风控 -> 调用 IAM 更改用户角色与状态。
* **增量更新用户统计数据 (IncrementUserStat)**
  * **Method/Path**: `POST /iam/users/{id}:incrementStat`
* **获取用户统计详情 (GetUserStat)**
  * **Method/Path**: `GET /iam/users/{id}/stat`
* **管理角色 (Roles)**
  * **Method/Path**: `GET /iam/roles`, `POST /iam/roles`, `PUT /iam/roles/{id}`, `DELETE /iam/roles/{id}`
* **管理邀请码 (Invites)**
  * **Method/Path**: `GET /iam/invites`, `POST /iam/invites`, `DELETE /iam/invites/{id}`

### 2. AdminForumCategoryUsecase, AdminForumNodeUsecase 和 AdminForumTopicUsecase (论坛管理后台)
* **管理版块分类 (NodeCategories)**
  * **Method/Path**: `GET /forum/categories`, `POST /forum/categories`, `PUT /forum/categories/{id}`, `DELETE /forum/categories/{id}`
* **管理版块 (Nodes)**
  * **Method/Path**: `GET /forum/nodes`, `POST /forum/nodes`, `PUT /forum/nodes/{id}`, `DELETE /forum/nodes/{id}`
* **强制锁定主题 (LockTopic)**
  * **Method/Path**: `POST /forum/topics/{id}:lock`
* **移动主题 (MoveTopic)**
  * **Method/Path**: `POST /forum/topics/{id}:move`

### 3. AdminModReportUsecase, AdminModCheaterUsecase 和 AdminModUserUsecase (风控与审核后台)
* **查询和处理举报 (Reports)**
  * **Method/Path**: `GET /mod/reports`, `POST /mod/reports/{id}:resolve`
* **查询和处理作弊记录 (CheaterLogs)**
  * **Method/Path**: `GET /mod/cheaters`, `POST /mod/cheaters/{id}:resolve`
* **手动添加处罚记录 (ApplyMod)**
  * **Method/Path**: `POST /mod/users/{id}:applyMod`
* **撤销处罚记录 (RemoveMod)**
  * **Method/Path**: `POST /mod/users/{id}:removeMod`

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
* **管理标签组 (TagGroups)**
  * **Method/Path**: `GET /catalog/tag-groups`, `POST /catalog/tag-groups`, `PATCH /catalog/tag-groups/{id}`, `DELETE /catalog/tag-groups/{id}`
* **管理标签 (Tags)**
  * **Method/Path**: `POST /catalog/tag-groups/{id}/tags`, `PATCH /catalog/tags/{id}`, `DELETE /catalog/tags/{id}`
  * 标签组仍有标签、或标签仍被种子使用时拒绝删除；所有操作使用 `admin:catalog/tag:*` 并写入审计日志。

### 5. AdminSysCronUsecase (系统基建后台)
* **查阅系统审计日志 (ListSiteAudits)**
  * **Method/Path**: `GET /site/audits`
* **查阅定时任务日志 (ListCronLogs)**
  * **Method/Path**: `GET /sys/cron-logs`

### 6. AdminSiteConfigUsecase 和 AdminSiteAuditUsecase (站点后台)
* **站点系统配置管理 (Configs)**
  * **Method/Path**: `GET /site/configs/{group}`, `PUT /site/configs/{group}/{key}`
