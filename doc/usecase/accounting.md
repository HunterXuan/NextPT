# Accounting (流量结算与统计) 域设计

## 核心实体 (Domain Entities)
- `UserStat` (用户全局数据统计：包含总上传、下载、做种时间等)
- `UserDailyStat` / `UserMonthlyStat` (用户流量历史按天/月快照)
- `Snatch` (用户的种子下载完成与做种记录)

## Usecase 划分及 RESTful 接口设计

> **路由前缀约定**: `/api/v1/accounting`
> 
> *注意：流量的「写操作」全部由 Tracker 域通过内部 Domain Service 或异步队列驱动，对外部前端仅暴露安全脱敏的「读操作」。*

### 1. TrafficUsecase (流量统计大盘)
负责提供用户全生命周期的流量与各项数据指标汇总及历史账单。

* **获取当前流量总览 (GetMyTraffic)**
  * **Method/Path**: `GET /users/me/traffic`
  * **参数概述**: 无 (从 ctx 获取 `user_id`)
  * **核心逻辑**: 读取 `user_stat` 表，返回上传量 (`uploaded`)、下载量 (`downloaded`)、分享率 (`share_ratio`)、总做种时间 (`seed_time`) 等信息。

* **获取流量明细趋势 (ListMyTrafficHistory)**
  * **Method/Path**: `GET /users/me/traffic-history`
  * **参数概述**: 
    - `period` (必填: `daily` 或 `monthly`)
    - `start_date`, `end_date` (选填)
  * **核心逻辑**: 基于 `user_daily_stat` 或 `user_monthly_stat` 提供趋势图表所需的数据点集。

### 2. SnatchUsecase (历史下载记录)
负责展示用户参与过的 BT 下载历史。

* **获取我的下载/做种记录列表 (ListMySnatches)**
  * **Method/Path**: `GET /users/me/snatches`
  * **参数概述**: `page`, `size`, `is_finished` (是否已完成下载), `is_active` (当前是否依然在线做种，可通过查询当前 `snatch` 的 `last_action` 与时间差阈值来粗略判定，或者跨域依赖 Tracker 域校验)。
  * **核心逻辑**: 返回对应的 `snatch` 列表，包含关联的 `torrent_id`、上传量、下载量、当前完成状态等。

* **获取单条下载记录详情 (GetMySnatch)**
  * **Method/Path**: `GET /users/me/snatches/{torrent_id}`
  * **参数概述**: `torrent_id`
  * **核心逻辑**: 当用户查看特定种子时，如果他已经下载过，可以显示他的专属进度/统计信息（如下载时长、赚取的魔力值等）。

## 跨域依赖与内部通信 (Cross-Domain)
1. **被依赖方 (Tracker Domain)**: 
   - `Tracker` 在处理客户端的心跳汇报后，异步写入 `user_stat`、`user_daily_stat` 和 `snatch` 等表，底层需调用 `Accounting` 域提供的领域服务接口，或者直接在 `Tracker` 模块完成 SQL 操作（目前为了性能 Tracker 已经内聚了这部分写操作逻辑）。
2. **被依赖方 (Admin Domain)**:
   - 管理后台需要查阅某人的流量、重置某人的流量，必须通过 `Accounting` 域提供的接口进行合规的增删改。
