package v1

import (
	"server/internal/model/in/forumin"
	"server/internal/model/out/forumout"

	"github.com/gogf/gf/v2/frame/g"
)

type TopicListReq struct {
	g.Meta `path:"/nodes/{slug}/topics" method:"get" tags:"Forum" summary:"获取节点下的主题列表" perm:"read:forum/topic:*"`
	forumin.TopicListInp
}

type TopicListRes struct {
	forumout.TopicListOut
}

type TopicGetHotReq struct {
	g.Meta `path:"/topics:getHot" method:"get" tags:"Forum" summary:"获取热门主题" perm:"read:forum/topic:*"`
	forumin.TopicGetHotInp
}

type TopicGetHotRes struct {
	forumout.TopicHotListOut
}

type TopicDetailReq struct {
	g.Meta `path:"/topics/{id}" method:"get" tags:"Forum" summary:"获取主题详情" perm:"read:forum/topic:*"`
	forumin.TopicDetailInp
}

type TopicDetailRes struct {
	forumout.TopicDetailOut
}

type TopicCreateReq struct {
	g.Meta `path:"/topics" method:"post" tags:"Forum" summary:"发布新主题" perm:"create:forum/topic:*"`
	forumin.TopicCreateInp
}

type TopicCreateRes struct {
	forumout.TopicCreateOut
}

type TopicUpdateReq struct {
	g.Meta `path:"/topics/{id}" method:"patch" tags:"Forum" summary:"编辑主题" perm:"update:forum/topic:{id}"`
	forumin.TopicUpdateInp
}

type TopicUpdateRes struct{}

type TopicAppendReq struct {
	g.Meta `path:"/topics/{id}:append" method:"post" tags:"Forum" summary:"追加主题附言" perm:"update:forum/topic:{id}"`
	forumin.TopicAppendInp
}

type TopicAppendRes struct{}

type TopicToggleLikeReq struct {
	g.Meta `path:"/topics/{id}:like" method:"post" tags:"Forum" summary:"点赞/取消点赞主题" perm:"read:forum/topic:*"`
	forumin.TopicToggleLikeInp
}
type TopicToggleLikeRes struct{}

type TopicRewardReq struct {
	g.Meta `path:"/topics/{id}:reward" method:"post" tags:"Forum" summary:"打赏主题" perm:"read:forum/topic:*"`
	forumin.TopicRewardInp
}
type TopicRewardRes struct{}

type TopicRewardListReq struct {
	g.Meta `path:"/topics/{id}/rewards" method:"get" tags:"Forum" summary:"获取主题赞赏列表" perm:"read:forum/topic:*"`
	forumin.TopicRewardListInp
}

type TopicRewardListRes struct {
	forumout.TopicRewardListOut
}

type TopicReportReq struct {
	g.Meta `path:"/topics/{id}:report" method:"post" tags:"Forum" summary:"举报主题" perm:"read:forum/topic:*"`
	forumin.TopicReportInp
}
type TopicReportRes struct{}

type TopicBookmarkReq struct {
	g.Meta `path:"/topics/{id}:bookmark" method:"post" tags:"Forum" summary:"收藏主题" perm:"read:forum/topic:*"`
	forumin.TopicBookmarkInp
}
type TopicBookmarkRes struct{}

type TopicUnbookmarkReq struct {
	g.Meta `path:"/topics/{id}:unbookmark" method:"post" tags:"Forum" summary:"取消收藏主题" perm:"read:forum/topic:*"`
	forumin.TopicUnbookmarkInp
}
type TopicUnbookmarkRes struct{}
