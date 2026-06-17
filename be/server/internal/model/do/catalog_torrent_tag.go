// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogTorrentTag is the golang structure of table catalog_torrent_tag for DAO operations like Where/Data.
type CatalogTorrentTag struct {
	g.Meta    `orm:"table:catalog_torrent_tag, do:true"`
	Id        any //
	TorrentId any //
	TagId     any //
}
