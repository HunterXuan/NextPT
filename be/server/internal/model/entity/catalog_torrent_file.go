// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// CatalogTorrentFile is the golang structure for table catalog_torrent_file.
type CatalogTorrentFile struct {
	Id        uint64 `json:"id"        orm:"id"         description:""`
	TorrentId uint64 `json:"torrentId" orm:"torrent_id" description:""`
	FilePath  string `json:"filePath"  orm:"file_path"  description:"文件路径"`
	Size      uint64 `json:"size"      orm:"size"       description:""`
}
