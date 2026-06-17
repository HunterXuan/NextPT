// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTorrentBookmark is the golang structure for table catalog_torrent_bookmark.
type CatalogTorrentBookmark struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`
	TorrentId uint64      `json:"torrentId" orm:"torrent_id" description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
