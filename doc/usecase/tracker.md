# Tracker (BT 追踪器) 域设计

> 本域已按照 Clean Architecture 进行全面重构，遵循“实体为王，动作分离”的设计原则。

## 核心实体 (Domain Entities)
- `Peer` (节点状态，存储于 Redis，包含 IP/Port/上传量/下载量/做种状态等)
- `AgentWhitelist` (客户端白名单规则，存储于 MySQL)

## 1. Domain 层 (核心底层引擎)

> 纯粹的领域逻辑封装在 `service.TrackerPeerDomain()` 中，供 Usecase 编排或跨域调用。

* **GetActivePeers(ctx, torrentId) -> []*Peer**
  * 获取种子当前全部在线 Peer（做种者 + 下载者），通过 ZRANGEBYSCORE 过滤过期节点 + MGET 批量拉取详情。
* **GetSeedingUsers(ctx) -> []uint64**
  * 从 `tracker:seeding_users` SET 获取所有当前有做种行为的用户 ID。
* **GetUserSeedingPeers(ctx, userId) -> []Peer**
  * 获取指定用户当前所有做种中的 Peer 详情。
* **CalculateTrafficDiff(ctx, event, oldPeer) -> (diffUp, diffDn)**
  * 包含**防作弊引擎**逻辑：首次 Peer 增量归零、硬限速 10MB/s 校验、异常回退检测。
* **UpsertPeer / RemovePeer**
  * 封装操作 Redis 的极度复杂的脏活累活：维护 `peer详情(JSON)`、`种子做种/下载队列(ZSET)`、`用户做种/下载反向索引(SET)`、`活跃种子池(SET)`、`活跃用户池(SET)` 等多级索引体系。

## 2. Usecase 层 (应用级编排)

### 2.1 协议交互: TrackerPeerUsecase
负责对外接收标准的 BT 客户端请求，组装领域事件。

* **客户端汇报状态 (Announce)**
  * **Method/Path**: `GET /announce`
  * **核心逻辑**: 解析双栈 IP -> 组装 `AnnounceEvent` -> 投递至 `EventUsecase` -> 拉取并洗牌活跃 Peer 返回给客户端。
* **健康度抓取 (Scrape)**
  * **Method/Path**: `GET /scrape`
  * **核心逻辑**: 复用 Catalog 的种子元数据缓存，极速返回种子完成数。
* **种子文件下载 (Download)**
  * **Method/Path**: `GET /download`
* **客户端白名单校验 (CheckClientWhitelist)**
  * 全量白名单内存缓存（1 小时过期），预编译正则。匹配规则：PeerID 前缀 + User-Agent 正则。由中间件统一调用。

### 2.2 异步流水线: TrackerEventUsecase
负责承接庞大的写入吞吐量，将写操作异步化。

* **10 个常驻协程** 从 Channel 消费 `AnnounceEvent`。
* **handleAnnounceEvent(event)**: 
  1. 获取 Redis 分布式锁 `SET NX EX 5s` 防止双花
  2. 获取旧 Peer 快照
  3. 调用 `PeerDomain.CalculateTrafficDiff()` 算流量
  4. 编排路由：调用 Accounting 记流量/快照 -> 调用 `PeerDomain.UpsertPeer / RemovePeer` 刷写 Redis 状态。

### 2.3 系统调度: TrackerSyncUsecase
消除基建词汇的定时任务编排层。

* **CleanupGhostPeers(ctx)** (每 5 分钟执行)
  * 扫描活跃种子 ZSET 中 score ≤ 当前时间的过期 Peer，使用 **Lua 原子脚本**安全驱逐，同步清理反向索引。
* **SyncTorrentData(ctx)** (每 1 分钟执行)
  * 批量计算 Redis 实时在线人数并回写 MySQL `torrent` 表，双数据源对齐合并。

## 3. 跨域依赖 (Cross-Domain)

1. **依赖 IAM 域**:
   - `CheckTrackerAuth` 中间件调用 `IamUserDomain().GetUserByPasskey()`（含 gcache 高并发本地防击穿缓存）进行鉴权。
2. **依赖 Accounting 域**:
   - 编排层消费事件时，调用 `AccountingTrafficDomain().RecordTraffic()` 记全局流量，`AccountingSnatchDomain().RecordSnatch()` 记种子级做种快照。
3. **被 Economy 域调用**:
   - 魔力值与时薪发放任务调用 `TrackerPeerDomain()` 获取全站或用户的实时做种详情列表。
