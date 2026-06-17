package v1

import (
	"server/internal/model/in/accountingin"
	"server/internal/model/out/accountingout"

	"github.com/gogf/gf/v2/frame/g"
)

type SnatchListReq struct {
	g.Meta `path:"/users/me/snatches" method:"get" tags:"AccountingSnatch" summary:"获取我的下载与做种历史"`
	accountingin.SnatchListInp
}

type SnatchListRes struct {
	accountingout.SnatchListOut
}

type SnatchGetReq struct {
	g.Meta `path:"/users/me/snatches/{torrentId}" method:"get" tags:"AccountingSnatch" summary:"获取单条种子的下载记录详情"`
	accountingin.SnatchGetInp
}

type SnatchGetRes struct {
	accountingout.SnatchGetOut
}
