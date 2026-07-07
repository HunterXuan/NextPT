package accountingin

import "github.com/gogf/gf/v2/os/gtime"

type TrafficHistoryListInp struct {
	Period    string      `json:"period" in:"query" v:"required|in:daily,monthly" description:"维度：daily / monthly"`
	StartDate *gtime.Time `json:"startDate" in:"query" description:"开始时间"`
	EndDate   *gtime.Time `json:"endDate" in:"query" description:"结束时间"`
}

type RecordTrafficInp struct {
	UserId            uint64
	UploadedDiff      int64
	DownloadedDiff    int64
	RawUploadedDiff   int64
	RawDownloadedDiff int64
	IsSeeder          bool
	TimeDiff          int
	EventTime         *gtime.Time
}
