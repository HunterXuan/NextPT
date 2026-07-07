package v1

import (
	"server/internal/model/in/trackerin"
	"server/internal/model/out/trackerout"

	"github.com/gogf/gf/v2/frame/g"
)

type ScrapeReq struct {
	g.Meta `path:"/scrape" tags:"Tracker" method:"get" summary:"Tracker Scrape"`
	trackerin.ScrapeInp
}

type ScrapeRes struct {
	trackerout.ScrapeOut
}
