// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogSubtitle is the golang structure for table catalog_subtitle.
type CatalogSubtitle struct {
	Id            uint64      `json:"id"            orm:"id"             description:""`
	TorrentId     uint64      `json:"torrentId"     orm:"torrent_id"     description:""`
	UserId        uint64      `json:"userId"        orm:"user_id"        description:""`
	Title         string      `json:"title"         orm:"title"          description:""`
	FileName      string      `json:"fileName"      orm:"file_name"      description:"原始文件名"`
	FileExt       string      `json:"fileExt"       orm:"file_ext"       description:""`
	FileSize      uint64      `json:"fileSize"      orm:"file_size"      description:""`
	StoragePath   string      `json:"storagePath"   orm:"storage_path"   description:"存储路径/S3 key"`
	Language      string      `json:"language"      orm:"language"       description:"语言代码 (如 zh-CN, en-US)"`
	DownloadCount uint        `json:"downloadCount" orm:"download_count" description:""`
	IsAnonymous   bool        `json:"isAnonymous"   orm:"is_anonymous"   description:""`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`
}
