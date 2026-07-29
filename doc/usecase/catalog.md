# Catalog (资源编目) 域设计

## 核心实体 (Domain Entities)
- `Torrent` (种子)
- `TorrentFile` (种子文件列表)
- `TorrentMeta` (种子外部资源身份绑定与评分快照)
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
  * **参数概述**: `file` (Torrent 文件), `name`, `subTitle`, `categoryId`, `description`, `releaseFields`, `metadata`, `tagIds`, `anonymous`
  * **核心逻辑**: Bencode 解析 -> 校验分类发布字段和标签 -> 在事务内写种子、文件、元数据绑定和标签关系 -> 事务外写对象存储。
* **局部更新种子信息 (UpdateTorrent)**
  * **Method/Path**: `PATCH /torrents/{id}`
* **删除种子 (DeleteTorrent)**
  * **Method/Path**: `DELETE /torrents/{id}`

### 2. TorrentQueryUsecase (种子查阅)
* **获取种子列表 (ListTorrents)**
  * **Method/Path**: `GET /torrents`
  * **参数概述**: `page`, `size`, `keyword`, `categoryIds`, `tagIds`, `promotion`, `seedStatus`, `featuredOnly`, `minSize`, `maxSize`, `publishedWithin`, `sort`, `imdbId`, `doubanId`, `bangumiId`, `tmdbId`, `tmdbType`
  * **筛选说明**:
    * `promotion` 支持全部、有优惠、无优惠和具体优惠类型；按全站优惠覆盖后的实际生效状态查询。
    * `seedStatus` 支持有做种和无做种，使用种子表中的 Tracker 缓存统计字段。
    * `publishedWithin` 表示最近发布天数；`minSize` / `maxSize` 使用字节。
    * `sort` 支持发布时间、做种数、下载数、完成数和体积排序；普通列表始终优先展示置顶种子。
    * 外部资源 ID 使用 `catalog_torrent_meta` 精确匹配；多个 ID 同时传入时按交集查询。`tmdbType` 是 `tmdbId` 的可选附加条件，用于区分电影和剧集。
    * `tagIds` 先按标签组归类；同组多个标签按 OR，跨组按 AND。列表和详情批量返回已关联标签，不逐条查询。
  * RSS 查询复用同一个筛选对象，保证页面筛选与 BT 客户端订阅条件一致。
* **获取种子详情 (GetTorrent)**
  * **Method/Path**: `GET /torrents/{id}`
  * 详情中的 `metadata` 包含外部身份绑定、固定优先级合并后的媒体资料和各来源资料。完整资料按 provider 缓存在 Redis，不在每个种子行中复制。
* **获取种子内部文件列表 (ListTorrentFiles)**
  * **Method/Path**: `GET /torrents/{id}/files`
* **获取种子做种/下载者列表 (ListTorrentPeers)**
  * **Method/Path**: `GET /torrents/{id}/peers`
* **获取种子收藏列表 (ListBookmarkedTorrents)**
  * **Method/Path**: `GET /bookmarks` (基于前缀即为 `/api/v1/catalog/bookmarks`)

### 3. TorrentMetadataUsecase (资源元数据增强)

* **搜索外部资源 (SearchMetadata)**
  * **Method/Path**: `GET /torrent-metadata:search`
  * 当前 provider 为 TMDB，参数为 `query`、`tmdbType(movie|tv)`、`page`。
  * 服务端 token 只放在服务配置中，不通过站点配置或前端暴露。
