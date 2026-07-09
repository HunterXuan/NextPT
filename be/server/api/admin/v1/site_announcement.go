package v1

import (
	"server/internal/model/in/sitein"
	"server/internal/model/out/siteout"

	"github.com/gogf/gf/v2/frame/g"
)

type SiteAnnouncementListReq struct {
	g.Meta `path:"/site/announcements" method:"get" tags:"AdminSite" summary:"查询站点公告" perm:"admin:site/announcement:*"`
	sitein.AdminAnnouncementListInp
}

type SiteAnnouncementListRes struct {
	siteout.AnnouncementListOut
}

type SiteAnnouncementCreateReq struct {
	g.Meta `path:"/site/announcements" method:"post" tags:"AdminSite" summary:"创建站点公告" perm:"admin:site/announcement:*"`
	sitein.AdminAnnouncementCreateInp
}

type SiteAnnouncementCreateRes struct {
	siteout.AnnouncementCreateOut
}

type SiteAnnouncementUpdateReq struct {
	g.Meta `path:"/site/announcements/{id}" method:"patch" tags:"AdminSite" summary:"更新站点公告" perm:"admin:site/announcement:*"`
	sitein.AdminAnnouncementUpdateInp
}

type SiteAnnouncementUpdateRes struct{}

type SiteAnnouncementDeleteReq struct {
	g.Meta `path:"/site/announcements/{id}" method:"delete" tags:"AdminSite" summary:"删除站点公告" perm:"admin:site/announcement:*"`
	sitein.AdminAnnouncementDeleteInp
}

type SiteAnnouncementDeleteRes struct{}
