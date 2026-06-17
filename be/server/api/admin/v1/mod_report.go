package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/modout"

	"github.com/gogf/gf/v2/frame/g"
)

type ModReportListReq struct {
	g.Meta `path:"/mod/reports" method:"get" tags:"Admin Mod" summary:"List moderation reports" perm:"admin:mod/report:*"`
	adminin.ModReportListInp
}

type ModReportListRes struct {
	modout.ListReportsOut
}

type ModReportResolveReq struct {
	g.Meta `path:"/mod/reports/{id}/resolve" method:"post" tags:"Admin Mod" summary:"Resolve a report" perm:"admin:mod/report:*"`
	adminin.ModReportResolveInp
}

type ModReportResolveRes struct{}
