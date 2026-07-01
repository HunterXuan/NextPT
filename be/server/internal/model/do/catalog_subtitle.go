// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogSubtitle is the golang structure of table catalog_subtitle for DAO operations like Where/Data.
type CatalogSubtitle struct {
	g.Meta        `orm:"table:catalog_subtitle, do:true"`
	Id            any         //
	TorrentId     any         //
	UserId        any         //
	Title         any         //
	FileName      any         // 原始文件名
	FileExt       any         //
	FileSize      any         //
	StoragePath   any         // 存储路径/S3 key
	Language      any         // 语言代码 (如 zh-CN, en-US)
	DownloadCount any         //
	Anonymous     any         // 匿名上传
	CreatedAt     *gtime.Time //
}
