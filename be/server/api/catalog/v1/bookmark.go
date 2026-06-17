package v1

import (
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"

	"github.com/gogf/gf/v2/frame/g"
)

type BookmarkListReq struct {
	g.Meta `path:"/bookmarks" method:"get" tags:"Catalog" summary:"获取种子收藏列表"`
	catalogin.TorrentBookmarkListInp
}

type BookmarkListRes struct {
	catalogout.TorrentBookmarkListOut
}
