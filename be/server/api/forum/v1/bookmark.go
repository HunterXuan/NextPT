package v1

import (
	"server/internal/model/in/forumin"
	"server/internal/model/out/forumout"

	"github.com/gogf/gf/v2/frame/g"
)

type BookmarkListReq struct {
	g.Meta `path:"/bookmarks" method:"get" tags:"Forum" summary:"获取主题收藏列表"`
	forumin.TopicBookmarkListInp
}

type BookmarkListRes struct {
	forumout.TopicBookmarkListOut
}
