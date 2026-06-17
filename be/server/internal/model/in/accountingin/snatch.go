package accountingin

import "github.com/gogf/gf/v2/os/gtime"

type SnatchListInp struct {
	Page       int   `json:"page" d:"1" v:"min:1" description:"页码"`
	Size       int   `json:"size" d:"20" v:"min:1|max:100" description:"每页数量"`
	IsFinished *bool `json:"isFinished" in:"query" description:"是否已完成下载"`
	IsActive   *bool `json:"isActive" in:"query" description:"是否在线做种"`
}

type SnatchGetInp struct {
	TorrentId uint64 `json:"torrentId" in:"path" v:"required" description:"种子ID"`
}

type RecordSnatchInp struct {
	TorrentId      uint64
	UserId         uint64
	Ipv4           string
	Ipv6           string
	Port           int
	UploadedDiff   int64
	DownloadedDiff int64
	Remaining      int64
	IsSeeder       bool
	IsFinished     bool
	EventTime      *gtime.Time
	TimeDiff       int
}
