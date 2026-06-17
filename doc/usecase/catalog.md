# Catalog (资源编目) 域设计

## 核心实体 (Domain Entities)
- `Torrent` (种子)
- `TorrentFile` (种子文件列表)
- `Category` (分类)
- `Tag` / `TagGroup` (标签与标签组)
- `Subtitle` (字幕)
- `Bookmark` (收藏)

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
  * **参数概述**: `page`, `size`, `category_id`, `keyword`, `tag_ids`, `bookmarked` (是否仅看收藏)
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
