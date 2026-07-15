# Catalog (资源编目) 域设计

## 核心实体 (Domain Entities)
- `Torrent` (种子)
- `TorrentFile` (种子文件列表)
- `Category` (分类)
- `Tag` / `TagGroup` (标签与标签组)
- `Subtitle` (字幕)
- `Bookmark` (收藏)
- `Request` (求种 / 续种请求)

## Usecase 划分及 RESTful 接口设计

> **路由前缀约定**: `/api/v1/catalog`

### 1. TorrentPublishUsecase (种子发布与管理)
* **发布种子 (CreateTorrent)**
  * **Method/Path**: `POST /torrents`
  * **参数概述**: `file` (Torrent 文件), `name`, `sub_title`, `category_id`, `description`, `anonymous`
  * **核心逻辑**: Bencode 解析 -> S3 存储 -> 写表。
* **局部更新种子信息 (UpdateTorrent)**
  * **Method/Path**: `PATCH /torrents/{id}`
* **删除种子 (DeleteTorrent)**
  * **Method/Path**: `DELETE /torrents/{id}`

### 2. TorrentQueryUsecase (种子查阅)
* **获取种子列表 (ListTorrents)**
  * **Method/Path**: `GET /torrents`
  * **参数概述**: `page`, `size`, `keyword`, `categoryIds`, `promotion`, `seedStatus`, `featuredOnly`, `minSize`, `maxSize`, `publishedWithin`, `sort`
  * **筛选说明**:
    * `promotion` 支持全部、有优惠、无优惠和具体优惠类型；按全站优惠覆盖后的实际生效状态查询。
    * `seedStatus` 支持有做种和无做种，使用种子表中的 Tracker 缓存统计字段。
    * `publishedWithin` 表示最近发布天数；`minSize` / `maxSize` 使用字节。
    * `sort` 支持发布时间、做种数、下载数、完成数和体积排序；普通列表始终优先展示置顶种子。
  * RSS 查询复用同一个筛选对象，保证页面筛选与 BT 客户端订阅条件一致。
* **获取种子详情 (GetTorrent)**
  * **Method/Path**: `GET /torrents/{id}`
* **获取种子内部文件列表 (ListTorrentFiles)**
  * **Method/Path**: `GET /torrents/{id}/files`
* **获取种子做种/下载者列表 (ListTorrentPeers)**
  * **Method/Path**: `GET /torrents/{id}/peers`
* **获取种子收藏列表 (ListBookmarkedTorrents)**
  * **Method/Path**: `GET /bookmarks` (基于前缀即为 `/api/v1/catalog/bookmarks`)

### 3. TorrentInteractionUsecase (种子互动)
* **下载私有种子文件 (DownloadTorrent)**
  * **Method/Path**: `GET /torrents/{id}:download`
* **收藏种子 (BookmarkTorrent)**
  * **Method/Path**: `POST /torrents/{id}:bookmark`
* **取消收藏 (UnbookmarkTorrent)**
  * **Method/Path**: `POST /torrents/{id}:unbookmark` (或 `DELETE /bookmarks/{id}`)
* **赞赏种子 (RewardTorrent)**
  * **Method/Path**: `POST /torrents/{id}:reward`
* **获取赞赏记录 (ListTorrentRewards)**
  * **Method/Path**: `GET /torrents/{id}/rewards`
* **点赞种子 (LikeTorrent)**
  * **Method/Path**: `POST /torrents/{id}:like`
* **获取点赞记录 (ListTorrentLikes)**
  * **Method/Path**: `GET /torrents/{id}/likes`
* **举报种子 (ReportTorrent)**
  * **Method/Path**: `POST /torrents/{id}:report`

### 4. TorrentCommentUsecase (种子评论)
* **发表评论 (CreateComment)**
  * **Method/Path**: `POST /torrents/{id}/comments`
* **获取评论列表 (ListComments)**
  * **Method/Path**: `GET /torrents/{id}/comments`
* **举报评论 (ReportComment)**
  * **Method/Path**: `POST /torrents/{id}/comments/{cid}:report`
* **点赞/取消点赞评论 (ToggleCommentLike)**
  * **Method/Path**: `POST /torrents/{id}/comments/{cid}:like`
* **打赏评论 (RewardComment)**
  * **Method/Path**: `POST /torrents/{id}/comments/{cid}:reward`

> 普通用户侧不提供评论删除入口。违规评论通过举报进入 Mod 域，再由 Admin/Mod 后台处理。

### 5. SubtitleUsecase (独立字幕中心)
> **注意**: 字幕强关联于种子，但在管理与查询上提供全局视角的接口。
* **获取全局字幕列表 (ListSubtitles)**
  * **Method/Path**: `GET /subtitles`
* **获取某个种子的字幕列表 (ListTorrentSubtitles)**
  * **Method/Path**: `GET /torrents/{id}/subtitles`
* **为指定种子上传字幕 (UploadSubtitle)**
  * **Method/Path**: `POST /torrents/{id}/subtitles`
