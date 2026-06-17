package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"server/internal/model/in/adminin"
)

type CatalogTorrentDeleteReq struct {
	g.Meta `path:"/catalog/torrents/{id}" method:"delete" tags:"AdminCatalog" summary:"强制删除种子" perm:"admin:catalog/torrent:*"`
	adminin.CatalogTorrentDeleteInp
}

type CatalogTorrentDeleteRes struct{}
