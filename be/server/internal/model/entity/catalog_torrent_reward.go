// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogTorrentReward is the golang structure for table catalog_torrent_reward.
type CatalogTorrentReward struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"赞赏者ID"`
	TorrentId uint64      `json:"torrentId" orm:"torrent_id" description:"种子ID"`
	Amount    float64     `json:"amount"    orm:"amount"     description:"赞赏的Bonus数量"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
