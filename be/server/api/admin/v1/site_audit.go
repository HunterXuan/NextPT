package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"
)

type SiteAuditListReq struct {
	g.Meta `path:"/site/audits" method:"get" tags:"AdminSite" summary:"查询审计日志" perm:"admin:site/audit:*"`
	adminin.SiteAuditListInp
}

type SiteAuditListRes struct {
	adminout.SiteAuditListOut
}
