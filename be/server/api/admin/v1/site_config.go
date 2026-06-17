package v1

import (
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"

	"github.com/gogf/gf/v2/frame/g"
)

type SiteConfigListReq struct {
	g.Meta `path:"/site/configs/{group}" method:"get" tags:"AdminSite" summary:"获取后台配置列表" perm:"admin:site/config:*"`
	adminin.SiteConfigListInp
}

type SiteConfigListRes struct {
	adminout.SiteConfigListOut
}

type SiteConfigUpdateReq struct {
	g.Meta `path:"/site/configs/{group}/{key}" method:"put" tags:"AdminSite" summary:"修改配置项" perm:"admin:site/config:*"`
	adminin.SiteConfigUpdateInp
}

type SiteConfigUpdateRes struct{}
