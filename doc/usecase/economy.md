# Economy (经济系统) 域设计

## 核心实体 (Domain Entities)
- `UserBonus` (用户魔力值余额，存储于 `user_stat.bonus`)
- `BonusLog` (魔力值变动流水)

## Usecase 划分及 RESTful 接口设计

> **路由前缀约定**: `/api/v1/economy`

### 1. BonusUsecase (魔力值应用服务)
* **获取我的魔力值与流水 (ListMyBonusLogs)**
  * **Method/Path**: `GET /users/me/bonus-logs`
  * **参数概述**: `page`, `size`, `action` (可选过滤类型)
  * **核心逻辑**: 展示魔力值收支明细。通过传入 action 参数，可以代替原本独立的赞赏查询接口，比如传入 action="seed_bonus" 过滤挂机收益，或者 "torrent_reward_sent" 过滤赞赏。

* **获取当前用户每小时预期魔力值 (GetMyHourlyBonus)**
  * **Method/Path**: `GET /users/me/hourly-bonus`
  * **参数概述**: 无
  * **核心逻辑**: 根据用户当前的所有活跃做种记录（批量查询避免 N+1），套用 NexusPHP 积分公式，实时计算该用户当前的做种“时薪” (Bonus/Hour)。

## Domain Service（内部调用，不直接暴露 API）

> 以下方法由其他域（Catalog、Cron、Admin 等）通过 `service.EconomyBonusDomain()` 调用，**所有魔力值变动必须经由此处**，禁止直接操作 `dao.UserStat`。

* **AddBonus(ctx, userId, amount, action, targetType, targetId, remark)**
  * 为指定用户增加魔力值（支持正负数），同时自动写入 `bonus_log` 流水。
  * 调用方：Cron 定时做种结算、Admin 手动调整余额等。

* **TransferBonus(ctx, fromUserId, toUserId, amount, targetType, targetId, remarkFrom, remarkTo)**
  * 在两个用户之间转移魔力值（扣除发送方、增加接收方），并为双方各写一条 `bonus_log`。
  * 调用方：Catalog 域的种子打赏 (`POST /torrents/{id}:reward`)。

* **DistributeBonusPoints(ctx)**
  * 定时任务入口，按 NexusPHP 公式批量计算并发放所有活跃做种用户的魔力值。
  * 调用方：Cron 定时任务。

* **CalculateHourlyBonus(ctx, userId) -> float64**
  * 纯计算函数，返回指定用户当前每小时可获得的魔力值预估。
  * 与 `DistributeBonusPoints` 共享底层公式 `calculateBonusForPeers`，保证算法一致性。
  * 调用方：BonusUsecase (`GetMyHourlyBonus`)。
