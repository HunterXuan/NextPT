package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type DownloadReq struct {
	g.Meta `path:"/download" tags:"Tracker" method:"get" summary:"Tracker 专用种子下载" perm:"download:catalog/torrent:*"`
	Id     uint64 `json:"id" in:"query" v:"required#种子ID不能为空"`
}

type DownloadRes struct {
	// 直接响应二进制文件
}
