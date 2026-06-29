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
