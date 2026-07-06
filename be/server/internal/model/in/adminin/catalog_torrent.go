package adminin

import "github.com/gogf/gf/v2/os/gtime"

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
