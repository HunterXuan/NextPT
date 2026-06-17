// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTorrentBookmark is the golang structure of table catalog_torrent_bookmark for DAO operations like Where/Data.
type CatalogTorrentBookmark struct {
	g.Meta    `orm:"table:catalog_torrent_bookmark, do:true"`
	Id        any         //
	UserId    any         //
	TorrentId any         //
	CreatedAt *gtime.Time //
}
