package catalogout

import (
	"server/internal/model"
	"server/internal/model/out/economyout"

	"github.com/gogf/gf/v2/os/gtime"
)

type TorrentUploadOut struct {
	TorrentId uint64 `json:"torrentId" description:"新种子的ID"`
	InfoHash  string `json:"infoHash" description:"新种子的InfoHash(Hex)"`
}

type TorrentListItem struct {
	Id         uint64               `json:"id" description:"种子ID"`
	Name       string               `json:"name" description:"标题"`
	SubTitle   string               `json:"subTitle" description:"副标题"`
	CategoryId uint                 `json:"categoryId" description:"分类ID"`
	Size       uint64               `json:"size" description:"总大小(字节)"`
	FileCount  uint                 `json:"fileCount" description:"文件数量"`
	SpState    int                  `json:"spState" description:"促销状态"`
	SpExpireAt string               `json:"spExpireAt" description:"促销到期时间"`
	IsFeatured bool                 `json:"isFeatured" description:"是否推荐"`
	IsPinned   bool                 `json:"isPinned" description:"是否置顶"`
	PinWeight  int                  `json:"pinWeight" description:"置顶权重"`
	Seeders    uint                 `json:"seeders" description:"做种数"`
	Leechers   uint                 `json:"leechers" description:"下载数"`
	Snatched   uint                 `json:"snatched" description:"完成数"`
	LikeCount  uint                 `json:"likeCount" description:"感谢数"`
	Owner      model.IamUserSummary `json:"owner" description:"发布者"`
	Anonymous  bool                 `json:"anonymous" description:"是否匿名"`
	CreatedAt  string               `json:"createdAt" description:"发布时间"`
}

type TorrentListOut struct {
	List  []TorrentListItem `json:"list" description:"种子列表"`
	Total int               `json:"total" description:"总记录数"`
}

type TorrentRssItem struct {
	Id        uint64      `json:"id"`
	Name      string      `json:"name"`
	SubTitle  string      `json:"subTitle"`
	Category  string      `json:"category"`
	Size      uint64      `json:"size"`
	Seeders   uint        `json:"seeders"`
	Leechers  uint        `json:"leechers"`
	Snatched  uint        `json:"snatched"`
	CreatedAt *gtime.Time `json:"createdAt"`
}

type TorrentRssOut struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Language    string           `json:"language"`
	List        []TorrentRssItem `json:"list"`
}

type TorrentHotItem struct {
	Id         uint64        `json:"id" description:"种子ID"`
	Name       string        `json:"name" description:"标题"`
	SubTitle   string        `json:"subTitle" description:"副标题"`
	CategoryId uint          `json:"categoryId" description:"分类ID"`
	Category   *CategoryItem `json:"category,omitempty" description:"分类信息"`
	Size       uint64        `json:"size" description:"总大小(字节)"`
	FileCount  uint          `json:"fileCount" description:"文件数量"`
	SpState    int           `json:"spState" description:"促销状态"`
	SpExpireAt string        `json:"spExpireAt" description:"促销到期时间"`
	IsFeatured bool          `json:"isFeatured" description:"是否推荐"`
	IsPinned   bool          `json:"isPinned" description:"是否置顶"`
	PinWeight  int           `json:"pinWeight" description:"置顶权重"`
	Seeders    uint          `json:"seeders" description:"做种数"`
	Leechers   uint          `json:"leechers" description:"下载数"`
	Snatched   uint          `json:"snatched" description:"完成数"`
	LikeCount  uint          `json:"likeCount" description:"感谢数"`
	CreatedAt  string        `json:"createdAt" description:"发布时间"`
}

type TorrentHotListOut struct {
	List  []TorrentHotItem `json:"list" description:"热门种子列表"`
	Total int              `json:"total" description:"总记录数"`
}

