# Tracker (BT 追踪器) 域设计

> 本域已按照 Clean Architecture 进行全面重构，遵循“实体为王，动作分离”的设计原则。

> **路由前缀约定**: `/api/tracker`

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
* **CalculateTrafficDiff(ctx, event, oldPeer) -> (diffUp, diffDn, rejectionReason)**
  * 首次 Peer 增量归零，客户端计数回退按重启后的计数计算；上传与下载增量合计超过 10 MiB/s 时，本次流量归零并返回原因。仅返回检测结果，不在 Domain 内跨域写作弊日志。
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
* **BT 客户端 RSS (Rss)**
  * **Method/Path**: `GET /rss`
  * 使用用户 passkey 和 `download:catalog/torrent:*` 权限，只经过 Tracker 客户端鉴权链路。
  * 支持页面高级筛选中的 `keyword`、`categoryIds`、`tagIds`、优惠、做种状态、体积、发布时间和外部资源 ID 等条件，筛选与种子可见范围复用 Catalog 逻辑。
  * 标签筛选保持同组 OR、跨组 AND，与 Web 列表一致。
  * RSS enclosure 指向 `GET /download`，不经过 Catalog Web 下载接口。
* **客户端白名单校验 (CheckClientWhitelist)**
  * 全量白名单内存缓存（1 小时过期），在 Usecase 层预编译 Peer ID / User-Agent 正则与内嵌例外；共享 library 只编译、匹配，不读配置、不访问数据库或缓存。
  * 只取 `tracker_agent_whitelist` 中启用的记录，逐条匹配允许正则，再匹配该记录自身启用的例外。被例外排除后继续下一条，任一规则通过则允许；全部未通过则拒绝。填写两侧正则必须同时匹配，留空一侧不限制，两侧都为空不合法。不解析版本，不依赖规则顺序；例外不是全局禁用，另一条重叠规则仍可能允许。
  * 无效正则或例外导致所属记录整体跳过并写运行日志，不把错误规则放宽为无限制规则。
  * 后台维护复用该缓存键，修改后清理本地缓存并广播所有实例失效；广播失败记录日志，其它实例仍受 1 小时 TTL 兜底。
  * `bootstrap.sql` 的客户端段按已有客户端排行保留 15 条规则，未出现在排行或计数为零的客户端不再默认支持：移除 Azureus 2/3/4、Deluge 0.x/1.x、Transmission 1.x、BitRocket、Rufus、MLDonkey、SymTorrent，以及 µTorrent 1.6/1.7/Beta。µTorrent 保留 1.x（仅稳定 1.8.x）、2.x、3.x、Mac 1.x 四组；3.x 支持标准四位版本编码和 3.5.x 字母后缀，不再兼容截断版本。Deluge 仅支持 2.0.3 起的 2.x 稳定版，统一一条规则，Peer ID 支持数字和字母版本编码且以 s 标识稳定版，User-Agent 排除预发布和开发版。BitTorrent 7.x 将 M7 与 -BT7 格式合并为一条，统一最低 7.0.1，拒绝旧 -BT7 规则曾允许的 7.0.0 和三位截断版本。保留其它排行中低使用量的客户端；qBittorrent 4.5.0 / 4.5.1 例外继续启用，移除不再支持的 Beta 例外。排行是导入策略依据，不是运行时自动淘汰规则。
  * 导入对旧规则作两项收紧：原版本分隔符中的裸 `.` 改为字面量点，qBittorrent 3.x / 4.x Peer ID 增加起始锚点。不复刻旧 HTTPS 限制与先匹配先报错的顺序逻辑；Deluge 2.x 按官方版本编码及稳定版标识匹配，不再沿用旧站的单数字版本限制。测试覆盖常见客户端、版本边界及归并后的匹配一致性。
  * bootstrap 仅用于空库初始化，客户端段使用单条 INSERT，不清空现有记录，也不单独包装事务。已有站点通过管理端维护白名单，不要执行整份开发 bootstrap。直接 SQL 更新不会广播缓存失效，更新后需重启后端或主动清理白名单缓存。

### 2.2 异步流水线: TrackerEventUsecase
负责承接庞大的写入吞吐量，将写操作通过 Redis Stream 异步化。

* **Start(ctx)** 由 `global.Init` 在服务启动时调用，初始化 Consumer Group 并启动消费者；构造函数只创建对象，不访问 Redis 或启动后台协程。
* **10 个常驻协程** 从 Redis Stream Consumer Group 消费 `AnnounceEvent`，另有 pending watcher 通过 `XAUTOCLAIM` 处理超时未确认事件。
* **handleAnnounceEvent(event)**: 
  1. 获取 Redis 分布式锁 `SET NX EX 5s` 防止双花
  2. 获取旧 Peer 快照
  3. 调用 `PeerDomain.CalculateTrafficDiff()` 算流量并取得异常原因；账务事务提交后，首次处理的异常事件通过 `ModCheaterUsecase.Record()` 记录嫌疑日志。重放事件不重复记录；日志失败只写运行日志，不中断 Peer 状态更新，也不重试已提交的账务。
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
