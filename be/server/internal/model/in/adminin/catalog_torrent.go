package adminin

import "github.com/gogf/gf/v2/os/gtime"

type CatalogTorrentReviewListInp struct {
	Keyword    string `json:"keyword" in:"query" v:"max-length:100"`
	CategoryId uint   `json:"categoryId" in:"query"`
	Status     int    `json:"status" in:"query" d:"0" v:"in:-1,0,1,2"`
	Page       int    `json:"page" in:"query" d:"1" v:"min:1"`
	Size       int    `json:"size" in:"query" d:"20" v:"min:1|max:100"`
}

type CatalogTorrentApproveInp struct {
	Id      uint64 `json:"id" in:"path" v:"required"`
	Comment string `json:"comment" v:"max-length:1000"`
}

type CatalogTorrentRejectInp struct {
	Id      uint64 `json:"id" in:"path" v:"required"`
	Comment string `json:"comment" v:"required|length:1,1000"`
}

type CatalogTorrentPinInp struct {
	Id        uint64 `json:"id" in:"path" v:"required"`
	PinWeight int    `json:"pinWeight" d:"0" v:"min:0|max:32767"`
}

type CatalogTorrentUnpinInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type CatalogTorrentFeatureInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type CatalogTorrentUnfeatureInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type CatalogTorrentPromotionInp struct {
	Id         uint64      `json:"id" in:"path" v:"required"`
	SpState    int         `json:"spState" v:"required|in:1,2,3,4,5,6"`
	SpExpireAt *gtime.Time `json:"spExpireAt"`
}

type CatalogTorrentClearPromotionInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type CatalogTorrentDeleteInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}
