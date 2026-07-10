package v1

import (
	"server/internal/model/in/trackerin"

	"github.com/gogf/gf/v2/frame/g"
)

type RssReq struct {
	g.Meta `path:"/rss" tags:"Tracker" method:"get" summary:"Tracker RSS" perm:"download:catalog/torrent:*"`
	trackerin.RssInp
}

type RssRes struct{}
