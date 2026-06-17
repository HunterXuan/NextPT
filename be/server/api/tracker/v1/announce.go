package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type AnnounceReq struct {
	g.Meta     `path:"/announce" tags:"Tracker" method:"get" summary:"Tracker Announce" perm:"download:catalog/torrent:*"`
	InfoHash   string `json:"info_hash" in:"query" v:"required"`
	PeerId     string `json:"peer_id" in:"query" v:"required"`
	Port       int    `json:"port" in:"query" v:"required|between:1,65535"`
	Uploaded   int64  `json:"uploaded" in:"query"`
	Downloaded int64  `json:"downloaded" in:"query"`
	Left       int64  `json:"left" in:"query"`
	Event      string `json:"event" in:"query"`
	Compact    int    `json:"compact" in:"query"`
	NumWant    int    `json:"numwant" in:"query"`
}

type AnnounceRes struct {
	// 响应由控制器直接接管并返回 Bencode，不需要定义 JSON 字段
}
