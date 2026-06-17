# Mod (风控与审核) 域设计

> **架构定位**: Mod 域是一个 **Headless (无外网路由)** 的纯内部枢纽域。它不直接向普通用户或客户端暴露任何 `/api/v1/mod` 路由，所有的外部交互均通过其他域的跨域调用或 Admin 域的 BFF 进行。

## 核心实体 (Domain Entities)
- `Report` (用户举报)
- `UserMod` (系统处罚记录，如警告、封禁)
- `CheaterLog` (作弊侦测日志)

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
