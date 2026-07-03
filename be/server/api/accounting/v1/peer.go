package v1

import (
	"server/internal/model/in/accountingin"
	"server/internal/model/out/accountingout"

	"github.com/gogf/gf/v2/frame/g"
)

type PeerListReq struct {
	g.Meta `path:"/users/me/peers" method:"get" tags:"AccountingPeer" summary:"获取我的当前活动种子"`
	accountingin.PeerListInp
}

type PeerListRes struct {
	accountingout.PeerListOut
}
