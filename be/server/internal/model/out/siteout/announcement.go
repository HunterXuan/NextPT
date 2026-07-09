package siteout

import "github.com/gogf/gf/v2/os/gtime"

type AnnouncementItem struct {
	Id          uint64      `json:"id"`
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	Status      int         `json:"status"`
	IsRead      bool        `json:"isRead"`
	CreatedBy   uint64      `json:"createdBy"`
	UpdatedBy   uint64      `json:"updatedBy"`
	PublishedAt *gtime.Time `json:"publishedAt"`
	CreatedAt   *gtime.Time `json:"createdAt"`
	UpdatedAt   *gtime.Time `json:"updatedAt"`
}

type AnnouncementListOut struct {
	List  []*AnnouncementItem `json:"list"`
	Total int                 `json:"total"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
}

type AnnouncementCreateOut struct {
	Id uint64 `json:"id"`
}
