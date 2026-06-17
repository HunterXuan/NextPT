// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTorrentReward is the golang structure of table catalog_torrent_reward for DAO operations like Where/Data.
type CatalogTorrentReward struct {
	g.Meta    `orm:"table:catalog_torrent_reward, do:true"`
	Id        any         //
	UserId    any         // 赞赏者ID
	TorrentId any         // 种子ID
	Amount    any         // 赞赏的Bonus数量
	CreatedAt *gtime.Time //
}
