// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTorrentMeta is the golang structure of table catalog_torrent_meta for DAO operations like Where/Data.
type CatalogTorrentMeta struct {
	g.Meta        `orm:"table:catalog_torrent_meta, do:true"`
	Id            any         //
	TorrentId     any         //
	ImdbId        any         //
	ImdbRating    any         //
	DoubanId      any         //
	DoubanRating  any         //
	BangumiId     any         //
	BangumiRating any         //
	TmdbId        any         //
	TmdbRating    any         //
	Extra         *gjson.Json // 其他元数据 (JSON)
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
