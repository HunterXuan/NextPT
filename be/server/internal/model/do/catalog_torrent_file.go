// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogTorrentFile is the golang structure of table catalog_torrent_file for DAO operations like Where/Data.
type CatalogTorrentFile struct {
	g.Meta    `orm:"table:catalog_torrent_file, do:true"`
	Id        any //
	TorrentId any //
	FilePath  any // 文件路径
	Size      any //
}
