package v1

import (
	"server/internal/model/in/trackerin"

	"github.com/gogf/gf/v2/frame/g"
)

type DownloadReq struct {
	g.Meta `path:"/download" tags:"Tracker" method:"get" summary:"Tracker 专用种子下载" perm:"download:catalog/torrent:*"`
	trackerin.DownloadInp
}

type DownloadRes struct {
	// 直接响应二进制文件
}
