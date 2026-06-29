// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTorrent is the golang structure of table catalog_torrent for DAO operations like Where/Data.
type CatalogTorrent struct {
	g.Meta         `orm:"table:catalog_torrent, do:true"`
	Id             any         //
	InfoHash       []byte      // BitTorrent info_hash
	Name           any         // 种子标题
	SubTitle       any         // 副标题
	CategoryId     any         //
	Description    any         // 详情描述 (BBCode/Markdown)
	Nfo            []byte      // NFO 文件内容
	FileName       any         // 种子文件名
	Size           any         // 总大小 (bytes)
	FileCount      any         // 文件数量
	OwnerId        any         // 上传者
	Anonymous      any         // 匿名上传
	SpState        any         // 0=normal 1=free 2=2x 3=2xfree 4=50%off 5=2x50% 6=30%off 7=custom
	SpExpireAt     *gtime.Time // 促销到期时间
	IsFeatured     any         // 是否推荐
	IsPinned       any         // 是否置顶
	PinWeight      any         // 置顶权重 (越大越靠前)
	Seeders        any         //
	Leechers       any         //
	TimesCompleted any         //
	CommentsCount  any         //
	Views          any         //
	LikeCount      any         //
	RewardsCount   any         // 收到的赞赏次数
	RewardsAmount  any         // 收到的赞赏总金额(Bonus)
	Visible        any         //
	Banned         any         //
	LastAction     *gtime.Time // Tracker 最后活动时间
	LastReseed     *gtime.Time //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
