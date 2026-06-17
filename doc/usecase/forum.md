# Forum (社区论坛) 域设计

## 核心实体 (Domain Entities)
- `NodeCategory` (论坛大区)
- `Node` (讨论板块)
- `Topic` (帖子/话题)
- `Reply` (回复)

## Usecase 划分及 RESTful 接口设计

> **路由前缀约定**: `/api/v1/forum`

### 1. TopicUsecase (话题应用服务)
* **发布帖子 (CreateTopic)**
  * **Method/Path**: `POST /topics`
  * **参数概述**: `node_id`, `title`, `content`
* **获取帖子列表 (ListTopics)**
  * **Method/Path**: `GET /topics`
  * **参数概述**: `node_id`, `page`, `size`
* **获取帖子详情 (GetTopic)**
  * **Method/Path**: `GET /topics/{id}`
* **追加帖子内容 (AppendTopic)**
  * **Method/Path**: `POST /topics/{id}:append`
  * **参数概述**: `content`
  * **核心逻辑**: 用户无法编辑原帖，但可以追加补充说明。追加的内容将以 JSON 数组追加到表的 `appends` 字段中。
### 2. ReplyUsecase (回复应用服务)
> 使用子资源形式表达关系
* **回复某贴 (CreateReply)**
  * **Method/Path**: `POST /topics/{topic_id}/replies`
  * **参数概述**: `content`, `quote_id`
### 3. TopicInteractionUsecase (话题互动应用服务)
处理跨域的打赏、举报、以及本域的点赞行为。
* **点赞/取消点赞帖子 (ToggleTopicLike)**
  * **Method/Path**: `POST /topics/{id}:like`
* **打赏帖子 (RewardTopic)**
  * **Method/Path**: `POST /topics/{id}:reward`
  * **核心逻辑**: 本接口为门面(Facade)，实际在内部调用 `Economy` 域的转账与打赏服务。
* **举报帖子 (ReportTopic)**
  * **Method/Path**: `POST /topics/{id}:report`
* **收藏帖子 (BookmarkTopic)**
  * **Method/Path**: `POST /topics/{id}:bookmark`
* **取消收藏帖子 (UnbookmarkTopic)**
  * **Method/Path**: `POST /topics/{id}:unbookmark`
* **获取主题收藏列表 (ListBookmarkedTopics)**
  * **Method/Path**: `GET /bookmarks` (基于前缀即为 `/api/v1/forum/bookmarks`)

### 4. ReplyInteractionUsecase (回复互动应用服务)
* **点赞/取消点赞回复 (ToggleReplyLike)**
  * **Method/Path**: `POST /replies/{id}:like`
* **打赏回复 (RewardReply)**
  * **Method/Path**: `POST /replies/{id}:reward`
* **举报回复 (ReportReply)**
  * **Method/Path**: `POST /replies/{id}:report`

### 5. NodeUsecase (板块应用服务)
* **获取板块树 (ListNodeTree)**
  * **Method/Path**: `GET /nodes`
