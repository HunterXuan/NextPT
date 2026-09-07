# Accounting (流量结算与统计) 域设计

## 核心实体 (Domain Entities)
- `IamUserStat` (用户全局数据统计：包含入账与真实上传/下载、做种时间等)
- `IamUserPeriodStat` (用户流量历史按日/月保存的增量统计)
- `Snatch` (用户的种子下载完成与做种记录)

## Usecase 划分及 RESTful 接口设计

> **路由前缀约定**: `/api/accounting`
> 
> *注意：流量的「写操作」全部由 Tracker 域通过内部 Domain Service 或异步队列驱动，对外部前端仅暴露安全脱敏的「读操作」。*

### 1. TrafficUsecase (流量统计大盘)
负责提供用户全生命周期的流量与各项数据指标汇总及历史账单。

* **获取当前流量总览 (GetMyTraffic)**
  * **Method/Path**: `GET /users/me/traffic`
  * **参数概述**: 无（使用当前认证用户）
  * **核心逻辑**: 读取 `iam_user_stat`，返回入账与真实上传/下载、分享率、总做种时间等信息。

* **获取流量明细趋势 (ListMyTrafficHistory)**
  * **Method/Path**: `GET /users/me/traffic-history`
  * **参数概述**: 
    - `period` (必填: `daily` 或 `monthly`)
    - `startDate`, `endDate` (选填)
  * **核心逻辑**: 基于 `iam_user_period_stat` 的日/月行提供趋势图表所需的数据点集，包含入账与真实上传/下载、做种时间、下载时间和魔力。

### 2. SnatchUsecase (历史下载记录)
负责展示用户参与过的 BT 下载历史。

* **获取我的下载/做种记录列表 (ListMySnatches)**
  * **Method/Path**: `GET /users/me/snatches`
  * **参数概述**: `page`, `size`, `isFinished`（是否已完成下载）。当前在线做种状态属于 Tracker 实时 Peer 数据，不作为下载历史筛选条件。
  * **核心逻辑**: 返回对应的 `snatch` 列表，包含关联的 `torrent_id`、上传量、下载量、当前完成状态等。

* **获取单条下载记录详情 (GetMySnatch)**
  * **Method/Path**: `GET /users/me/snatches/{torrentId}`
  * **参数概述**: `torrentId`
  * **核心逻辑**: 当用户查看特定种子时，如果他已经下载过，可以显示他的专属进度/统计信息（如下载时长、赚取的魔力值等）。

### 3. PeerUsecase (当前活动种子)

* **获取我的当前活动种子 (ListMyPeers)**
  * **Method/Path**: `GET /users/me/peers`
  * **参数概述**: `page`, `size`, `status`（`all`、`seeding` 或 `leeching`）
  * **核心逻辑**: 返回当前用户在 Tracker 中仍处于活动状态的下载或做种记录。

## 跨域依赖与内部通信 (Cross-Domain)
1. **被依赖方 (Tracker Domain)**: 
   - `Tracker` 在消费客户端 announce 事件后，通过 `AccountingTrafficDomain.RecordTraffic` 更新 `iam_user_stat` 与 `iam_user_period_stat`，并维护下载完成与做种记录。
2. **被依赖方 (Admin Domain)**:
   - 管理后台需要查阅某人的流量、重置某人的流量，必须通过 `Accounting` 域提供的接口进行合规的增删改。
