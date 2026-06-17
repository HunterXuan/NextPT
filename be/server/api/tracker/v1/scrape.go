package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ScrapeReq struct {
	g.Meta   `path:"/scrape" tags:"Tracker" method:"get" summary:"Tracker Scrape"`
	InfoHash []string `json:"info_hash" in:"query"`
}

type ScrapeRes struct {
	// 响应由控制器直接接管并返回 Bencode，不需要定义 JSON 字段
}
