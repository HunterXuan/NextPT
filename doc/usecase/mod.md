# Mod (风控与审核) 域设计

> **架构定位**: Mod 域负责举报、处罚、作弊记录和管理组信箱。风控记录主要由其他域或 Admin 域调用；管理组信箱则直接提供用户和管理端的轻量 API。

## 核心实体 (Domain Entities)
- `Report` (用户举报)
- `UserMod` (系统处罚记录，如警告、封禁)
- `CheaterLog` (作弊侦测日志)
- `StaffMessage` (用户发给管理组的消息)

## Usecase 划分及跨域调用设计

### 1. ReportUsecase (举报流转服务)
* **提交举报 (CreateReport)** 
  * **入口**: `Catalog` 域 (种子、评论举报), `Forum` 域 (帖子举报) 内部跨域调用。
* **获取待处理举报 / 处理举报 (List/Resolve Reports)**
  * **入口**: 由 `Admin` 域组装提供给管理后台使用。

### 2. UserModUsecase (用户处罚与警告服务)
* **对用户执行处罚 (ApplyMod)**
  * **核心逻辑**: 在 DB 写入处罚记录。
  * **入口**: 由 `Admin` 域调用，Admin 域负责联动 `IAM` 域禁用账户。
* **解除处罚 (RemoveMod)**
  * **入口**: 由 `Admin` 域调用。
* **获取某用户的历史处罚记录 (ListUserMods)**
  * **入口**: 由 `Admin` 域调用（版主查阅），或由 `IAM` 域的 `GET /users/me/mods` 跨域调用（用户查阅自己的历史）。

### 3. AntiCheatUsecase (防作弊监控服务)
* **记录作弊嫌疑 (RecordCheaterLog)**
  * **入口**: 提供给 `Tracker` 域跨域调用的底层接口。当探测到异常上报时被动写入。
* **获取 / 处理作弊嫌疑 (List/Resolve CheaterLogs)**
  * **入口**: 由 `Admin` 域组装提供给管理后台使用。

### 4. StaffMessageUsecase (管理组信箱)

> **定位**：这是一个轻量的管理组收件箱，不引入工单分类、关联目标、申诉流程或复杂会话。用户发送一条消息，管理组处理，可选择填写回复，用户在自己的列表中查看处理结果。

* **用户侧**
  * 用户只能查看自己发出的消息。
  * 新消息创建后状态为 `pending`。
  * 管理组处理后状态为 `processed`，回复内容可填可不填。
  * 管理组处理时，通过 SiteMessage 通知发信用户；回复内容为空时只发送处理通知。
* **管理侧**
  * 管理员按状态、用户 ID分页查看所有消息。
  * `processed` 表示已处理，回复内容可填可不填。
  * 已处理消息不再修改，处理动作记录到 SiteAudit。

#### API

- `GET /api/mod/staff-messages`
  - 权限：`read:mod/staff-message:*`
  - 查询当前用户发给管理组的消息。
- `POST /api/mod/staff-messages`
  - 权限：`create:mod/staff-message:*`
  - 创建管理组消息，参数为 `subject` 和 `content`。
- `GET /api/admin/mod/staff-messages`
  - 权限：`admin:mod/staff-message:*`
  - 管理端分页查看，可按 `status`、`senderId` 筛选。
- `PATCH /api/admin/mod/staff-messages/{id}`
  - 权限：`admin:mod/staff-message:*`
  - 处理消息，可选保存回复内容；处理后状态统一为 `processed`。
