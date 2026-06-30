package catalogout

import "github.com/gogf/gf/v2/os/gtime"

type TorrentUploadOut struct {
	TorrentId uint64 `json:"torrentId" description:"新种子的ID"`
	InfoHash  string `json:"infoHash" description:"新种子的InfoHash(Hex)"`
}

type TorrentListItem struct {
	Id         uint64 `json:"id" description:"种子ID"`
	Name       string `json:"name" description:"标题"`
	SubTitle   string `json:"subTitle" description:"副标题"`
	CategoryId uint   `json:"categoryId" description:"分类ID"`
	Size       uint64 `json:"size" description:"总大小(字节)"`
	FileCount  uint   `json:"fileCount" description:"文件数量"`
	SpState    int    `json:"spState" description:"促销状态"`
	SpExpireAt string `json:"spExpireAt" description:"促销到期时间"`
	IsFeatured bool   `json:"isFeatured" description:"是否推荐"`
	IsPinned   bool   `json:"isPinned" description:"是否置顶"`
	Seeders    uint   `json:"seeders" description:"做种数"`
	Leechers   uint   `json:"leechers" description:"下载数"`
	Snatched   uint   `json:"snatched" description:"完成数"`
	LikeCount  uint   `json:"likeCount" description:"感谢数"`
	OwnerId    uint64 `json:"ownerId" description:"发布者ID"`
	OwnerName  string `json:"ownerName" description:"发布者用户名"`
	Anonymous  bool   `json:"anonymous" description:"是否匿名"`
	CreatedAt  string `json:"createdAt" description:"发布时间"`
}

type TorrentListOut struct {
	List  []TorrentListItem `json:"list" description:"种子列表"`
	Total int               `json:"total" description:"总记录数"`
}

type TorrentDetailOut struct {
	TorrentListItem
	Description   string         `json:"description" description:"详细描述"`
	ReleaseFields map[string]any `json:"releaseFields" description:"发布结构化字段值"`
	IsBookmarked  bool           `json:"isBookmarked" description:"是否收藏"`
	IsLiked       bool           `json:"isLiked" description:"是否已感谢(点赞)"`
}

type TorrentDownloadOut struct {
	Bytes    []byte `json:"-"`
	FileName string `json:"-"`
}

type TorrentRewardOut struct {
	// 空结构体
}

type TorrentRewardSummary struct {
	UserId       uint64      `json:"userId" orm:"user_id"`
	Amount       float64     `json:"amount" orm:"amount"`
	RewardCount  uint        `json:"rewardCount" orm:"reward_count"`
	LastRewardAt *gtime.Time `json:"lastRewardAt" orm:"last_reward_at"`
}

type TorrentRewardItem struct {
	UserId       uint64  `json:"userId" description:"赞赏者ID"`
	Username     string  `json:"username" description:"赞赏者名称"`
	Amount       float64 `json:"amount" description:"累计赞赏金额"`
	RewardCount  uint    `json:"rewardCount" description:"赞赏次数"`
	LastRewardAt string  `json:"lastRewardAt" description:"最近赞赏时间"`
}

type TorrentRewardListOut struct {
	List  []TorrentRewardItem `json:"list" description:"赞赏列表"`
	Total int                 `json:"total" description:"总记录数"`
}

type TorrentBookmarkListOut struct {
	List  []TorrentListItem `json:"list" description:"种子列表"`
	Total int               `json:"total" description:"总记录数"`
}

type TorrentToggleLikeOut struct {
	IsLiked bool `json:"isLiked" description:"当前状态是否为已点赞"`
}

type TorrentLikeItem struct {
	Id        uint64 `json:"id" description:"点赞记录ID"`
	UserId    uint64 `json:"userId" description:"点赞者ID"`
	Username  string `json:"username" description:"点赞者名称"`
	CreatedAt string `json:"createdAt" description:"点赞时间"`
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
	UserId     uint64 `json:"userId"`
	Username   string `json:"username"`
	IsSeeder   bool   `json:"isSeeder"`
	Uploaded   uint64 `json:"uploaded"`
	Downloaded uint64 `json:"downloaded"`
	StartedAt  string `json:"startedAt"`
}

type TorrentPeerListOut struct {
	List []TorrentPeerItem `json:"list"`
}

type TorrentReportOut struct {
	Success bool `json:"success"`
}
