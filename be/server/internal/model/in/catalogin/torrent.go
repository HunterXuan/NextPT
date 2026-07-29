package catalogin

import "github.com/gogf/gf/v2/net/ghttp"

type TorrentUploadInp struct {
	File          *ghttp.UploadFile `json:"file" v:"required#{#catalog.torrent.file_req}" type:"file" description:"Torrent 文件"`
	Name          string            `json:"name" description:"自定义标题"`
	SubTitle      string            `json:"subTitle" description:"副标题"`
	CategoryId    uint              `json:"categoryId" v:"required#{#catalog.torrent.category_req}" description:"分类ID"`
	Description   string            `json:"description" description:"种子详情描述"`
	ReleaseFields string            `json:"releaseFields" description:"发布结构化字段 JSON"`
	Metadata      string            `json:"metadata" description:"外部资源元数据绑定 JSON"`
	TagIds        []uint            `json:"tagIds" description:"标签 ID 列表"`
	Anonymous     bool              `json:"anonymous" description:"是否匿名上传"`
}

type TorrentMetadataFilterInp struct {
	ImdbId    string `json:"imdbId" v:"max-length:20" description:"IMDb ID 精确筛选"`
	DoubanId  string `json:"doubanId" v:"max-length:20" description:"豆瓣 ID 精确筛选"`
	BangumiId string `json:"bangumiId" v:"max-length:20" description:"Bangumi ID 精确筛选"`
	TmdbId    string `json:"tmdbId" v:"max-length:20" description:"TMDB ID 精确筛选"`
	TmdbType  string `json:"tmdbType" v:"max-length:10" description:"TMDB 资源类型(movie/tv)"`
}

type TorrentListInp struct {
	TorrentMetadataFilterInp
	Page            int    `json:"page" d:"1" v:"min:1" description:"页码"`
	Size            int    `json:"size" d:"50" v:"min:1|max:100" description:"每页数量"`
	Keyword         string `json:"keyword" v:"max-length:100" description:"标题关键词"`
	CategoryIds     []uint `json:"categoryIds" description:"分类ID列表(可选)"`
	TagIds          []uint `json:"tagIds" description:"标签ID列表(同组或、跨组且)"`
	Promotion       string `json:"promotion" d:"all" v:"in:all,promoted,normal,free,2x,2x_free,50_percent,2x_50_percent,30_percent" description:"优惠状态"`
	SeedStatus      string `json:"seedStatus" d:"all" v:"in:all,seeded,unseeded" description:"做种状态"`
	FeaturedOnly    bool   `json:"featuredOnly" description:"仅推荐种子"`
	MinSize         uint64 `json:"minSize" description:"最小体积(bytes)"`
	MaxSize         uint64 `json:"maxSize" description:"最大体积(bytes)"`
	PublishedWithin int    `json:"publishedWithin" v:"min:0|max:3650" description:"最近发布天数"`
	Sort            string `json:"sort" d:"newest" v:"in:newest,oldest,seeders,leechers,completed,size_asc,size_desc" description:"排序方式"`
}

type TorrentRssInp struct {
	TorrentMetadataFilterInp
	Size            int    `json:"size" d:"50" v:"min:1|max:100" description:"返回数量"`
	Keyword         string `json:"keyword" v:"max-length:100" description:"标题关键词"`
	CategoryIds     []uint `json:"categoryIds" description:"分类ID列表(可选)"`
	TagIds          []uint `json:"tagIds" description:"标签ID列表(同组或、跨组且)"`
	Promotion       string `json:"promotion" d:"all" v:"in:all,promoted,normal,free,2x,2x_free,50_percent,2x_50_percent,30_percent" description:"优惠状态"`
	PromotionOnly   bool   `json:"promotionOnly" description:"仅返回优惠种子"`
	SeedStatus      string `json:"seedStatus" d:"all" v:"in:all,seeded,unseeded" description:"做种状态"`
	FeaturedOnly    bool   `json:"featuredOnly" description:"仅推荐种子"`
	MinSize         uint64 `json:"minSize" description:"最小体积(bytes)"`
	MaxSize         uint64 `json:"maxSize" description:"最大体积(bytes)"`
	PublishedWithin int    `json:"publishedWithin" v:"min:0|max:3650" description:"最近发布天数"`
}

type TorrentGetHotInp struct {
	Size int `json:"size" d:"5" v:"min:1|max:10" description:"返回数量"`
}

type TorrentMetadataSearchInp struct {
	Query    string `json:"query" v:"required|max-length:100" description:"搜索关键词"`
	TmdbType string `json:"tmdbType" v:"required|in:movie,tv" description:"TMDB 资源类型"`
	Page     int    `json:"page" d:"1" v:"min:1|max:500" description:"页码"`
}

type TorrentDownloadInp struct {
	Id uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
}

type TorrentGetInp struct {
	Id uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
}

type TorrentRewardInp struct {
	Id     uint64  `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
	Amount float64 `json:"amount" v:"required|min:1#{#catalog.reward.amount_req}|{#catalog.reward.amount_min}" description:"赞赏金额"`
}

type TorrentRewardListInp struct {
	Id   uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
	Page int    `json:"page" d:"1" v:"min:1" description:"页码"`
	Size int    `json:"size" d:"20" v:"min:1|max:100" description:"每页数量"`
}

type TorrentBookmarkInp struct {
	Id uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
}

type TorrentUnbookmarkInp struct {
	Id uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
}

type TorrentBookmarkListInp struct {
	Page int `json:"page" d:"1" v:"min:1" description:"页码"`
	Size int `json:"size" d:"20" v:"min:1|max:100" description:"每页数量"`
}

type TorrentToggleLikeInp struct {
	Id uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
}

type TorrentLikeListInp struct {
	Id   uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
	Page int    `json:"page" d:"1" v:"min:1" description:"页码"`
	Size int    `json:"size" d:"20" v:"max:100" description:"每页数量"`
}

type TorrentUpdateInp struct {
	Id            uint64  `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
	Name          string  `json:"name" description:"种子标题"`
	SubTitle      string  `json:"subTitle" description:"副标题"`
	CategoryId    uint    `json:"categoryId" description:"分类ID"`
	Description   string  `json:"description" description:"详情描述"`
	ReleaseFields string  `json:"releaseFields" description:"发布结构化字段 JSON"`
	Metadata      *string `json:"metadata" description:"外部资源元数据绑定 JSON"`
	TagIds        *[]uint `json:"tagIds" description:"标签 ID 列表"`
	Anonymous     *bool   `json:"anonymous" description:"匿名上传"`
}

type TorrentFileListInp struct {
	Id uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
}

type TorrentPeerListInp struct {
	Id uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
}

type TorrentReportInp struct {
	Id     uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
	Reason string `json:"reason" v:"required" description:"举报原因"`
}