* **下载字幕 (DownloadSubtitle)**
  * **Method/Path**: `GET /subtitles/{id}:download`
* **更新字幕信息 (UpdateSubtitle)**
  * **Method/Path**: `PATCH /subtitles/{id}`
* **举报字幕 (ReportSubtitle)**
  * **Method/Path**: `POST /subtitles/{id}:report`

> 普通用户侧不提供字幕删除入口。违规字幕通过举报进入 Mod 域，再由 Admin/Mod 后台处理。

### 6. CategoryUsecase (分类与标签)
* **获取分类列表 (ListCategories)**
  * **Method/Path**: `GET /categories`
* **获取标签组 (ListTagGroups)**
  * **Method/Path**: `GET /tag-groups`

### 7. RequestUsecase (求种与续种)

> 求种和续种共用 `catalog_request` 实体，通过 `request_type` 区分。请求奖励在创建时从请求人余额扣除，完成时发给认领人，取消时退回请求人。所有余额和状态变更必须由 RequestUsecase 在同一事务内编排，RequestDomain 与 EconomyBonusDomain 只提供原子操作。

* **获取请求列表 (ListRequests)**
  * **Method/Path**: `GET /requests`
  * **参数概述**: `page`, `size`, `keyword`, `requestType`, `status`, `categoryId`, `view`
  * `view=created` 仅查询当前用户创建的请求，`view=claimed` 仅查询当前用户认领的请求。
* **获取请求详情 (GetRequest)**
  * **Method/Path**: `GET /requests/{id}`
* **创建请求 (CreateRequest)**
  * **Method/Path**: `POST /requests`
  * 求种需要 `categoryId`、`title`、`description` 和 `rewardAmount`。
  * 续种需要 `targetTorrentId`、`description` 和 `rewardAmount`；目标种子必须可见且当前没有做种者，也不能存在其他进行中的续种请求。
  * 前端提供固定的魔力档位，后端仍校验奖励金额为正数、精度合法且余额充足。
  * 创建成功后向请求人授予 `update:catalog/request:{id}` 资源权限。
* **认领请求 (ClaimRequest)**
  * **权限**: `read:catalog/request:*`；认领资格由 Usecase 按请求状态、请求人和发布权限继续校验。
  * **Method/Path**: `POST /requests/{id}:claim`
  * 请求人不能认领自己的请求；求种认领人必须具有种子发布权限。
  * 求种认领有效期为 72 小时，续种认领有效期为 24 小时。
  * 认领成功后向认领人授予 `update:catalog/request:{id}`，用于提交结果；放弃或过期时回收。
* **放弃认领 (AbandonRequest)**
  * **权限**: `read:catalog/request:*`；仅当前认领人可以放弃。
  * **Method/Path**: `POST /requests/{id}:abandon`
  * 仅当前认领人可以在提交结果前放弃，状态重新回到开放。
* **提交结果 (SubmitRequest)**
  * **权限**: `update:catalog/request:{id}`。
  * **Method/Path**: `POST /requests/{id}:submit`
  * 求种提交 `resultTorrentId`，后端校验种子存在且可见。
  * 续种不提交结果种子，后端校验当前认领人已经是目标种子的活跃 Seeder。
  * 提交成功后回收认领人的资源 update 权限。
* **确认完成 (CompleteRequest)**
  * **用户接口权限**: `update:catalog/request:{id}`。
  * **Method/Path**: `POST /requests/{id}:complete`
  * 仅请求人可以通过用户接口确认；确认后向认领人发放奖励。
  * Staff 代确认使用 `POST /admin/catalog/requests/{id}:complete`，权限为 `admin:catalog/request:*`。
  * Staff 操作通过 AdminCatalogRequestUsecase 调用统一事务，并写入 `catalog_request` 审计日志。
* **取消请求 (CancelRequest)**
  * **用户接口权限**: `update:catalog/request:{id}`。
  * **Method/Path**: `POST /requests/{id}:cancel`
  * 请求人只能取消尚未认领的请求；取消后奖励退回请求人。
  * Staff 取消使用 `POST /admin/catalog/requests/{id}:cancel`，权限为 `admin:catalog/request:*`，可以取消任何未完成请求。
  * 完成或取消后回收该请求相关的用户资源权限。

#### 状态流

```text
open -> claimed -> submitted -> completed
  |         |
  |         +-> open (主动放弃或认领超时)
  +---------------------------> cancelled
```

#### 认领超时

`sys` 域定时任务周期性释放 `claim_expires_at` 已经过期且仍处于 `claimed` 状态的请求。释放采用带状态和认领人条件的原子更新，避免与提交动作竞争。释放成功后通知请求人和原认领人。

#### 评论

请求评论复用 `catalog_comment`，使用 `target_type=catalog_request`。API 保持请求资源层级：

* `GET /requests/{id}/comments`
* `POST /requests/{id}/comments`
* `POST /requests/{id}/comments/{cid}:like`
* `POST /requests/{id}/comments/{cid}:reward`
* `POST /requests/{id}/comments/{cid}:report`
