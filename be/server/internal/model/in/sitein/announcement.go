package sitein

import "github.com/gogf/gf/v2/os/gtime"

type AnnouncementListInp struct {
	Page   int   `json:"page" d:"1" v:"min:1"`
	Size   int   `json:"size" d:"10" v:"min:1|max:50"`
	IsRead *bool `json:"isRead" in:"query"`
}

type AnnouncementReadInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type AdminAnnouncementListInp struct {
	Page   int  `json:"page" d:"1" v:"min:1"`
	Size   int  `json:"size" d:"20" v:"min:1|max:100"`
	Status *int `json:"status" in:"query"`
}

type AdminAnnouncementCreateInp struct {
	Title       string      `json:"title" v:"required|length:1,200"`
	Content     string      `json:"content" v:"required"`
	Status      int         `json:"status" v:"in:0,1,2"`
	PublishedAt *gtime.Time `json:"publishedAt"`
}

type AdminAnnouncementUpdateInp struct {
	Id          uint64      `json:"id" in:"path" v:"required"`
	Title       string      `json:"title" v:"required|length:1,200"`
	Content     string      `json:"content" v:"required"`
	Status      int         `json:"status" v:"in:0,1,2"`
	PublishedAt *gtime.Time `json:"publishedAt"`
}

type AdminAnnouncementDeleteInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}
