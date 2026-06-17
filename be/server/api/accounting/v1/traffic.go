package v1

import (
	"server/internal/model/in/accountingin"
	"server/internal/model/out/accountingout"

	"github.com/gogf/gf/v2/frame/g"
)

type TrafficGetMeReq struct {
	g.Meta `path:"/users/me/traffic" method:"get" tags:"AccountingTraffic" summary:"获取我的流量总览"`
}

type TrafficGetMeRes struct {
	accountingout.TrafficGetMeOut
}

type TrafficHistoryListReq struct {
	g.Meta `path:"/users/me/traffic-history" method:"get" tags:"AccountingTraffic" summary:"获取我的流量明细趋势"`
	accountingin.TrafficHistoryListInp
}

type TrafficHistoryListRes struct {
	accountingout.TrafficHistoryListOut
}
