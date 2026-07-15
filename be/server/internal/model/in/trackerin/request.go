package trackerin

type AnnounceInp struct {
	InfoHash   string `json:"info_hash" in:"query" v:"required"`
	PeerId     string `json:"peer_id" in:"query" v:"required"`
	Port       int    `json:"port" in:"query" v:"required|between:1,65535"`
	Uploaded   int64  `json:"uploaded" in:"query"`
	Downloaded int64  `json:"downloaded" in:"query"`
	Left       int64  `json:"left" in:"query"`
	Event      string `json:"event" in:"query"`
	Compact    int    `json:"compact" in:"query"`
	NumWant    int    `json:"numwant" in:"query"`
}

type ScrapeInp struct {
	InfoHash []string `json:"info_hash" in:"query"`
}

type DownloadInp struct {
	Id uint64 `json:"id" in:"query" v:"required#种子ID不能为空"`
}

type RssInp struct {
	Size            int    `json:"size" in:"query" d:"50" v:"min:1|max:100" description:"返回数量"`
	Keyword         string `json:"keyword" in:"query" v:"max-length:100" description:"标题关键词"`
	CategoryIds     []uint `json:"categoryIds" in:"query" description:"分类ID列表(可选)"`
	Promotion       string `json:"promotion" in:"query" d:"all" v:"in:all,promoted,normal,free,2x,2x_free,50_percent,2x_50_percent,30_percent" description:"优惠状态"`
	PromotionOnly   bool   `json:"promotionOnly" in:"query" description:"仅返回优惠种子"`
	SeedStatus      string `json:"seedStatus" in:"query" d:"all" v:"in:all,seeded,unseeded" description:"做种状态"`
	FeaturedOnly    bool   `json:"featuredOnly" in:"query" description:"仅推荐种子"`
	MinSize         uint64 `json:"minSize" in:"query" description:"最小体积(bytes)"`
	MaxSize         uint64 `json:"maxSize" in:"query" description:"最大体积(bytes)"`
	PublishedWithin int    `json:"publishedWithin" in:"query" v:"min:0|max:3650" description:"最近发布天数"`
}