type TorrentDetailOut struct {
	TorrentListItem
	Description   string              `json:"description" description:"详细描述"`
	ReleaseFields map[string]any      `json:"releaseFields" description:"发布结构化字段值"`
	Metadata      *TorrentMetadataOut `json:"metadata,omitempty" description:"外部资源元数据"`
	IsBookmarked  bool                `json:"isBookmarked" description:"是否收藏"`
	IsLiked       bool                `json:"isLiked" description:"是否已感谢(点赞)"`
}

type TorrentMetadataBinding struct {
	ImdbId        string  `json:"imdbId"`
	ImdbRating    float64 `json:"imdbRating"`
	DoubanId      string  `json:"doubanId"`
	DoubanRating  float64 `json:"doubanRating"`
	BangumiId     string  `json:"bangumiId"`
	BangumiRating float64 `json:"bangumiRating"`
	TmdbId        string  `json:"tmdbId"`
	TmdbType      string  `json:"tmdbType"`
	TmdbRating    float64 `json:"tmdbRating"`
}

type TorrentMetadataItem struct {
	Provider      string   `json:"provider"`
	ProviderId    string   `json:"providerId"`
	TmdbType      string   `json:"tmdbType"`
	Title         string   `json:"title"`
	OriginalTitle string   `json:"originalTitle"`
	Year          string   `json:"year"`
	ReleaseDate   string   `json:"releaseDate"`
	Overview      string   `json:"overview"`
	PosterUrl     string   `json:"posterUrl"`
	BackdropUrl   string   `json:"backdropUrl"`
	Rating        float64  `json:"rating"`
	Genres        []string `json:"genres"`
	ImdbId        string   `json:"imdbId"`
}

type TorrentMetadataOut struct {
	Binding TorrentMetadataBinding `json:"binding"`
	Data    *TorrentMetadataItem   `json:"data,omitempty"`
	Sources []TorrentMetadataItem  `json:"sources"`
}

type TorrentMetadataSearchOut struct {
	List         []TorrentMetadataItem `json:"list"`
	Page         int                   `json:"page"`
	TotalPages   int                   `json:"totalPages"`
	TotalResults int                   `json:"totalResults"`
}

type TorrentDownloadOut struct {
	Bytes    []byte `json:"-"`
	FileName string `json:"-"`
}

type TorrentRewardOut struct {
	// 空结构体
}

type TorrentRewardItem = economyout.RewardItem

type TorrentRewardListOut = economyout.RewardListOut

type TorrentBookmarkListOut struct {
	List  []TorrentListItem `json:"list" description:"种子列表"`
	Total int               `json:"total" description:"总记录数"`
}

type TorrentToggleLikeOut struct {
	IsLiked bool `json:"isLiked" description:"当前状态是否为已点赞"`
}

type TorrentLikeItem struct {
	Id        uint64               `json:"id" description:"点赞记录ID"`
	User      model.IamUserSummary `json:"user" description:"点赞者"`
	CreatedAt string               `json:"createdAt" description:"点赞时间"`
}

type TorrentLikeListOut struct {
	List  []TorrentLikeItem `json:"list" description:"点赞列表"`
	Total int               `json:"total" description:"总记录数"`
}

type TorrentUpdateOut struct {
	Success bool `json:"success"`
}

type TorrentFileItem struct {
	Path string `json:"path" description:"文件相对路径"`
	Size uint64 `json:"size" description:"文件大小"`
}

type TorrentFileListOut struct {
	List []TorrentFileItem `json:"list"`
}

type TorrentPeerItem struct {
	User       model.IamUserSummary `json:"user"`
	IsSeeder   bool                 `json:"isSeeder"`
	Uploaded   uint64               `json:"uploaded"`
	Downloaded uint64               `json:"downloaded"`
	StartedAt  string               `json:"startedAt"`
}

type TorrentPeerListOut struct {
	List []TorrentPeerItem `json:"list"`
}

type TorrentReportOut struct {
	Success bool `json:"success"`
}
