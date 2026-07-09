package v1

import (
	"server/internal/model/in/sitein"
	"server/internal/model/out/siteout"

	"github.com/gogf/gf/v2/frame/g"
)

type AnnouncementListReq struct {
	g.Meta `path:"/announcements" method:"get" tags:"Site" summary:"获取站点公告" perm:"read:site/announcement:*"`
	sitein.AnnouncementListInp
}

type AnnouncementListRes struct {
	siteout.AnnouncementListOut
}

type AnnouncementReadReq struct {
	g.Meta `path:"/announcements/{id}:read" method:"post" tags:"Site" summary:"标记公告已读" perm:"read:site/announcement:*"`
	sitein.AnnouncementReadInp
}

type AnnouncementReadRes struct{}