* **绑定外部资源**
  * 上传和编辑种子时通过 `metadata` JSON 保存 `imdbId`、`doubanId`、`bangumiId`、`tmdbId`、`tmdbType`。
  * 绑定只保存外部身份和评分快照；标题、年份、简介、海报、类型等可变资料由 provider 按需解析并写入 Redis。
  * provider 优先级固定为 `TMDB -> IMDb -> 豆瓣 -> Bangumi`，不随界面语言变化。四个来源均保留独立评分快照，并按顺序补齐主资料缺失字段。
  * 用户只填写 IMDb ID 且没有 TMDB 绑定时，服务端通过 TMDB `find` 接口自动补全 TMDB ID 和 movie/TV 类型。IMDb provider 从页面使用的 GraphQL JSON 获取资料和评分，并以标题页 JSON-LD 作为回退；TMDB 详情返回的 IMDb ID 也会反向补回绑定。
  * 豆瓣 provider 优先读取移动页面使用的 ReXXar JSON，移动 HTML 的 JSON-LD、Schema meta 和媒体摘要结构作为回退。能够提取 IMDb ID 时继续尝试映射 TMDB；抓取失败不会阻止种子发布、编辑或详情加载。
  * 用户只填写 Bangumi ID 时，服务端通过 Bangumi API 获取资料，并从 infobox 中提取 IMDb ID 后尝试映射 TMDB。Bangumi 请求失败同样不会阻止主流程。
* **代码边界**
  * `internal/library/metadata` 只负责请求第三方接口、解析响应并转换成统一 `Item`，不依赖 Redis、service、catalog inp/out 或持久化逻辑。
  * `CatalogMetadataUsecase` 负责缓存、固定 provider 优先级、跨来源补缺、身份映射和评分快照更新。
  * `CatalogMetadataDomain` 只负责 `catalog_torrent_meta` 的查询与写入。
* **缓存策略**
  * TMDB 详情缓存键为 `catalog:metadata:tmdb:{tmdbType}:{tmdbId}:{locale}`；IMDb、豆瓣与 Bangumi 分别使用 `catalog:metadata:imdb:{imdbId}`、`catalog:metadata:douban:{doubanId}`、`catalog:metadata:bangumi:{bangumiId}`。
  * 各 provider 默认缓存 7 至 30 天；缓存失效后重新请求 provider，并更新对应来源的评分快照。
* **服务配置**
  ```yaml
  catalog:
    metadata:
      userAgent: "Mozilla/5.0 (...) Version/17.0 Mobile/15E148 Safari/604.1"
      tmdb:
        token: "<server-side-token>"
        baseUrl: "https://api.themoviedb.org/3"
        imageBaseUrl: "https://image.tmdb.org/t/p/w500"
        cacheTtl: "168h"
      imdb:
        enabled: true
        baseUrl: "https://www.imdb.com"
        graphqlUrl: "https://api.graphql.imdb.com/"
        cacheTtl: "168h"
      douban:
        enabled: true
        baseUrl: "https://m.douban.com/movie"
        apiBaseUrl: "https://m.douban.com/rexxar/api/v2"
        cacheTtl: "720h"
      bangumi:
        enabled: true
        baseUrl: "https://api.bgm.tv/v0"
        cacheTtl: "720h"
  ```

> 元数据 provider 必须输出统一结构。新增来源时在 usecase 中明确固定优先级与补缺策略，不根据界面语言切换主数据源。

### 4. TorrentInteractionUsecase (种子互动)
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

### 5. TorrentCommentUsecase (种子评论)
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

### 6. SubtitleUsecase (独立字幕中心)
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

### 7. CategoryUsecase (分类与标签)
* **获取分类列表 (ListCategories)**
  * **Method/Path**: `GET /categories`
* **获取标签组 (ListTagGroups)**
  * **Method/Path**: `GET /tag-groups`
  * 返回标签组、适用分类和组内标签，供发布表单、筛选器及管理页面复用。

#### 标签关系规则

* 分类发布字段中 `options.source=tagGroup` 的字段会按稳定 `slug/value` 自动生成 `catalog_torrent_tag` 关系。
* 发布和编辑页同时提供独立标签选择，但隐藏已经由发布字段承载的标签组，避免同一信息出现两套控件。
* 编辑时，发布字段承载的分组以当前字段值为准；其它标签在仍适用于新分类时保留。
* 标签组 `slug` 和标签 `value` 创建后不可修改，避免分类发布配置及历史 `releaseFields` 失效。

### 8. RequestUsecase (求种与续种)

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
