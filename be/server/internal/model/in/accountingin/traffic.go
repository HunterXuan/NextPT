package accountingin

import "github.com/gogf/gf/v2/os/gtime"

type TrafficHistoryListInp struct {
	Period    string      `json:"period" in:"query" v:"required|in:daily,monthly" description:"维度：daily / monthly"`
	StartDate *gtime.Time `json:"startDate" in:"query" description:"开始时间"`
	EndDate   *gtime.Time `json:"endDate" in:"query" description:"结束时间"`
}
