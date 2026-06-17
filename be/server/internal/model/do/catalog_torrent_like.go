// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTorrentLike is the golang structure of table catalog_torrent_like for DAO operations like Where/Data.
type CatalogTorrentLike struct {
	g.Meta    `orm:"table:catalog_torrent_like, do:true"`
	Id        any         //
	UserId    any         //
	TorrentId any         //
	CreatedAt *gtime.Time //
}
