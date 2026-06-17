package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/modout"

	"github.com/gogf/gf/v2/frame/g"
)

type ModCheaterListReq struct {
	g.Meta `path:"/mod/cheaters" method:"get" tags:"Admin Mod" summary:"List cheater logs" perm:"admin:mod/cheater:*"`
	adminin.ModCheaterListInp
}

type ModCheaterListRes struct {
	modout.ListCheaterLogsOut
}

type ModCheaterResolveReq struct {
	g.Meta `path:"/mod/cheaters/{id}/resolve" method:"post" tags:"Admin Mod" summary:"Resolve a cheater log" perm:"admin:mod/cheater:*"`
	adminin.ModCheaterResolveInp
}

type ModCheaterResolveRes struct{}
