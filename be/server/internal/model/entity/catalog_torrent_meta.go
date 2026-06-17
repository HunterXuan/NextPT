// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTorrentMeta is the golang structure for table catalog_torrent_meta.
type CatalogTorrentMeta struct {
	Id            uint64      `json:"id"            orm:"id"             description:""`
	TorrentId     uint64      `json:"torrentId"     orm:"torrent_id"     description:""`
	ImdbId        string      `json:"imdbId"        orm:"imdb_id"        description:""`
	ImdbRating    float64     `json:"imdbRating"    orm:"imdb_rating"    description:""`
	DoubanId      string      `json:"doubanId"      orm:"douban_id"      description:""`
	DoubanRating  float64     `json:"doubanRating"  orm:"douban_rating"  description:""`
	BangumiId     string      `json:"bangumiId"     orm:"bangumi_id"     description:""`
	BangumiRating float64     `json:"bangumiRating" orm:"bangumi_rating" description:""`
	TmdbId        string      `json:"tmdbId"        orm:"tmdb_id"        description:""`
	TmdbRating    float64     `json:"tmdbRating"    orm:"tmdb_rating"    description:""`
	Extra         *gjson.Json `json:"extra"         orm:"extra"          description:"其他元数据 (JSON)"`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""`
}
