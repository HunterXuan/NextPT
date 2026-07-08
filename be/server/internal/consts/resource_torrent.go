package consts

const (
	ResourceTorrentSpNormal  = iota // 0: 正常
	ResourceTorrentSpFree           // 1: 免费
	ResourceTorrentSp2x             // 2: 2倍上传
	ResourceTorrentSp2xFree         // 3: 2倍上传且免费
	ResourceTorrentSp50Off          // 4: 50% 下载
	ResourceTorrentSp2x50Off        // 5: 2倍上传 50% 下载
	ResourceTorrentSp30Off          // 6: 30% 下载
)

const (
	ResourceTorrentPromotionStateNormal      = "normal"
	ResourceTorrentPromotionStateFree        = "free"
	ResourceTorrentPromotionState2x          = "2x"
	ResourceTorrentPromotionState2xFree      = "2x_free"
	ResourceTorrentPromotionState50Percent   = "50_percent"
	ResourceTorrentPromotionState2x50Percent = "2x_50_percent"
	ResourceTorrentPromotionState30Percent   = "30_percent"
)

var ResourceTorrentPromotionStateToSp = map[string]int{
	ResourceTorrentPromotionStateNormal:      ResourceTorrentSpNormal,
	ResourceTorrentPromotionStateFree:        ResourceTorrentSpFree,
	ResourceTorrentPromotionState2x:          ResourceTorrentSp2x,
	ResourceTorrentPromotionState2xFree:      ResourceTorrentSp2xFree,
	ResourceTorrentPromotionState50Percent:   ResourceTorrentSp50Off,
	ResourceTorrentPromotionState2x50Percent: ResourceTorrentSp2x50Off,
	ResourceTorrentPromotionState30Percent:   ResourceTorrentSp30Off,
}
