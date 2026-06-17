// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogComment is the golang structure of table catalog_comment for DAO operations like Where/Data.
type CatalogComment struct {
	g.Meta      `orm:"table:catalog_comment, do:true"`
	Id          any         //
	TargetType  any         // torrent/offer/request
	TargetId    any         //
	UserId      any         //
	Content     any         // 评论内容 (Markdown/BBCode)
	RewardCount any         // 获得的魔力值打赏
	LikeCount   any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
