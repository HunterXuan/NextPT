package v1

import (
	"server/internal/model/in/trackerin"
	"server/internal/model/out/trackerout"

	"github.com/gogf/gf/v2/frame/g"
)

type AnnounceReq struct {
	g.Meta `path:"/announce" tags:"Tracker" method:"get" summary:"Tracker Announce" perm:"download:catalog/torrent:*"`
	trackerin.AnnounceInp
}

type AnnounceRes struct {
	trackerout.AnnounceOut
}
