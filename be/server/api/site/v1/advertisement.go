package v1

import (
	"server/internal/model/out/siteout"

	"github.com/gogf/gf/v2/frame/g"
)

type AdvertisementListReq struct {
	g.Meta `path:"/advertisements" method:"get" tags:"Site" summary:"获取站点广告" noPerm:"true"`
}

type AdvertisementListRes struct {
	siteout.AdvertisementListOut
}
