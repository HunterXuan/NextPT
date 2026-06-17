// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// CatalogTorrentTag is the golang structure for table catalog_torrent_tag.
type CatalogTorrentTag struct {
	Id        uint64 `json:"id"        orm:"id"         description:""`
	TorrentId uint64 `json:"torrentId" orm:"torrent_id" description:""`
	TagId     uint   `json:"tagId"     orm:"tag_id"     description:""`
}
